package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UnlockRequest struct {
	Session string `form:"session" json:"session"`
}

type UnlockResponse struct {
	Code int64  `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
}

func (c2 controller) Unlock(c *gin.Context) {
	var request UnlockRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, UnlockResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Session == "" {
			c.IndentedJSON(http.StatusBadRequest, UnlockResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().Unlock(request.Session)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, UnlockResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, UnlockResponse{Code: http.StatusOK, Msg: ""})
			}
		}
	}
	log.Debug("blob/Unlock ", request)
}
