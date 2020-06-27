package poller_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/poller"
	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner/fakerunner"
	"github.com/jenkins-x/jx-helpers/pkg/files"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestPoller(t *testing.T) {
	ns := "jx"
	repoName := "fake-repository"
	gitURL := "https://github.com/jenkins-x/fake-repository.git"
	gitSha := "dummysha1234"

	resourcesDir, err := filepath.Abs(filepath.Join("test_data", "somerepo", ".jx", "git-operator", "resources"))
	require.NoError(t, err, "failed to get absolute dir %s", resourcesDir)

	tmpDir, err := ioutil.TempDir("", "test-jx-git-operator-")
	require.NoError(t, err, "failed to create temp dir")

	t.Logf("running in dir %s", tmpDir)

	// lets copy the dummy git clone to the temp dir
	err = files.CopyDirOverwrite(filepath.Join("test_data",  repoName), filepath.Join(tmpDir, repoName))
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
		CommandRunner: func (c *cmdrunner.Command) (string, error) {
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


	// TODO verify we have a Job


	err = p.Run()
	require.NoError(t, err, "failed to run poller")


	// TODO verify we don't create another one
	for _, c := range runner.OrderedCommands {
		t.Logf("created command: %s\n", c.CLI())
	}
	/*
	runner.ExpectResults(t,
		fakerunner.FakeResult{
			CLI: "kubectl apply -f " + resourcesDir,
		},
	)
	*/
}
