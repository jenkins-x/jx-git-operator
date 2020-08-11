package repo_test

import (
	"testing"

	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddGitURLUserPassword(t *testing.T) {
	testCases := []struct {
		url, username, password, expected string
	}{
		{
			url:      "https://myuser:mypwd@gitub.com/some/repo.git",
			username: "",
			password: "",
			expected: "https://myuser:mypwd@gitub.com/some/repo.git",
		},
		{
			url:      "https://gitub.com/some/repo.git",
			username: "myuser",
			password: "mypwd",
			expected: "https://myuser:mypwd@gitub.com/some/repo.git",
		},
		{
			url:      "https://myuser:mypwd@gitub.com/some/repo.git",
			username: "anotheruser",
			password: "anotherpwd",
			expected: "https://anotheruser:anotherpwd@gitub.com/some/repo.git",
		},
		{
			url:      "https://myuser:mypwd@gitub.com/some/repo.git",
			password: "anotherpwd",
			expected: "https://myuser:anotherpwd@gitub.com/some/repo.git",
		},
	}

	for _, tc := range testCases {
		actual, err := repo.AddGitURLUserPassword(tc.url, tc.username, tc.password)
		require.NoError(t, err, "should not create an error for URL %s user %s password %s", tc.url, tc.username, tc.password)
		assert.Equal(t, tc.expected, actual, "for URL %s user %s password %s", tc.url, tc.username, tc.password)

		t.Logf("URL %s user %s password %s generated: %s\n", tc.url, tc.username, tc.password, actual)
	}
}
