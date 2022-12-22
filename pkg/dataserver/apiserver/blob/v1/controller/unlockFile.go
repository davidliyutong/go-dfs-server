package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type unlockFileRequest struct {
	Path string `form:"path" json:"path"`
}

type unlockFileResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) UnlockFile(c *gin.Context) {
	var request unlockFileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(400, unlockFileResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, unlockFileResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().UnlockFile(request.Path)
			if err != nil {
				c.IndentedJSON(500, unlockFileResponse{Code: 500, Msg: err.Error()})
			} else {
				c.IndentedJSON(200, unlockFileResponse{Code: 200, Msg: ""})
			}
		}
		log.Debug("blob/LockFile", request)
	}
}
