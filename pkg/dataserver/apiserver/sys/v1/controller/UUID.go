package v1

import (
	"github.com/gin-gonic/gin"
	"go-dfs-server/pkg/dataserver/server"
	"net/http"
)

type UUIDResponse struct {
	Code int16  `json:"code"`
	Msg  string `json:"msg"`
	UUID string `json:"uuid"`
}

func (o *controller) UUID(c *gin.Context) {
	c.JSON(http.StatusOK, UUIDResponse{
		Code: 200,
		Msg:  "",
		UUID: server.GlobalServerDesc.Opt.UUID,
	})
}
