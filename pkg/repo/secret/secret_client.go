package secret

import (
	"context"
	"fmt"

	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"github.com/jenkins-x/jx-kube-client/v3/pkg/kubeclient"

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
func NewClient(kubeClient kubernetes.Interface, ns, selector string) (repo.Interface, error) {
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
	return &client{
		kubeClient: kubeClient,
		ns:         ns,
		selector:   selector,
	}, nil
}

func (c *client) List() ([]repo.Repository, error) {
	ctx := context.Background()
	list, err := c.kubeClient.CoreV1().Secrets(c.ns).List(ctx, metav1.ListOptions{
		LabelSelector: c.selector,
	})
	if err != nil && apierrors.IsNotFound(err) {
		err = nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find Secrets in namespace %s with selector %s: %w", c.ns, c.selector, err)
	}

	var answer []repo.Repository
	for i := range list.Items {
		s := list.Items[i]
		r, err := c.toRepository(&s)
		if err != nil {
			return answer, fmt.Errorf("failed to create repo.Repository: %w", err)
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
		return repo.Repository{}, fmt.Errorf("failed to create git URL from url %s username: %s password %s: %w", rawurl, username, password, err)
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
