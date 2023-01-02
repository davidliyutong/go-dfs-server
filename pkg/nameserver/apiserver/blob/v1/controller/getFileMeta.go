package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"net/http"
)

type GetFileMetaRequest struct {
	Path string `form:"path" json:"path"`
}

type GetFileMetaResponse struct {
	Code int64           `json:"code"`
	Msg  string          `json:"msg"`
	Blob v1.BlobMetaData `json:"blob"`
}

func (o *controller) GetFileMeta(c *gin.Context) {
	var request GetFileMetaRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, GetFileMetaResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, GetFileMetaResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			blob, err := o.srv.NewBlobService().GetFileMeta(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, GetFileMetaResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, GetFileMetaResponse{Code: http.StatusOK, Msg: "", Blob: blob})
			}
		}
		log.Debug("blob/GetFileMeta ", request)
	}
}
