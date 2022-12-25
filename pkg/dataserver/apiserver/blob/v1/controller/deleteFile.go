package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type deleteFileRequest struct {
	Path string `form:"path" json:"path"`
}

type deleteFileResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) DeleteFile(c *gin.Context) {
	var request deleteFileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, deleteFileResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, deleteFileResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().DeleteFile(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, deleteFileResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, deleteFileResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/DeleteFile ", request)
	}
}
