package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type FlushRequest struct {
	Session string `form:"session" json:"session"`
}

type FlushResponse struct {
	Code int64  `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
}

func (c2 controller) Flush(c *gin.Context) {
	var request FlushRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, FlushResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Session == "" {
			c.IndentedJSON(http.StatusBadRequest, FlushResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err := c2.srv.NewBlobService().Flush(request.Session)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, FlushResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, FlushResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/Write ", request)
	}
}
