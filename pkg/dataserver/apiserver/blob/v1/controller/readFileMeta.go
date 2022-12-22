package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type readFileMetaRequest struct {
	Path string `form:"path" json:"path"`
}

type readFileMetaResponse struct {
	Code int64            `json:"code"`
	Msg  string           `json:"msg"`
	MD5  map[int64]string `json:"md5"`
}

func (c2 controller) ReadFileMeta(c *gin.Context) {
	var request readFileMetaRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(400, readFileMetaResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, readFileMetaResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			meta, err := c2.srv.NewBlobService().ReadFileMeta(request.Path)
			if err != nil {
				c.IndentedJSON(500, readFileMetaResponse{Code: 500, Msg: err.Error()})
			} else {
				c.IndentedJSON(200, readFileMetaResponse{Code: 200, Msg: "", MD5: meta.Content})
			}
		}
		log.Debug("blob/readFileMeta ", request)
	}
}
