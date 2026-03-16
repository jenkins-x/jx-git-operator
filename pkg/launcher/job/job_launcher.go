package job

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/naming"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"
	"github.com/jenkins-x/jx-helpers/v3/pkg/yamls"
	"github.com/jenkins-x/jx-kube-client/v3/pkg/kubeclient"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"k8s.io/apimachinery/pkg/util/strategicpatch"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1 "k8s.io/api/batch/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v12 "k8s.io/client-go/kubernetes/typed/batch/v1"
)

type client struct {
	kubeClient kubernetes.Interface
	ns         string
	selector   string
	runner     cmdrunner.CommandRunner
}

// NewLauncher creates a new launcher for Jobs using the given kubernetes client and namespace
// if nil is passed in the kubernetes client will be lazily created
func NewLauncher(kubeClient kubernetes.Interface, ns, selector string, runner cmdrunner.CommandRunner) (launcher.Interface, error) {
	if kubeClient == nil {
		f := kubeclient.NewFactory()
		cfg, err := f.CreateKubeConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create kube config: %w", err)
		}

		kubeClient, err = kubernetes.NewForConfig(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create the kube client: %w", err)
		}

		if ns == "" {
			ns, err = kubeclient.CurrentNamespace()
			if err != nil {
				return nil, fmt.Errorf("failed to find the current namespace: %w", err)
			}
		}
	}
	if runner == nil {
		runner = cmdrunner.DefaultCommandRunner
	}
	return &client{
		kubeClient: kubeClient,
		ns:         ns,
		selector:   selector,
		runner:     runner,
	}, nil
}

// Launch launches a job for the given commit
func (c *client) Launch(opts *launcher.LaunchOptions) ([]runtime.Object, error) {
	ctx := context.Background()
	ns := opts.Repository.Namespace
	if ns == "" {
		ns = c.ns
	}
	safeGitURL := stringhelpers.SanitizeURL(opts.Repository.GitURL)
	if opts.LastCommitURL == "" && opts.Repository.GitURL != "" && opts.GitSHA != "" {
		opts.LastCommitURL = stringhelpers.UrlJoin(strings.TrimSuffix(safeGitURL, ".git"), "commit", opts.GitSHA)
	}
	safeName := naming.ToValidValue(opts.Repository.Name)
	safeSha := naming.ToValidValue(opts.GitSHA)
	selector := fmt.Sprintf("%s,%s=%s", c.selector, launcher.RepositoryLabelKey, safeName)
	jobInterface := c.kubeClient.BatchV1().Jobs(ns)
	list, err := jobInterface.List(ctx, metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil && apierrors.IsNotFound(err) {
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find Jobs in namespace %s with selector %s: %w", ns, selector, err)
	}

	var jobsForSha []v1.Job
	var activeJobs []v1.Job
	for i := range list.Items {
		r := list.Items[i]
		log.Logger().Infof("found Job %s", r.Name)

		if r.Labels[launcher.CommitShaLabelKey] == safeSha && r.Labels[launcher.RerunLabelKey] != "true" {
			jobsForSha = append(jobsForSha, r)
		}

		// is the job active
		if IsJobActive(&r) {
			activeJobs = append(activeJobs, r)
		}
	}

	if len(jobsForSha) == 0 {
		if len(activeJobs) > 0 {
			log.Logger().Infof("not creating a Job in namespace %s for repo %s sha %s yet as there is an active job %s", ns, safeName, safeSha, activeJobs[0].Name)
			return nil, nil
		}
		return c.startNewJob(ctx, opts, jobInterface, ns, safeName, safeSha, safeGitURL)
	}
	return nil, nil
}

// IsJobActive returns true if the job has not completed or terminated yet
func IsJobActive(r *v1.Job) bool {
	for _, con := range r.Status.Conditions {
		if con.Status == corev1.ConditionTrue {
			return false
		}
	}
	return true
}

// startNewJob lets create a new Job resource
func (c *client) startNewJob(ctx context.Context, opts *launcher.LaunchOptions, jobInterface v12.JobInterface, ns, safeName, safeSha, safeGitURL string) ([]runtime.Object, error) {
	log.Logger().Infof("about to create a new job for name %s and sha %s", safeName, safeSha)

	// lets see if we are using a version stream to store the git operator configuration
	folder := filepath.Join(opts.Dir, "versionStream", "git-operator")
	exists, err := files.DirExists(folder)
	if err != nil {
		return nil, fmt.Errorf("failed to check if folder exists %s: %w", folder, err)
	}
	if !exists {
		// lets try the original location
		folder = filepath.Join(opts.Dir, ".jx", "git-operator")
	}

	jobFileName := "job.yaml"
	if os.Getenv("JX_CUSTOM_FILE") == "true" {
		fileNamePath := filepath.Join(opts.Dir, ".jx", "git-operator", "filename.txt")
		exists, err = files.FileExists(fileNamePath)
		if err != nil {
			return nil, fmt.Errorf("failed to check for file %s: %w", fileNamePath, err)
		}
		if exists {
			data, err := os.ReadFile(fileNamePath)
			if err != nil {
				return nil, fmt.Errorf("failed to load file %s: %w", fileNamePath, err)
			}
			jobFileName = strings.TrimSpace(string(data))
			if jobFileName == "" {
				return nil, fmt.Errorf("the job name file %s is empty", fileNamePath)
			}
		}
	}

	fileName := filepath.Join(folder, jobFileName)
	exists, err = files.FileExists(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to find file %s in repository %s: %w", fileName, safeName, err)
	}
	if !exists {
		return nil, fmt.Errorf("repository %s does not have a Job file: %s", safeName, fileName)
	}

	jsonMap := &strategicpatch.JSONMap{}
	err = yamls.LoadFile(fileName, jsonMap)
	if err != nil {
		return nil, fmt.Errorf("failed to load Job file %s in repository %s: %w", fileName, safeName, err)
	}

	enriched, err := c.enrichJob(opts, jsonMap, safeName)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich the Job: %w", err)
	}
	converter := runtime.DefaultUnstructuredConverter
	resource := &v1.Job{}
	if err := converter.FromUnstructuredWithValidation(enriched, resource, true); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	if !opts.NoResourceApply {
		// now lets check if there is a resources dir
		resourcesDir := filepath.Join(folder, "resources")
		exists, err = files.DirExists(resourcesDir)
		if err != nil {
			return nil, fmt.Errorf("failed to check if resources directory %s exists in repository %s: %w", resourcesDir, safeName, err)
		}
		if exists {
			absDir, err := filepath.Abs(resourcesDir)
			if err != nil {
				return nil, fmt.Errorf("failed to get absolute resources dir %s: %w", resourcesDir, err)
			}

			cmd := &cmdrunner.Command{
				Name: "kubectl",
				Args: []string{"apply", "-f", absDir},
			}
			log.Logger().Infof("running command: %s", cmd.CLI())
			_, err = c.runner(cmd)
			if err != nil {
				return nil, fmt.Errorf("failed to apply resources in dir %s: %w", absDir, err)
			}
		}
	}

	// lets try use a maximum of 31 characters and a minimum of 10 for the sha
	namePrefix := trimLength(safeName, 20)

	id := uuid.New().String()
	resourceName := namePrefix + "-" + id

	resource.Name = resourceName

	if resource.Annotations == nil {
		resource.Annotations = map[string]string{}
	}
	if resource.Labels == nil {
		resource.Labels = map[string]string{}
	}
	resource.Labels[constants.DefaultSelectorKey] = constants.DefaultSelectorValue
	resource.Labels[launcher.RepositoryLabelKey] = safeName
	resource.Labels[launcher.CommitShaLabelKey] = safeSha
	if opts.LastCommitAuthor != "" {
		resource.Annotations[launcher.CommitAuthorAnnotation] = opts.LastCommitAuthor
	}
	if opts.LastCommitAuthorEmail != "" {
		resource.Annotations[launcher.CommitAuthorEmailAnnotation] = opts.LastCommitAuthorEmail
	}
	if opts.LastCommitDate != "" {
		resource.Annotations[launcher.CommitDateAnnotation] = opts.LastCommitDate
	}
	if opts.LastCommitMessage != "" {
		resource.Annotations[launcher.CommitMessageAnnotation] = opts.LastCommitMessage
	}
	if opts.LastCommitURL != "" {
		resource.Annotations[launcher.CommitURLAnnotation] = opts.LastCommitURL
	}
	if safeGitURL != "" {
		resource.Annotations[launcher.RepositoryURLAnnotation] = safeGitURL
	}

	r2, err := jobInterface.Create(ctx, resource, metav1.CreateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to create Job %s in namespace %s: %w", resourceName, ns, err)
	}
	log.Logger().Infof("created Job %s in namespace %s", resourceName, ns)
	return []runtime.Object{r2}, nil
}

func (c *client) enrichJob(opts *launcher.LaunchOptions, job *strategicpatch.JSONMap, safeName string) (strategicpatch.JSONMap, error) {
	path := filepath.Join(opts.Dir, ".jx", "git-operator", "job-overlay.yaml")
	exists, err := files.FileExists(path)
	if err != nil {
		return nil, fmt.Errorf("failed to check for file %s: %w", path, err)
	}
	if !exists {
		return nil, nil
	}
	overlay := &strategicpatch.JSONMap{}
	err = yamls.LoadFile(path, overlay)
	if err != nil {
		return nil, fmt.Errorf("failed to load Job file %s in repository %s: %w", path, safeName, err)
	}

	enriched, err := OverlayJob(job, overlay)
	if err != nil {
		return nil, fmt.Errorf("failed to apply overlay from file %s to Job: %w", path, err)
	}
	return enriched, nil
}

// OverlayJob applies the given overlay to the job
func OverlayJob(job, overlay *strategicpatch.JSONMap) (strategicpatch.JSONMap, error) {
	if overlay == nil {
		return nil, nil
	}
	schemaReferenceObj := v1.Job{}
	return strategicpatch.StrategicMergeMapPatch(*job, *overlay, schemaReferenceObj)
}

func trimLength(text string, length int) string {
	if len(text) <= length {
		return text
	}
	return text[0:length]
}
