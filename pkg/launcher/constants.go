package launcher

const (
	// RepositoryLabelKey the label key for associating resources (Job/Task/Pipeline) to a repository
	RepositoryLabelKey = "git-operator.jenkins.io/repository"

	// CommitShaLabelKey the label key for associating the commit sha
	CommitShaLabelKey = "git-operator.jenkins.io/commit-sha"

	// RerunLabelKey the label key to force this job to retrigger
	RerunLabelKey = "git-operator.jenkins.io/rerun"

	// CommitAuthorAnnotation the annotation key for the last commit author
	CommitAuthorAnnotation = "git-operator.jenkins.io/commit-author"

	// CommitAuthorEmailAnnotation the annotation key for the last commit author email
	CommitAuthorEmailAnnotation = "git-operator.jenkins.io/commit-author-email"

	// CommitDateAnnotation the annotation key for the last commit date
	CommitDateAnnotation = "git-operator.jenkins.io/commit-date"

	// CommitMessageAnnotation the annotation key for the last commit message
	CommitMessageAnnotation = "git-operator.jenkins.io/commit-message"

	// CommitURLAnnotation the annotation key for the last commit URL
	CommitURLAnnotation = "git-operator.jenkins.io/commit-url"
)
