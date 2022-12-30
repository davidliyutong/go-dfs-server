package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type MkdirRequest struct {
	Path string `form:"path" json:"path"`
}

type MkdirResponse struct {
	Code int64  `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
}

func (c2 controller) Mkdir(c *gin.Context) {
	var request MkdirRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, MkdirResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, MkdirResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().Mkdir(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, MkdirResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, MkdirResponse{Code: http.StatusOK, Msg: ""})
			}
		}
	}
	log.Debug("blob/Mkdir ", request)
}
