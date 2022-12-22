package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
		c.IndentedJSON(400, deleteDirectoryResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, deleteDirectoryResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().DeleteDirectory(request.Path)
			if err != nil {
				c.IndentedJSON(500, deleteDirectoryResponse{Code: 500, Msg: err.Error()})
			} else {
				c.IndentedJSON(200, deleteDirectoryResponse{Code: 200, Msg: ""})
			}
		}
		log.Debug("blob/DeleteDirectory ", request)
	}
}
