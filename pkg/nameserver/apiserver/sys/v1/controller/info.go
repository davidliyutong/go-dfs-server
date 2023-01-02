package v1

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/auth"
	"go-dfs-server/pkg/nameserver/server"
	"net/http"
)

type InfoResponse struct {
	Code      int64  `form:"code" json:"code"`
	Msg       string `form:"msg" json:"msg"`
	AccessKey string `form:"accessKey" json:"accessKey"`
	Role      string `form:"role" json:"role"`
	Version   string `form:"version" json:"version"`
}

func (o *controller) Info(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Debugln("the claims is:", claims)
	user, _ := c.Get(auth.IdentityKeyStr)
	c.IndentedJSON(http.StatusOK, InfoResponse{
		Code:      http.StatusOK,
		Msg:       "",
		AccessKey: user.(*auth.User).AccessKey,
		Role:      "nameserver",
		Version:   server.NameServerAPIVersion,
	})
	log.Debug("blob/info ")
}
