package v1

import (
	model "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/model"
	"go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/repo"
)

type DataserverService interface {
	Create(server *model.DataServer) error
	DeleteByUUID(uuid int64) error
	DeleteByName(uuid int64) error
	Update(server *model.DataServer) error
	GetByUUID(uuid int64) error
	List() (*model.DataServerList, error)
}

type dataserverService struct {
	repo repo.Repo
}

func (o *dataserverService) Create(server *model.DataServer) error {
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

func (o *dataserverService) Update(server *model.DataServer) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverService) GetByUUID(uuid int64) error {
	//TODO implement me
	panic("implement me")
}

func (o *dataserverService) List() (*model.DataServerList, error) {
	//TODO implement me
	panic("implement me")
}

func newDataserverService(repo repo.Repo) DataserverService {
	return &dataserverService{repo}
}

type Service interface {
	NewDataserverService() DataserverService
}

type service struct {
	repo repo.Repo
}

var _ Service = (*service)(nil)

func NewService(repo repo.Repo) Service {
	return &service{repo}
}

func (o *service) NewDataserverService() DataserverService {
	return newDataserverService(o.repo)
}
