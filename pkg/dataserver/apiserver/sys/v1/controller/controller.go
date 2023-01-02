package v1

import (
	"github.com/gin-gonic/gin"
	"go-dfs-server/pkg/dataserver/apiserver/sys/v1/repo"
	srv "go-dfs-server/pkg/dataserver/apiserver/sys/v1/service"
)

type Controller interface {
	Info(c *gin.Context)
	Config(c *gin.Context)
	UUID(c *gin.Context)
	Register(c *gin.Context)
}

type controller struct {
	srv srv.Service
}

var _ Controller = &controller{}

func NewController(repo repo.Repo) Controller {
	return &controller{
		srv: srv.NewService(repo),
	}
}
