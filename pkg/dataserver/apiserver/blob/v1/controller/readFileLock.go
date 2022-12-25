package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type readFileLockRequest struct {
	Path string `form:"path" json:"path"`
}

type readFileLockResponse struct {
	Code int64    `json:"code"`
	Msg  string   `json:"msg"`
	ID   []string `json:"id"`
}

func (c2 controller) ReadFileLock(c *gin.Context) {
	var request readFileLockRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, readFileLockResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, readFileLockResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			locks, err := c2.srv.NewBlobService().ReadFileLock(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, readFileLockResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, readFileLockResponse{Code: http.StatusOK, Msg: "", ID: locks})
			}
		}
		log.Debug("blob/readFileMeta ", request)
	}
}
