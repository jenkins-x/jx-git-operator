package poller_test

import (
	"path/filepath"
	"testing"

	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/poller"
	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner"
	"github.com/jenkins-x/jx-helpers/pkg/cmdrunner/fakerunner"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestPoller(t *testing.T) {
	// TODO
	t.SkipNow()

	ns := "jx"
	repoName := "fake-repository"
	gitURL := "https://github.com/jenkins-x/fake-repository.git"
	gitSha := "dummysha1234"

	resourcesDir, err := filepath.Abs(filepath.Join("test_data", "somerepo", ".jx", "git-operator", "resources"))
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
		Dir:           "",
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
