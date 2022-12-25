package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type deleteDirectoryRequest struct {
	Path string `form:"path" json:"path"`
}

type deleteDirectoryResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) DeleteDirectory(c *gin.Context) {
	var request deleteDirectoryRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, deleteDirectoryResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, deleteDirectoryResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().DeleteDirectory(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, deleteDirectoryResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, deleteDirectoryResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/DeleteDirectory ", request)
	}
}
