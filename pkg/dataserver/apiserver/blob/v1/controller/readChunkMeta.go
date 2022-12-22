package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type readChunkMetaRequest struct {
	Path string `form:"path" json:"path"`
	ID   int64  `form:"id" json:"id"`
}

type readChunkMetaResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	MD5  string `json:"md5"`
}

func (c2 controller) ReadChunkMeta(c *gin.Context) {
	var request readChunkMetaRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(400, readChunkMetaResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, readChunkMetaResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			MD5String, err := c2.srv.NewBlobService().ReadChunkMeta(request.Path, request.ID)
			if err != nil {
				c.IndentedJSON(500, writeChunkResponse{Code: 500, Msg: err.Error()})
			} else {
				c.IndentedJSON(200, writeChunkResponse{Code: 200, Msg: "", MD5: MD5String})
			}
		}
		log.Debug("blob/readChunkMeta ", request)
	}
}
