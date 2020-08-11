package secret

import (
	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"github.com/jenkins-x/jx-kube-client/pkg/kubeclient"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type client struct {
	kubeClient kubernetes.Interface
	ns         string
	selector   string
}

// NewClient creates a new client using the given kubernetes client and namespace
// if nil is passed in the kubernetes client will be lazily created
func NewClient(kubeClient kubernetes.Interface, ns string, selector string) (repo.Interface, error) {
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

func (c *client) List() ([]repo.Repository, error) {
	list, err := c.kubeClient.CoreV1().Secrets(c.ns).List(metav1.ListOptions{
		LabelSelector: c.selector,
	})
	if err != nil && apierrors.IsNotFound(err) {
		err = nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find Secrets in namespace %s with selector %s", c.ns, c.selector)
	}

	var answer []repo.Repository
	for _, s := range list.Items {
		r, err := c.toRepository(&s)
		if err != nil {
			return answer, errors.Wrapf(err, "failed to create repo.Repository")
		}
		if r.GitURL != "" {
			answer = append(answer, r)
		}
	}
	return answer, nil
}

func (c *client) toRepository(s *v1.Secret) (repo.Repository, error) {
	if s.Data == nil {
		s.Data = map[string][]byte{}
	}

	rawurl := string(s.Data["url"])
	username := string(s.Data["username"])
	password := string(s.Data["password"])
	gitURL, err := repo.AddGitURLUserPassword(rawurl, username, password)
	if err != nil {
		return repo.Repository{}, errors.Wrapf(err, "failed to create git URL from url %s username: %s password %s", rawurl, username, password)
	}
	ns := s.Namespace
	if ns == "" {
		ns = c.ns
	}
	return repo.Repository{
		Name:      s.Name,
		Namespace: ns,
		GitURL:    gitURL,
	}, nil
}
