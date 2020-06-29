package repo

// Interface the interface for querying the git repositories to be operated
type Interface interface {
	// List lists the repositories enabled for the git operator
	List() ([]Repository, error)
}
