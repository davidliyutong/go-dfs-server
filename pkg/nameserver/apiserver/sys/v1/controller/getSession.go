package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type GetSessionRequest struct {
	Session string `form:"session" json:"session"`
}

type GetSessionResponse struct {
	Code   int64     `form:"code" json:"code"`
	Msg    string    `form:"msg" json:"msg"`
	Path   string    `form:"path" json:"path"`
	Mode   int       `form:"mode" json:"mode"`
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
		if request.Session == "" {
			c.IndentedJSON(http.StatusBadRequest, GetSessionResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			session, err := o.srv.NewSysService().GetSession(request.Session)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, GetSessionResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				c.IndentedJSON(http.StatusOK, GetSessionResponse{
					Code:   http.StatusOK,
					Msg:    "",
					Path:   *session.GetPath(),
					Mode:   *session.GetMode(),
					Time:   session.GetTime(),
					Offset: *session.GetOffset(),
					Opened: session.IsOpened(),
					Size:   session.GetBlobMetaData().Size})
			}
		}
	}
	log.Debug("blob/GetSession ", request)
}
