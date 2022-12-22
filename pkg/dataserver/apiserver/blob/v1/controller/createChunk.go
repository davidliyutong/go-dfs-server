package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type createChunkRequest struct {
	Path string `form:"path" json:"path"`
	ID   int64  `form:"id" json:"id"`
}

type createChunkResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) CreateChunk(c *gin.Context) {
	var request createChunkRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(400, createChunkResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, createChunkResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().CreateChunk(request.Path, request.ID)
			if err != nil {
				c.IndentedJSON(500, createChunkResponse{Code: 500, Msg: err.Error()})
			} else {
				c.IndentedJSON(200, createChunkResponse{Code: 200, Msg: ""})
			}
		}
		log.Debug("blob/CreateChunk ", request)
	}

}
