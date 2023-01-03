package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"net/http"
)

type OpenRequest struct {
	Path string `form:"path" json:"path"`
	Mode int    `form:"mode" json:"mode"`
}

type OpenResponse struct {
	Code int64           `form:"code" json:"code"`
	Msg  string          `form:"msg" json:"msg"`
	Blob v1.BlobMetaData `json:"blob"`
}

func (c2 controller) Open(c *gin.Context) {
	var request OpenRequest

	err := c.ShouldBind(&request)
	var blob v1.BlobMetaData
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, OpenResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, OpenResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			blob, err = c2.srv.NewBlobService().Open(request.Path, request.Mode)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, OpenResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, OpenResponse{Code: http.StatusOK, Msg: "", Blob: blob})
			}
		}
	}
	log.Debug("blob/Open ", request)

}
