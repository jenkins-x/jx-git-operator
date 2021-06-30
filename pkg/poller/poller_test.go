package poller_test

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/launcher"
	"github.com/jenkins-x/jx-git-operator/pkg/poller"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/cmdrunner/fakerunner"
	"github.com/jenkins-x/jx-helpers/v3/pkg/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestPoller(t *testing.T) {
	if os.Getenv("CI") == "true" {
		t.Skip("Skipping testing in CI environment")
	}
	ns := "jx"
	ctx := context.Background()
	repoName := "fake-repository"
	gitURL := "https://github.com/jenkins-x/fake-repository.git"
	gitSha := "dummysha1234"

	resourcesDir, err := filepath.Abs(filepath.Join("test_data", "somerepo", ".jx", "git-operator", "resources"))
	require.NoError(t, err, "failed to get absolute dir %s", resourcesDir)

	tmpDir, err := ioutil.TempDir("", "test-jx-git-operator-")
	require.NoError(t, err, "failed to create temp dir")

	t.Logf("running in dir %s", tmpDir)

	// lets copy the dummy git clone to the temp dir
	err = files.CopyDirOverwrite(filepath.Join("test_data", repoName), filepath.Join(tmpDir, repoName))
	require.NoError(t, err, "failed to copy git clone data to temp dir")

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
	runner := &fakerunner.FakeRunner{
		CommandRunner: func(c *cmdrunner.Command) (string, error) {
			if c.Name == "git" && len(c.Args) > 0 && c.Args[0] == "rev-parse" {
				return gitSha, nil
			}
			return "", nil
		},
	}

	p := &poller.Options{
		CommandRunner: runner.Run,
		KubeClient:    kubeClient,
		Dir:           tmpDir,
		Namespace:     ns,
		NoLoop:        true,
	}

	err = p.Run()
	require.NoError(t, err, "failed to run poller")

	assertHasJobCountForRepoAndSha(t, ctx, kubeClient, ns, repoName, gitSha, 1)

	err = p.Run()
	require.NoError(t, err, "failed to run poller")

	assertHasJobCountForRepoAndSha(t, ctx, kubeClient, ns, repoName, gitSha, 1)

	// now lets do a new commit
	firstGitSha := gitSha
	gitSha = "new-commit-sha"
	t.Logf("now creating a second commit with sha %s", gitSha)

	err = p.Run()
	require.NoError(t, err, "failed to run poller")

	oldJobs := assertHasJobCountForRepoAndSha(t, ctx, kubeClient, ns, repoName, firstGitSha, 1)
	assertHasJobCountForRepoAndSha(t, ctx, kubeClient, ns, repoName, gitSha, 0)

	// now lets make the first job as completed
	require.Len(t, oldJobs, 1, "should have one job for the old git commit")

	job := oldJobs[0]
	job.Status.Succeeded = 1
	_, err = kubeClient.BatchV1().Jobs(ns).Update(ctx, &job, metav1.UpdateOptions{})
	require.NoError(t, err, "failed to update the job %s in namespace %s to succeeded", job.Name, ns)

	err = p.Run()
	require.NoError(t, err, "failed to run poller")

	assertHasJobCountForRepoAndSha(t, ctx, kubeClient, ns, repoName, gitSha, 1)

	for _, c := range runner.OrderedCommands {
		t.Logf("created command: %s\n", c.CLI())
	}
}

func assertHasJobCountForRepoAndSha(t *testing.T, ctx context.Context, kubeClient kubernetes.Interface, ns string, repoName string, sha string, expectedCount int) []v1.Job {
	selector := constants.DefaultSelectorKey
	jobs, err := kubeClient.BatchV1().Jobs(ns).List(ctx, metav1.ListOptions{
		LabelSelector: selector,
	})
	require.NoError(t, err, "failed to list jobs in namespace %s with selector %s", ns, selector)
	require.NotNil(t, jobs, "no JobsList object returned for namespace %s with selector %s", ns, selector)

	count := 0
	for _, j := range jobs.Items {
		name := j.Name
		labels := j.Labels
		if labels == nil {
			t.Logf("Job %s in namespace %s does not have any labels", name, ns)
			continue
		}
		if labels[launcher.RepositoryLabelKey] == repoName && labels[launcher.CommitShaLabelKey] == sha {
			count++
			t.Logf("found Job %s in namespace %s for repoName %s and sha %s", name, ns, repoName, sha)
		}
	}
	assert.Equal(t, expectedCount, count, "number of Jobs in namespace %s with selector %s with repo %s and git sha %s", ns, selector, repoName, sha)
	return jobs.Items
}

func TestLazyCreatePoller(t *testing.T) {
	p := &poller.Options{}

	err := p.ValidateOptions()
	require.NoError(t, err, "failed to ValidateOptions()")

	assert.NotNil(t, p.GitClient, "GitClient")
	assert.NotNil(t, p.Launcher, "Launcher")
}
