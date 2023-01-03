package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type WriteRequest struct {
	Path        string `form:"path" json:"path"`
	ChunkID     int64  `form:"chunk_id" json:"chunk_id"`
	ChunkOffset int64  `form:"chunk_offset" json:"chunk_offset"`
	Size        int64  `form:"size" json:"size"`
	Version     int64  `form:"version" json:"version"`
}

type WriteResponse struct {
	Code     int64  `form:"code" json:"code"`
	Msg      string `form:"msg" json:"msg"`
	Checksum string `form:"md5" json:"md5"`
	Written  int    `form:"written" json:"written"`
}

func (c2 controller) Write(c *gin.Context) {
	var request WriteRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, WriteResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, WriteResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			file, err := c.FormFile("file")
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, WriteResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				checksum, size, err := c2.srv.NewBlobService().Write(request.Path, request.ChunkID, request.ChunkOffset, request.Size, request.Version, file)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, WriteResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
				} else {
					c.IndentedJSON(http.StatusOK, WriteResponse{Code: http.StatusOK, Msg: "", Checksum: checksum, Written: int(size)})
				}
			}
		}
		log.Debug("blob/Write ", request)
	}
}
