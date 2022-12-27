package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type DeleteChunkRequest struct {
	Path string `form:"path" json:"path"`
	ID   int64  `form:"id" json:"id"`
}

type DeleteChunkResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) DeleteChunk(c *gin.Context) {
	var request DeleteChunkRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, DeleteChunkResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, DeleteChunkResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().DeleteChunk(request.Path, request.ID)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, DeleteChunkResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, DeleteChunkResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/DeleteChunk ", request)
	}
}
