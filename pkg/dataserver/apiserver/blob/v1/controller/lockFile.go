package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type lockFileRequest struct {
	Path string `form:"path" json:"path"`
}

type lockFileResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) LockFile(c *gin.Context) {
	var request lockFileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(400, lockFileResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, lockFileResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().LockFile(request.Path)
			if err != nil {
				c.IndentedJSON(500, lockFileResponse{Code: 500, Msg: err.Error()})
			} else {
				c.IndentedJSON(200, lockFileResponse{Code: 200, Msg: ""})
			}
		}
		log.Debug("blob/LockFile ", request)
	}
}
