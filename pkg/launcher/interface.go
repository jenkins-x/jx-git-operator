package launcher

import (
	"github.com/jenkins-x/jx-git-operator/pkg/repo"
	"k8s.io/apimachinery/pkg/runtime"
)

// LaunchOptions the options for launching a new Job/Task/Pipeline for a git repository
type LaunchOptions struct {
	// Repository the repository name, namespace and git URL
	Repository repo.Repository

	// GitSHA the latest git commit SHA of the repository
	GitSHA string

	// LastCommitMessage the latest git commit message
	LastCommitMessage string

	// LastCommitAuthor the latest git commit author
	LastCommitAuthor string

	// LastCommitAuthorEmail the latest git commit author email
	LastCommitAuthorEmail string

	// LastCommitDate the latest git commit date
	LastCommitDate string

	// LastCommitURL the URL to the last git commit
	LastCommitURL string

	// 	Dir the root directory of the git clone of the repository
	Dir string

	// NoResourceApply if specified disable applying resources found in `.jx/git-operator/resources/*.yaml`
	NoResourceApply bool
}

// Interface the interface for launching Jobs/Tasks when there is a git commit in a repository
type Interface interface {
	// Launch launches a new Job/Task/Pipeline for the given repository if there has been a git commit
	// since the last
	Launch(opts LaunchOptions) ([]runtime.Object, error)
}
