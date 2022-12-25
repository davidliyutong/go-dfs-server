package repo

import model "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/model"

type DataServerRepo interface {
	Create(server *model.DataServer) error
	DeleteByUUID(uuid int64) error
	DeleteByName(uuid int64) error
	Update(server *model.DataServer) error
	GetByUUID(uuid int64) error
	List() (*model.DataServerList, error)
}
