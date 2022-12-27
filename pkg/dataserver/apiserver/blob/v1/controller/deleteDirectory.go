package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type DeleteDirectoryRequest struct {
	Path string `form:"path" json:"path"`
}

type DeleteDirectoryResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) DeleteDirectory(c *gin.Context) {
	var request DeleteDirectoryRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, DeleteDirectoryResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, DeleteDirectoryResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().DeleteDirectory(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, DeleteDirectoryResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, DeleteDirectoryResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/DeleteDirectory ", request)
	}
}
