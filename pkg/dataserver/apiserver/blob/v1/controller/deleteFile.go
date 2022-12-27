package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type DeleteFileRequest struct {
	Path string `form:"path" json:"path"`
}

type DeleteFileResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) DeleteFile(c *gin.Context) {
	var request DeleteFileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, DeleteFileResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, DeleteFileResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().DeleteFile(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, DeleteFileResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, DeleteFileResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/DeleteFile ", request)
	}
}
