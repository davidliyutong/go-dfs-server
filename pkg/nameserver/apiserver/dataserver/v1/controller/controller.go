package v1

import (
	"github.com/gin-gonic/gin"
	"go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/repo"
	srv "go-dfs-server/pkg/nameserver/apiserver/dataserver/v1/service"
)

type Controller interface {
	Create(c *gin.Context)
	DeleteByUUID(c *gin.Context)
	DeleteByName(c *gin.Context)
	Update(c *gin.Context)
	GetByUUID(c *gin.Context)
	List(c *gin.Context)
}

type controller struct {
	srv srv.Service
}

func NewController(repo repo.Repo) Controller {
	return &controller{
		srv: srv.NewService(repo),
	}
}
