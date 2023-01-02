package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/dataserver/server"
	"net/http"
)

type InfoResponse struct {
	Code    int16  `json:"code"`
	Msg     string `json:"msg"`
	Role    string `json:"role"`
	Version string `json:"version"`
}

func (o *controller) Info(c *gin.Context) {
	c.JSON(http.StatusOK, InfoResponse{
		Code:    200,
		Msg:     "",
		Role:    "dataserver",
		Version: server.DataServerAPIVersion,
	})
	log.Debugln("sys/Info ")
}
