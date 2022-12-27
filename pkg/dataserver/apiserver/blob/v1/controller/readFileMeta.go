package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type ReadFileMetaRequest struct {
	Path string `form:"path" json:"path"`
}

type ReadFileMetaResponse struct {
	Code      int64            `json:"code"`
	Msg       string           `json:"msg"`
	Checksums map[int64]string `json:"checksums"`
}

func (c2 controller) ReadFileMeta(c *gin.Context) {
	var request ReadFileMetaRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, ReadFileMetaResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, ReadFileMetaResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			meta, err := c2.srv.NewBlobService().ReadFileMeta(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, ReadFileMetaResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, ReadFileMetaResponse{Code: http.StatusOK, Msg: "", Checksums: meta.Content})
			}
		}
		log.Debug("blob/ReadFileMeta ", request)
	}
}
