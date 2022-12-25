package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type createChunkRequest struct {
	Path string `form:"path" json:"path"`
	ID   int64  `form:"id" json:"id"`
}

type createChunkResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) CreateChunk(c *gin.Context) {
	var request createChunkRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, createChunkResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, createChunkResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().CreateChunk(request.Path, request.ID)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, createChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, createChunkResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/CreateChunk ", request)
	}

}
