package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/config"
	"net/http"
)

type GetServersRequest struct {
}

type GetServersResponse struct {
	Code    int64                         `form:"code" json:"code"`
	Msg     string                        `form:"msg" json:"msg"`
	Servers []config.RegisteredDataServer `form:"servers" json:"servers"`
}

func (o *controller) GetServers(c *gin.Context) {
	var request GetServersRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, GetServersResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {

		ds, err := o.srv.NewSysService().GetServers()
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, GetServersResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
		} else {
			c.IndentedJSON(http.StatusOK, GetServersResponse{
				Code:    http.StatusOK,
				Msg:     "",
				Servers: ds})
		}

	}
	log.Debug("blob/GetSession ", request)
}
