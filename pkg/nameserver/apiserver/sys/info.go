package sys

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/auth"
)

type Controller interface {
	Get(c *gin.Context)
}

func (o *controller) Get(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Debugln("the claims is:", claims)
	user, _ := c.Get(auth.IdentityKeyStr)
	c.IndentedJSON(200, gin.H{
		"accessKey": user.(*auth.User).AccessKey,
		"message":   "OK",
	})
}

type controller struct {
}

type repo interface {
	InfoQuery()
}

func NewController(repo repo) Controller {
	return &controller{}
}
