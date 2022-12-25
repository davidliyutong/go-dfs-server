package v1

import (
	"github.com/gin-gonic/gin"
	srv "go-dfs-server/pkg/dataserver/apiserver/sys/v1/service"
	"go-dfs-server/pkg/nameserver/apiserver/sys/v1/repo"
)

type Controller interface {
	Info(c *gin.Context)
}

type controller struct {
	srv srv.Service
}

func NewController(repo repo.Repo) Controller {
	return &controller{srv: srv.NewService(repo)}
}
