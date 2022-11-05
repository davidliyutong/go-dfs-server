package repo

type Repo interface {
	DataserverRepo() DataserverRepo
	Close() error
}
