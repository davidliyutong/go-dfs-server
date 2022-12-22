package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type createFileRequest struct {
	Path string `form:"path" json:"path"`
}

type createFileResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) CreateFile(c *gin.Context) {
	var request createFileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(400, createFileResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, createFileResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().CreateFile(request.Path)
			if err != nil {
				c.IndentedJSON(500, createFileResponse{Code: 500, Msg: err.Error()})
			} else {
				c.IndentedJSON(200, createFileResponse{Code: 200, Msg: ""})
			}
		}
	}
	log.Debug("blob/CreateFile ", request)

}
