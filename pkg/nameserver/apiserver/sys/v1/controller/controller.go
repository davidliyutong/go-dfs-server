package v1

import (
	"github.com/gin-gonic/gin"
	"go-dfs-server/pkg/nameserver/apiserver/sys/v1/repo"
	srv "go-dfs-server/pkg/nameserver/apiserver/sys/v1/service"
)

type Controller interface {
	Info(c *gin.Context)
	GetSession(c *gin.Context)
	GetSessions(c *gin.Context)
}

type controller struct {
	srv srv.Service
}

func NewController(repo repo.Repo) Controller {
	return &controller{srv: srv.NewService(repo)}
}
