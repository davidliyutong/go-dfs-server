package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ReadChunkMetaRequest struct {
	Path   string `form:"path" json:"path"`
	ID     int64  `form:"id" json:"id"`
	Offset uint64 `form:"offset" json:"offset"`
	Size   uint64 `form:"size" json:"size"`
}

type ReadChunkMetaResponse struct {
	Code     int64  `json:"code"`
	Msg      string `json:"msg"`
	Version  int64  `json:"version"`
	Checksum string `json:"checksum"`
}

func (c2 controller) ReadChunkMeta(c *gin.Context) {
	var request ReadChunkMetaRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, ReadChunkMetaResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, ReadChunkMetaResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			version, MD5String, err := c2.srv.NewBlobService().ReadChunkMeta(request.Path, request.ID)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, ReadChunkMetaResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, ReadChunkMetaResponse{Code: http.StatusOK, Msg: "", Version: version, Checksum: MD5String})
			}
		}
		log.Debug("blob/ReadChunkMeta ", request)
	}
}
