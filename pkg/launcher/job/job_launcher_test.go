package job_test

import (
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher/job"
	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner/fakerunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/testhelpers"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestJobLauncher(t *testing.T) {
	ns := "jx"
	repoName := "fake-repository"
	gitURL := "https://github.com/jenkins-x/fake-repository.git"
	gitSha := "dummysha1234"

	resourcesDir, err := filepath.Abs(filepath.Join("test_data", "somerepo", "versionStream", "git-operator", "resources"))
	require.NoError(t, err, "failed to get absolute dir %s", resourcesDir)

	kubeClient := fake.NewSimpleClientset(
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      repoName,
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
	runner := &fakerunner.FakeRunner{}

	client, err := job.NewLauncher(kubeClient, ns, constants.DefaultSelector, runner.Run)
	require.NoError(t, err, "failed to create launcher client")

	o := launcher.LaunchOptions{
		Repository: repo.Repository{
			Name:      repoName,
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

	msg := "created Job"
	testhelpers.AssertLabel(t, constants.DefaultSelectorKey, constants.DefaultSelectorValue, j1.ObjectMeta, msg)
	testhelpers.AssertLabel(t, launcher.RepositoryLabelKey, repoName, j1.ObjectMeta, msg)
	testhelpers.AssertLabel(t, launcher.CommitShaLabelKey, gitSha, j1.ObjectMeta, msg)

	runner.ExpectResults(t,
		fakerunner.FakeResult{
			CLI: "kubectl apply -f " + resourcesDir,
		},
	)

	// we should not recreate the Job if we try to launch again as it already exists
	objects, err = client.Launch(o)
	require.NoError(t, err, "failed to launch the job")
	require.Len(t, objects, 0, "should not have a created a runtime.Object as we already have one for the commit sha")
}
