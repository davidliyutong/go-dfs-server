package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type RmRequest struct {
	Path      string `form:"path" json:"path"`
	Recursive bool   `form:"recursive" json:"recursive"`
}

type RmResponse struct {
	Code int64  `form:"code" json:"code"`
	Msg  string `form:"msg" json:"msg"`
}

func (c2 controller) Rm(c *gin.Context) {
	var request RmRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, RmResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, RmResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			err = c2.srv.NewBlobService().Rm(request.Path, request.Recursive)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, RmResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, RmResponse{Code: http.StatusOK, Msg: ""})
			}
		}
	}
	log.Debug("blob/Rm ", request)
}
