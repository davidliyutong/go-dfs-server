package mem

import (
	model "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/model"
	repoInterface "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/repo"
)

type dataServerRepo struct {
	data []*model.DataServer
}

var _ repoInterface.DataServerRepo = (*dataServerRepo)(nil)

func (o *dataServerRepo) Create(server *model.DataServer) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataServerRepo) DeleteByUUID(uuid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataServerRepo) DeleteByName(uuid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataServerRepo) Update(server *model.DataServer) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataServerRepo) GetByUUID(uuid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataServerRepo) List() (*model.DataServerList, error) {
	//TODO implement me
	panic("implement me")
}

func newDataServerRepo() repoInterface.DataServerRepo {
	return &dataServerRepo{nil}
}
