package repo

type Interface interface {
	List() ([]Repository, error)
}
