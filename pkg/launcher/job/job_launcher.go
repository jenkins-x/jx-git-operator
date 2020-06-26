package job

import (
	"fmt"
	"path/filepath"

	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher"
	"github.com/jenkins-x/jx-helpers/pkg/files"
	"github.com/jenkins-x/jx-helpers/pkg/kube/naming"
	"github.com/jenkins-x/jx-helpers/pkg/yamls"
	"github.com/jenkins-x/jx-kube-client/pkg/kubeclient"
	"github.com/jenkins-x/jx-logging/pkg/log"
	"github.com/pkg/errors"
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
}

// NewLauncher creates a new launcher for Jobs using the given kubernetes client and namespace
// if nil is passed in the kubernetes client will be lazily created
func NewLauncher(kubeClient kubernetes.Interface, ns string, selector string) (launcher.Interface, error) {
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
	return &client{
		kubeClient: kubeClient,
		ns:         ns,
		selector:   selector,
	}, nil
}

// Launch launches a job for the given commit
func (c *client) Launch(opts launcher.LaunchOptions) error {
	ns := opts.Repository.Namespace
	if ns == "" {
		ns = c.ns
	}
	safeName := naming.ToValidValue(opts.Repository.Name)
	safeSha := naming.ToValidValue(opts.GitSHA)
	selector := fmt.Sprintf("%s,%s=%s,%s=%s", c.selector,
		launcher.RepositoryLabelKey, safeName,
		launcher.CommitShaLabelKey, safeSha)
	jobInterface := c.kubeClient.BatchV1().Jobs(ns)
	list, err := jobInterface.List(metav1.ListOptions{
		LabelSelector: selector,
	})
	if err != nil && apierrors.IsNotFound(err) {
		err = nil
	}
	if err != nil {
		return errors.Wrapf(err, "failed to find Jobs in namespace %s with selector %s", ns, selector)
	}

	for _, r := range list.Items {
		log.Logger().Infof("found Job %s", r.Name)

		// TODO should we do anything with the status of existing jobs? e.g. report status somewhere
	}

	if len(list.Items) == 0 {
		return c.startNewJob(opts, jobInterface, ns, safeName, safeSha)
	}

	return nil
}

// startNewJob lets create a new Job resource
func (c *client) startNewJob(opts launcher.LaunchOptions, jobInterface v12.JobInterface, ns string, safeName string, safeSha string) error {
	log.Logger().Infof("about to create a new job for name %s and sha %s", safeName, safeSha)

	fileName := filepath.Join(opts.Dir, ".jx", "git-operator", "job.yaml")
	exists, err := files.FileExists(fileName)
	if err != nil {
		return errors.Wrapf(err, "failed to find file %s in repository %s", fileName, safeName)
	}
	if !exists {
		return errors.Errorf("repository %s does not have a Job file: %s", safeName, fileName)
	}

	resource := &v1.Job{}
	err = yamls.LoadFile(fileName, resource)
	if err != nil {
		return errors.Wrapf(err, "failed to load Job file %s in repository %s", fileName, safeName)
	}

	// lets try use a maximum of 31 characters and a minimum of 10 for the sha
	namePrefix := trimLength(safeName, 20)
	maxShaLen := 30 - len(namePrefix)

	resourceName := namePrefix + "-" + trimLength(safeSha, maxShaLen)
	resource.Name = resourceName

	if resource.Labels == nil {
		resource.Labels = map[string]string{}
	}
	resource.Labels[constants.DefaultSelectorKey] = constants.DefaultSelectorValue
	resource.Labels[launcher.RepositoryLabelKey] = safeName
	resource.Labels[launcher.CommitShaLabelKey] = safeSha

	_, err = jobInterface.Create(resource)
	if err != nil {
		return errors.Wrapf(err, "failed to create Job %s in namespace %s", resourceName, ns)
	}
	log.Logger().Infof("created Job %s in namespace %s", resourceName, ns)
	return nil
}

func trimLength(text string, length int) string {
	if len(text) <= length {
		return text
	}
	return text[0:length]
}
