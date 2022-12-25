package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
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
		c.IndentedJSON(http.StatusBadRequest, readChunkMetaResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, readChunkMetaResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			MD5String, err := c2.srv.NewBlobService().ReadChunkMeta(request.Path, request.ID)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, writeChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, writeChunkResponse{Code: http.StatusOK, Msg: "", MD5: MD5String})
			}
		}
		log.Debug("blob/readChunkMeta ", request)
	}
}
