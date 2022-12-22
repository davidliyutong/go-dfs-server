package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type deleteChunkRequest struct {
	Path string `form:"path" json:"path"`
	ID   int64  `form:"id" json:"id"`
}

type deleteChunkResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) DeleteChunk(c *gin.Context) {
	var request deleteChunkRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(400, deleteChunkResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, deleteChunkResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().DeleteChunk(request.Path, request.ID)
			if err != nil {
				c.IndentedJSON(500, deleteChunkResponse{Code: 500, Msg: err.Error()})
			} else {
				c.IndentedJSON(200, deleteChunkResponse{Code: 200, Msg: ""})
			}
		}
		log.Debug("blob/DeleteChunk ", request)
	}
}
