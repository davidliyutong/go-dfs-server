package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type writeChunkRequest struct {
	Path string `form:"path" json:"path"`
	ID   int64  `form:"id" json:"id"`
}

type writeChunkResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
	MD5  string `json:"md5"`
}

func (c2 controller) WriteChunk(c *gin.Context) {
	var request writeChunkRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, writeChunkResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, writeChunkResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			//err = c2.srv.NewBlobService().WriteChunk(request.Path, request.ID)
			file, err := c.FormFile("file")
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, writeChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				err = c2.srv.NewBlobService().WriteChunk(request.Path, request.ID, c, file)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, writeChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
				} else {
					MD5String, err := c2.srv.NewBlobService().GetChunkMD5(request.Path, request.ID)
					if err != nil {
						c.IndentedJSON(http.StatusInternalServerError, writeChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
					} else {
						err = c2.srv.NewBlobService().UpdateMeta(request.Path, request.ID, MD5String)
						if err != nil {
							c.IndentedJSON(http.StatusInternalServerError, writeChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
						} else {
							c.IndentedJSON(http.StatusOK, writeChunkResponse{Code: http.StatusOK, Msg: "", MD5: MD5String})
						}
					}
				}
			}
		}
		log.Debug("blob/writeChunk ", request)
	}

}
