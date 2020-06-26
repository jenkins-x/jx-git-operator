package secret_test

import (
	"testing"

	"github.com/jenkins-x/jx-git-operator/pkg/constants"
	"github.com/jenkins-x/jx-git-operator/pkg/repo/secret"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestSecretClient(t *testing.T) {
	ns := "jx"
	secretName := "my-secret"
	gitURL := "https://github.com/jenkins-x/fake-repository.git"

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

	client, err := secret.NewClient(kubeClient, ns, constants.DefaultSelector)
	require.NoError(t, err, "failed to create repo client")

	repos, err := client.List()
	require.NoError(t, err, "failed to list repositories")

	assert.Len(t, repos, 1, "should have found 1 repo")

	r1 := repos[0]
	assert.Equal(t, secretName, r1.Name, "repo.Name")
	assert.Equal(t, ns, r1.Namespace, "repo.Namespace")
	assert.Equal(t, gitURL, r1.GitURL, "repo.GitURL")

	t.Logf("found Repository %s in namespace %s with git URL %s", r1.Name, r1.Namespace, r1.GitURL)
}
