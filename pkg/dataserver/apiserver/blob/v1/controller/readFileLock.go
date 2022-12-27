package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ReadFileLockRequest struct {
	Path string `form:"path" json:"path"`
}

type ReadFileLockResponse struct {
	Code     int64    `json:"code"`
	Msg      string   `json:"msg"`
	Sessions []string `json:"id"`
}

func (c2 controller) ReadFileLock(c *gin.Context) {
	var request ReadFileLockRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, ReadFileLockResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, ReadFileLockResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			locks, err := c2.srv.NewBlobService().ReadFileLock(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, ReadFileLockResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, ReadFileLockResponse{Code: http.StatusOK, Msg: "", Sessions: locks})
			}
		}
		log.Debug("blob/readFileLock ", request)
	}
}
