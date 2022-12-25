package v1

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/auth"
	"net/http"
)

func (o *controller) Info(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Debugln("the claims is:", claims)
	user, _ := c.Get(auth.IdentityKeyStr)
	c.IndentedJSON(http.StatusOK, gin.H{
		"accessKey": user.(*auth.User).AccessKey,
		"message":   "OK",
	})
}
