package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type WriteChunkRequest struct {
	Path    string `form:"path" json:"path"`
	ID      int64  `form:"id" json:"id"`
	Version int64  `form:"version" json:"version"`
}

type WriteChunkResponse struct {
	Code     int64  `json:"code"`
	Msg      string `json:"msg"`
	Checksum string `json:"md5"`
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
		MD5String, err := c2.srv.NewBlobService().WriteChunk(request.Path, request.ID, request.Version, c, file)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, WriteChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			return
		}
		c.IndentedJSON(http.StatusOK, WriteChunkResponse{Code: http.StatusOK, Msg: "", Checksum: MD5String})

		log.Debug("blob/WriteChunk ", request)
	}

}
