package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RmdirRequest struct {
	Path string `form:"path" json:"path"`
}

type RmdirResponse struct {
	Code int64  `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
}

func (c2 controller) Rmdir(c *gin.Context) {
	var request RmdirRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, RmdirResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, RmdirResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().Rmdir(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, RmdirResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, RmdirResponse{Code: http.StatusOK, Msg: ""})
			}
		}
	}
	log.Debug("blob/Rmdir ", request)
}
