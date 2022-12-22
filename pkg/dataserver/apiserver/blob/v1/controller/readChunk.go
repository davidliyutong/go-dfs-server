package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type readChunkRequest struct {
	Path string `form:"path" json:"path"`
	ID   int64  `form:"id" json:"id"`
}

type readChunkResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) ReadChunk(c *gin.Context) {
	var request readChunkRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(400, readChunkResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, readChunkResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().ReadChunk(request.Path, request.ID, c)
			if err != nil {
				c.IndentedJSON(500, writeChunkResponse{Code: 500, Msg: err.Error()})
			}
		}
		log.Debug("blob/writeChunk ", request)
	}

}
