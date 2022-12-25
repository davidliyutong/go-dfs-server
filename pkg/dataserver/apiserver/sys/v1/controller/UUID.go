package v1

import (
	"github.com/gin-gonic/gin"
	"go-dfs-server/pkg/dataserver/server"
	"net/http"
)

func (o *controller) UUID(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "",
		"uuid": server.GlobalServerDesc.Opt.UUID,
	})
}
