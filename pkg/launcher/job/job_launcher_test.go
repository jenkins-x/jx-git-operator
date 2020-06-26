package job_test

import (
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher/job"
	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestJobLauncher(t *testing.T) {
	ns := "jx"
	secretName := "my-secret"
	gitURL := "https://github.com/jenkins-x/fake-repository.git"
	gitSha := "dummysha1234"

	kubeClient := fake.NewSimpleClientset(
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: ns,
				Labels: map[string]string{
					constants.DefaultSelectorKey: constants.DefaultSelectorValue,
				},
			},
			Data: map[string][]byte{
				"url": []byte(gitURL),
			},
		},
	)

	client, err := job.NewLauncher(kubeClient, ns, constants.DefaultSelector)
	require.NoError(t, err, "failed to create laucher client")

	o := launcher.LaunchOptions{
		Repository: repo.Repository{
			Name:      secretName,
			Namespace: ns,
			GitURL:    gitURL,
		},
		GitSHA: gitSha,
		Dir:    filepath.Join("test_data", "somerepo"),
	}
	err = client.Launch(o)
	require.NoError(t, err, "failed to launch the job")
}
