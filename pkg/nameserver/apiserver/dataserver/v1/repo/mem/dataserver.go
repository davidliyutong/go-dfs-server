package mem

import (
	model "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/model"
	repoInterface "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/repo"
)

type dataserverRepo struct {
	data []*model.Dataserver
}

var _ repoInterface.DataserverRepo = (*dataserverRepo)(nil)

func (o *dataserverRepo) Create(server *model.Dataserver) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverRepo) DeleteByUUID(uuid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverRepo) DeleteByName(uuid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverRepo) Update(server *model.Dataserver) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverRepo) GetByUUID(uuid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverRepo) List() (*model.DataserverList, error) {
	//TODO implement me
	panic("implement me")
}

func newDataserverRepo() repoInterface.DataserverRepo {
	return &dataserverRepo{nil}
}
