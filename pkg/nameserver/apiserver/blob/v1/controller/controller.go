package v1

import (
	"github.com/gin-gonic/gin"
	"go-dfs-server/pkg/nameserver/apiserver/blob/v1/repo"
	srv "go-dfs-server/pkg/nameserver/apiserver/blob/v1/service"
)

type Controller interface {
	Sync(c *gin.Context)
	Ls(c *gin.Context)
	Mkdir(c *gin.Context)
	Open(c *gin.Context)
	Read(c *gin.Context)
	Rm(c *gin.Context)
	Write(c *gin.Context)
}

type controller struct {
	srv srv.Service
}

func NewController(repo repo.Repo, err error) Controller {
	return &controller{
		srv: srv.NewService(repo),
	}
}
