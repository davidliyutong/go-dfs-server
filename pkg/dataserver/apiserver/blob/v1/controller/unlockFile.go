package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type unlockFileRequest struct {
	Path string `form:"path" json:"path"`
}

type unlockFileResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) UnlockFile(c *gin.Context) {
	var request unlockFileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, unlockFileResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, unlockFileResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().UnlockFile(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, unlockFileResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, unlockFileResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/LockFile", request)
	}
}
