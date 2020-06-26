package job_test

import (
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher/job"
	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/batch/v1"
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
	objects, err := client.Launch(o)
	require.NoError(t, err, "failed to launch the job")
	require.Len(t, objects, 1, "should have created one runtime.Object after launching")

	o1 := objects[0]
	j1, ok := o1.(*v1.Job)
	require.True(t, ok, "could not convert object %#v to a Job")

	t.Logf("created Job with name %s", j1.Name)
}
