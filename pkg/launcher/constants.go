package launcher

const (
	// RepositoryLabelKey the label key for associating resources (Job/Task/Pipeline) to a repository
	RepositoryLabelKey = "git-operator.jenkins.io/repository"

	// CommitShaLabelKey the label key for associating the commit sha
	CommitShaLabelKey = "git-operator.jenkins.io/commit-sha"
)
