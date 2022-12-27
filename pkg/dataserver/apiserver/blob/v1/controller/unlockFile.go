package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UnlockFileRequest struct {
	Path string `form:"path" json:"path"`
}

type UnlockFileResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) UnlockFile(c *gin.Context) {
	var request UnlockFileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, UnlockFileResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, UnlockFileResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().UnlockFile(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, UnlockFileResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, UnlockFileResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/UnlockFile", request)
	}
}
