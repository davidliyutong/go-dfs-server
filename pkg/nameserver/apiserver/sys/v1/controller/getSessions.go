package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type GetSessionsRequest struct {
}

type GetSessionsResponse struct {
	Code     int64    `form:"code" json:"code"`
	Msg      string   `form:"msg" json:"msg"`
	Sessions []string `form:"sessions" json:"sessions"`
}

func (o *controller) GetSessions(c *gin.Context) {

	sessions, err := o.srv.NewSysService().GetSessions()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, GetSessionsResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, GetSessionsResponse{Code: http.StatusOK, Msg: "", Sessions: sessions})
	}
	log.Debug("blob/GetSessions ")
}
