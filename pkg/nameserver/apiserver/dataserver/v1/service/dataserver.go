package v1

import (
	model "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/model"
	"go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/repo"
)

type DataserverService interface {
	Create(server *model.Dataserver) error
	DeleteByUUID(uuid int64) error
	DeleteByName(uuid int64) error
	Update(server *model.Dataserver) error
	GetByUUID(uuid int64) error
	List() (*model.DataserverList, error)
}

type dataserverService struct {
	repo repo.Repo
}

func (o *dataserverService) Create(server *model.Dataserver) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverService) DeleteByUUID(uuid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverService) DeleteByName(uuid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverService) Update(server *model.Dataserver) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverService) GetByUUID(uuid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverService) List() (*model.DataserverList, error) {
	//TODO implement me
	panic("implement me")
}

func newDataserverService(repo repo.Repo) DataserverService {
	return &dataserverService{repo}
}
