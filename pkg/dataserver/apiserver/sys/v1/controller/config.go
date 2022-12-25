package v1

import (
	"github.com/gin-gonic/gin"
	"go-dfs-server/pkg/dataserver/server"
	"net/http"
)

func (o *controller) Config(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"msg":      "",
		"port":     server.GlobalServerDesc.Opt.Network.Port,
		"endpoint": server.GlobalServerDesc.Opt.Network.Endpoint,
		"volume":   server.GlobalServerDesc.Opt.Volume,
	})
}
