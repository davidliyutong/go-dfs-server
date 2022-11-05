package repo

import model "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/model"

type DataserverRepo interface {
	Create(server *model.Dataserver) error
	DeleteByUUID(uuid int64) error
	DeleteByName(uuid int64) error
	Update(server *model.Dataserver) error
	GetByUUID(uuid int64) error
	List() (*model.DataserverList, error)
}
