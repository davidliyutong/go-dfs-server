package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ReadChunkRequest struct {
	Path   string `form:"path" json:"path"`
	ID     int64  `form:"id" json:"id"`
	Offset int64  `form:"offset" json:"offset"`
	Size   int64  `form:"size" json:"size"`
}

type ReadChunkResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) ReadChunk(c *gin.Context) {
	var request ReadChunkRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, ReadChunkResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, ReadChunkResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().ReadChunk(request.Path, request.ID, request.Offset, request.Size, c)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, ReadChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			}
		}
		log.Debug("blob/WriteChunk ", request)
	}

}
