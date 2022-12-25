package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type createDirectoryRequest struct {
	Path string `form:"path" json:"path"`
}

type createDirectoryResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) CreateDirectory(c *gin.Context) {
	var request createDirectoryRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, createDirectoryResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, createDirectoryResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().CreateDirectory(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, createDirectoryResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, createDirectoryResponse{Code: http.StatusOK, Msg: ""})
			}
		}
		log.Debug("blob/CreateDirectory ", request)
	}

}
