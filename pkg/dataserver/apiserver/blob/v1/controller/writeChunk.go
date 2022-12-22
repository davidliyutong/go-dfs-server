package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
		c.IndentedJSON(400, writeChunkResponse{Code: 400, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(400, writeChunkResponse{Code: 400, Msg: "wrong parameter"})
		} else {
			//err = c2.srv.NewBlobService().WriteChunk(request.Path, request.ID)
			file, err := c.FormFile("file")
			if err != nil {
				c.IndentedJSON(500, writeChunkResponse{Code: 500, Msg: err.Error()})
			} else {
				err = c2.srv.NewBlobService().WriteChunk(request.Path, request.ID, c, file)
				if err != nil {
					c.IndentedJSON(500, writeChunkResponse{Code: 500, Msg: err.Error()})
				} else {
					MD5String, err := c2.srv.NewBlobService().GetChunkMD5(request.Path, request.ID)
					if err != nil {
						c.IndentedJSON(500, writeChunkResponse{Code: 500, Msg: err.Error()})
					} else {
						err = c2.srv.NewBlobService().UpdateMeta(request.Path, request.ID, MD5String)
						if err != nil {
							c.IndentedJSON(500, writeChunkResponse{Code: 500, Msg: err.Error()})
						} else {
							c.IndentedJSON(200, writeChunkResponse{Code: 200, Msg: "", MD5: MD5String})
						}
					}
				}
			}
		}
		log.Debug("blob/writeChunk ", request)
	}

}
