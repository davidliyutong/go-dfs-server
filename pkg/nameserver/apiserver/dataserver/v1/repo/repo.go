package repo

type Repo interface {
	DataServerRepo() DataServerRepo
	Close() error
}
