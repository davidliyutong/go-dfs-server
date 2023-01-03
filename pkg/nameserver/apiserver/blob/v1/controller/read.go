package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ReadRequest struct {
	Path        string `form:"path" json:"path"`
	ChunkID     int64  `form:"chunk_id" json:"chunk_id"`
	ChunkOffset int64  `form:"chunk_offset" json:"chunk_offset"`
	Size        int64  `form:"size" json:"size"`
}

type ReadResponse struct {
	Code int64  `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
}

func (c2 controller) Read(c *gin.Context) {
	var request ReadRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, ReadResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, ReadResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			_, err := c2.srv.NewBlobService().Read(request.Path, request.ChunkID, request.ChunkOffset, request.Size, c)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, ReadResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			}
		}
	}
	log.Debug("blob/Read ", request)
}
