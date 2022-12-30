package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type LockRequest struct {
	Session string `form:"session" json:"session"`
}

type LockResponse struct {
	Code int64  `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
}

func (c2 controller) Lock(c *gin.Context) {
	var request LockRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, LockResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Session == "" {
			c.IndentedJSON(http.StatusBadRequest, LockResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err := c2.srv.NewBlobService().Lock(request.Session)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, LockResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, LockResponse{Code: http.StatusOK, Msg: ""})
			}
		}
	}
	log.Debug("blob/Lock ", request)
}
