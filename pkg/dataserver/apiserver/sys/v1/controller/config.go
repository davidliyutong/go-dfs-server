package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/dataserver/server"
	"net/http"
)

type ConfigResponse struct {
	Code     int16  `json:"code"`
	Msg      string `json:"msg"`
	Port     int64  `json:"port"`
	Endpoint string `json:"endpoint"`
	Volume   string `json:"volume"`
}

func (o *controller) Config(c *gin.Context) {
	c.JSON(http.StatusOK, ConfigResponse{
		Code:     200,
		Msg:      "",
		Port:     server.GlobalServerDesc.Opt.Network.Port,
		Endpoint: server.GlobalServerDesc.Opt.Network.Endpoint,
		Volume:   server.GlobalServerDesc.Opt.Volume,
	})
	log.Debugln("sys/Config ")
}
