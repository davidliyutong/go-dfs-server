package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type LockFileRequest struct {
	Path    string `form:"path" json:"path"`
	Session string `form:"id" json:"id"`
}

type LockFileResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) LockFile(c *gin.Context) {
	var request LockFileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, LockFileResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, LockFileResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().LockFile(request.Path, request.Session)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, LockFileResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, LockFileResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/LockFile ", request)
	}
}
