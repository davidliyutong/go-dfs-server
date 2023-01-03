package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type LsRequest struct {
	Path string `form:"path" json:"path"`
}

type LsFileInfo struct {
	BaseName string `form:"basename" json:"basename"`
	Type     string `form:"type" json:"type"`
}
type LsResponse struct {
	Code int64        `form:"code" json:"code"`
	Msg  string       `form:"msg" json:"msg"`
	List []LsFileInfo `form:"list" json:"list"`
}

func (c2 controller) Ls(c *gin.Context) {
	var request LsRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, LsResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {
		if request.Path == "" {
			c.IndentedJSON(http.StatusBadRequest, LsResponse{Code: http.StatusBadRequest, Msg: "wrong parameter"})
		} else {
			info, err := c2.srv.NewBlobService().Ls(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, LsResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				resp := LsResponse{Code: http.StatusOK, Msg: ""}
				for _, v := range info {
					log.Debug(v)
					resp.List = append(resp.List, LsFileInfo{BaseName: v.BaseName, Type: v.Type})
				}
				c.IndentedJSON(http.StatusOK, resp)
			}
		}
	}
	log.Debug("blob/Ls ", request)
}
