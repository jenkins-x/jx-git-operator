package job

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"github.com/jenkins-x/jx-helpers/v3/pkg/stringhelpers"

	"github.com/google/uuid"
	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/naming"
	"github.com/jenkins-x/jx-helpers/v3/pkg/yamls"
	"github.com/jenkins-x/jx-kube-client/v3/pkg/kubeclient"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
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
func NewLauncher(kubeClient kubernetes.Interface, ns string, selector string, runner cmdrunner.CommandRunner) (launcher.Interface, error) {
	if kubeClient == nil {
		f := kubeclient.NewFactory()
		cfg, err := f.CreateKubeConfig()
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create kube config")
		}

		kubeClient, err = kubernetes.NewForConfig(cfg)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create the kube client")
		}

		if ns == "" {
			ns, err = kubeclient.CurrentNamespace()
			if err != nil {
				return nil, errors.Wrapf(err, "failed to find the current namespace")
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
func (c *client) Launch(opts launcher.LaunchOptions) ([]runtime.Object, error) {
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
		return nil, errors.Wrapf(err, "failed to find Jobs in namespace %s with selector %s", ns, selector)
	}

	var jobsForSha []v1.Job
	var activeJobs []v1.Job
	for _, r := range list.Items {
		log.Logger().Infof("found Job %s", r.Name)

		if r.Labels[launcher.CommitShaLabelKey] == safeSha && r.Labels[launcher.RerunLabelKey] != "true" {
			jobsForSha = append(jobsForSha, r)
		}

		// is the job active
		if IsJobActive(r) {
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
func IsJobActive(r v1.Job) bool {
	for _, con := range r.Status.Conditions {
		if con.Status == corev1.ConditionTrue {
			return false
		}
	}
	return true
}

// startNewJob lets create a new Job resource
func (c *client) startNewJob(ctx context.Context, opts launcher.LaunchOptions, jobInterface v12.JobInterface, ns string, safeName string, safeSha, safeGitURL string) ([]runtime.Object, error) {
	log.Logger().Infof("about to create a new job for name %s and sha %s", safeName, safeSha)

	// lets see if we are using a version stream to store the git operator configuration
	folder := filepath.Join(opts.Dir, "versionStream", "git-operator")
	exists, err := files.DirExists(folder)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to check if folder exists %s", folder)
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
			return nil, errors.Wrapf(err, "failed to check for file %s", fileNamePath)
		}
		if exists {
			data, err := ioutil.ReadFile(fileNamePath)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to load file %s", fileNamePath)
			}
			jobFileName = strings.TrimSpace(string(data))
			if jobFileName == "" {
				return nil, errors.Errorf("the job name file %s is empty", fileNamePath)
			}
		}
	}

	fileName := filepath.Join(folder, jobFileName)
	exists, err = files.FileExists(fileName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find file %s in repository %s", fileName, safeName)
	}
	if !exists {
		return nil, errors.Errorf("repository %s does not have a Job file: %s", safeName, fileName)
	}

	resource := &v1.Job{}
	err = yamls.LoadFile(fileName, resource)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load Job file %s in repository %s", fileName, safeName)
	}

	err = c.enrichJob(ctx, opts, resource, safeName)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to enrich the Job")
	}

	if !opts.NoResourceApply {
		// now lets check if there is a resources dir
		resourcesDir := filepath.Join(folder, "resources")
		exists, err = files.DirExists(resourcesDir)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to check if resources directory %s exists in repository %s", resourcesDir, safeName)
		}
		if exists {
			absDir, err := filepath.Abs(resourcesDir)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to get absolute resources dir %s", resourcesDir)
			}

			cmd := &cmdrunner.Command{
				Name: "kubectl",
				Args: []string{"apply", "-f", absDir},
			}
			log.Logger().Infof("running command: %s", cmd.CLI())
			_, err = c.runner(cmd)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to apply resources in dir %s", absDir)
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
		return nil, errors.Wrapf(err, "failed to create Job %s in namespace %s", resourceName, ns)
	}
	log.Logger().Infof("created Job %s in namespace %s", resourceName, ns)
	return []runtime.Object{r2}, nil
}

func (c *client) enrichJob(ctx context.Context, opts launcher.LaunchOptions, job *v1.Job, safeName string) error {
	path := filepath.Join(opts.Dir, ".jx", "git-operator", "job-overlay.yaml")
	exists, err := files.FileExists(path)
	if err != nil {
		return errors.Wrapf(err, "failed to check for file %s", path)
	}
	if !exists {
		return nil
	}
	overlay := &v1.Job{}
	err = yamls.LoadFile(path, overlay)
	if err != nil {
		return errors.Wrapf(err, "failed to load Job file %s in repository %s", path, safeName)
	}

	err = OverlayJob(job, overlay)
	if err != nil {
		return errors.Wrapf(err, "failed to apply overlay from file %s to Job", path)
	}
	return nil
}

// OverlayJob applies the given overlay to the job
func OverlayJob(job *v1.Job, overlay *v1.Job) error {
	if overlay == nil {
		return nil
	}
	err := mergo.Merge(job, overlay)
	if err != nil {
		return errors.Wrap(err, "error merging Job with overlay")
	}

	// mergeo can't handle container and env vars yet so lets help...
	for i := range overlay.Spec.Template.Spec.Containers {
		oc := &overlay.Spec.Template.Spec.Containers[i]

		found := false
		for j := range job.Spec.Template.Spec.Containers {
			jc := &job.Spec.Template.Spec.Containers[j]
			if jc.Name == oc.Name {
				err = overlayJobContainer(jc, oc)
				if err != nil {
					return errors.Wrapf(err, "failed to merge overlay job container %s", oc.Name)
				}
				found = true
				break
			}
		}
		if !found {
			errors.Errorf("could not find container called %s in the Job definition from the overlay", oc.Name)
		}
	}
	return nil
}

func overlayJobContainer(jc *corev1.Container, oc *corev1.Container) error {
	err := mergo.Merge(jc, oc)
	if err != nil {
		return errors.Wrap(err, "error merging Container with overlay")
	}
	for i := range oc.Env {
		oe := &oc.Env[i]

		found := false
		for j := range jc.Env {
			je := &jc.Env[j]
			if je.Name == oe.Name {
				*je = *oe
				found = true
				break
			}
		}
		if !found {
			jc.Env = append(jc.Env, *oe)
		}
	}
	return nil
}

func trimLength(text string, length int) string {
	if len(text) <= length {
		return text
	}
	return text[0:length]
}
