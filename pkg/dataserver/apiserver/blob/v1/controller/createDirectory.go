package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type createDirectoryRequest struct {
	Path string `form:"path" json:"path"`
}

type createDirectoryResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) CreateDirectory(c *gin.Context) {
	var request createDirectoryRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(400, createDirectoryResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, createDirectoryResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().CreateDirectory(request.Path)
			if err != nil {
				c.IndentedJSON(500, createDirectoryResponse{Code: 500, Msg: err.Error()})
			} else {
				c.IndentedJSON(200, createDirectoryResponse{Code: 200, Msg: ""})
			}
		}
		log.Debug("blob/CreateDirectory ", request)
	}

}
