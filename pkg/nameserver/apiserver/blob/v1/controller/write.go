package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type WriteRequest struct {
	Session string `form:"session" json:"session"`
	Sync    bool   `form:"sync" json:"sync"`
}

type WriteResponse struct {
	Code int64  `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
	Size int64  `form:"size" json:"size"`
}

func (c2 controller) Write(c *gin.Context) {
	var request WriteRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, WriteResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Session == "" {
			c.IndentedJSON(http.StatusBadRequest, WriteResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			file, err := c.FormFile("file")
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, WriteResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				size, err := c2.srv.NewBlobService().Write(request.Session, request.Sync, file)
				if err != nil {
					c.IndentedJSON(http.StatusInternalServerError, WriteResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
				} else {
					c.IndentedJSON(http.StatusOK, WriteResponse{Code: http.StatusOK, Msg: "", Size: size})
				}
			}
		}
		log.Debug("blob/Write ", request)
	}
}
