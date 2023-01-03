package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type GetSessionRequest struct {
	Path string `form:"path" json:"path"`
}

type GetSessionResponse struct {
	Code   int64     `form:"code" json:"code"`
	Msg    string    `form:"msg" json:"msg"`
	Path   string    `form:"path" json:"path"`
	Time   time.Time `form:"time" json:"time"`
	Opened bool      `form:"opened" json:"opened"`
	Offset int64     `form:"offset" json:"offset"`
	Size   int64     `form:"size" json:"size"`
}

func (o *controller) GetSession(c *gin.Context) {
	var request GetSessionRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, GetSessionResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, GetSessionResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			session, err := o.srv.NewSysService().GetSession(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, GetSessionResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, GetSessionResponse{
					Code:   http.StatusOK,
					Msg:    "",
					Path:   *session.Path(),
					Time:   session.GetTime(),
					Opened: session.IsOpened(),
					Size:   session.GetBlobMetaData().Size})
			}
		}
	}
	log.Debug("blob/GetSession ", request)
}
