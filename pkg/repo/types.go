package repo

// Repository represents a git repository to clone
type Repository struct {
	// Name name of the repository
	Name string

	// Namespace of the repository - where the `Job` will be created
	Namespace string

	// GitURL the URL to git clone the repository
	GitURL string
}
