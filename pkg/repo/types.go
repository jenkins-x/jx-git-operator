package repo

// Repository represents a git repository to clone
type Repository struct {
	Name      string
	Namespace string
	GitURL    string
}
