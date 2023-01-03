package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type WriteChunkRequest struct {
	Path    string `form:"path" json:"path"`
	ID      int64  `form:"id" json:"id"`
	Offset  int64  `form:"offset" json:"offset"`
	Size    int64  `form:"size" json:"size"`
	Version int64  `form:"version" json:"version"`
}

type WriteChunkResponse struct {
	Code     int64  `form:"code" json:"code"`
	Msg      string `form:"msg" json:"msg"`
	Checksum string `form:"md5" json:"md5"`
	Written  int64  `form:"written" json:"written"`
}

func (c2 controller) WriteChunk(c *gin.Context) {
	var request WriteChunkRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, WriteChunkResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, WriteChunkResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, WriteChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			return
		}
		MD5String, written, err := c2.srv.NewBlobService().WriteChunk(request.Path, request.ID, request.Offset, request.Size, request.Version, file)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, WriteChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, WriteChunkResponse{Code: http.StatusOK, Msg: "", Checksum: MD5String, Written: written})

		log.Debug("blob/WriteChunk ", request)
	}

}
