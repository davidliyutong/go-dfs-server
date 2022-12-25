package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type lockFileRequest struct {
	Path string `form:"path" json:"path"`
	ID   string `form:"id" json:"id"`
}

type lockFileResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) LockFile(c *gin.Context) {
	var request lockFileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, lockFileResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, lockFileResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().LockFile(request.Path, request.ID)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, lockFileResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, lockFileResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/LockFile ", request)
	}
}
