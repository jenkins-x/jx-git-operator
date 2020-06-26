package launcher

import "github.com/jenkins-x/jx-git-operator/pkg/repo"

// LaunchOptions the options for launching a new Job/Task/Pipeline for a git repository
type LaunchOptions struct {
	// Repository the repository name, namespace and git URL
	Repository repo.Repository

	// GitSHA the latest git commit SHA of the repository
	GitSHA string

	// 	Dir the root directory of the git clone of the repository
	Dir string
}

// Interface the interface for launching Jobs/Tasks when there is a git commit in a repository
type Interface interface {
	// Launch lauches a new Job/Task/Pipeline for the given repository if there has been a git commit
	// since the last
	Launch(opts LaunchOptions) error
}
