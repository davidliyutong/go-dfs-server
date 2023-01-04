package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"net/http"
)

type LsRequest struct {
	Path string `form:"path" json:"path"`
}

type LsFileInfo struct {
	BaseName string `form:"basename" json:"basename"`
	Type     string `form:"type" json:"type"`
	Size     int64  `form:"size" json:"size"`
}

func (o *LsFileInfo) IsDir() bool {
	return o.Type == v1.BlobDirTypeName
}

type LsResponse struct {
	Code  int64        `form:"code" json:"code"`
	Msg   string       `form:"msg" json:"msg"`
	IsDir bool         `form:"is_dir" json:"is_dir"`
	List  []LsFileInfo `form:"list" json:"list"`
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
			isDir, info, err := c2.srv.NewBlobService().Ls(request.Path)
			if err != nil {
				c.IndentedJSON(http.StatusInternalServerError, LsResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
			} else {
				resp := LsResponse{Code: http.StatusOK, Msg: "", IsDir: isDir}
				for _, v := range info {
					log.Debug(v)
					resp.List = append(resp.List, LsFileInfo{BaseName: v.BaseName, Type: v.Type, Size: v.Size})
				}
				c.IndentedJSON(http.StatusOK, resp)
			}
		}
	}
	log.Debug("blob/Ls ", request)
}
