package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
		c.IndentedJSON(400, deleteFileResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, deleteFileResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().DeleteFile(request.Path)
			if err != nil {
				c.IndentedJSON(500, deleteFileResponse{Code: 500, Msg: err.Error()})
			} else {
				c.IndentedJSON(200, deleteFileResponse{Code: 200, Msg: ""})
			}
		}
		log.Debug("blob/DeleteFile ", request)
	}
}
