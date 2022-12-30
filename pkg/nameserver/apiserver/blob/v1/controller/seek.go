package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type SeekRequest struct {
	Session string `form:"session" json:"session"`
	Offset  int64  `form:"offset" json:"offset"`
	Whence  int    `form:"whence" json:"whence"`
}

type SeekResponse struct {
	Code   int64  `form:"code" json:"code"`
	Msg    string `form:"msg" json:"msg"`
	Offset int64  `form:"offset" json:"offset"`
}

func (c2 controller) Seek(c *gin.Context) {
	var request SeekRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, SeekResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Session == "" {
			c.IndentedJSON(http.StatusBadRequest, SeekResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			offset, err := c2.srv.NewBlobService().Seek(request.Session, request.Offset, request.Whence)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, SeekResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, SeekResponse{Code: http.StatusOK, Msg: "", Offset: offset})
			}
		}
	}
	log.Debug("blob/Seek ", request)

}
