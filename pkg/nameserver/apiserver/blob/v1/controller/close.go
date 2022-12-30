package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type CloseRequest struct {
	Session string `form:"session" json:"session"`
}

type CloseResponse struct {
	Code int64  `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
}

func (c2 controller) Close(c *gin.Context) {
	var request CloseRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, CloseResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Session == "" {
			c.IndentedJSON(http.StatusBadRequest, CloseResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err := c2.srv.NewBlobService().Close(request.Session)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, CloseResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, CloseResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/Close ", request)
	}
}
