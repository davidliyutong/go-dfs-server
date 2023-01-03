package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	model "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"net/http"
)

type SyncRequest struct {
	Path string             `form:"path" json:"path"`
	Blob model.BlobMetaData `form:"blob" json:"blob"`
}

type SyncResponse struct {
	Code int64              `form:"code" json:"code"`
	Msg  string             `form:"msg" json:"msg"`
	Blob model.BlobMetaData `form:"blob" json:"blob"`
}

func (c2 controller) Sync(c *gin.Context) {
	var request SyncRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, SyncResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, SyncResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			blob, err := c2.srv.NewBlobService().Sync(request.Path, request.Blob)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, SyncResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, SyncResponse{Code: http.StatusOK, Msg: "", Blob: blob})
			}
		}
		log.Debug("blob/Sync ", request)
	}
}
