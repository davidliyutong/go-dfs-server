package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type CreateFileRequest struct {
	Path string `form:"path" json:"path"`
}

type CreateFileResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

func (c2 controller) CreateFile(c *gin.Context) {
	var request CreateFileRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, CreateFileResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, CreateFileResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().CreateFile(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, CreateFileResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, CreateFileResponse{Code: http.StatusOK, Msg: ""})
			}
		}
	}
	log.Debug("blob/CreateFile ", request)

}
