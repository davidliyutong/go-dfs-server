package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type GetLockRequest struct {
	Path string `form:"path" json:"path"`
}

type GetLockResponse struct {
	Code     int64    `form:"code" json:"code"`
	Msg      string   `form:"msg" json:"msg"`
	Sessions []string `form:"locks" json:"locks"`
}

func (c2 controller) GetLock(c *gin.Context) {
	var request GetLockRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, GetLockResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, GetLockResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			sessions, err := c2.srv.NewBlobService().GetLock(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, GetLockResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, GetLockResponse{Code: http.StatusOK, Msg: "", Sessions: sessions})
			}
		}
	}
	log.Debug("blob/GetLock ", request)
}
