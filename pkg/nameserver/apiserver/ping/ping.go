package ping

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Controller interface {
	Get(c *gin.Context)
}

func (o *controller) Get(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Debugln("the claims is:", claims)
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type controller struct {
}

type repo interface {
}

func NewController(repo repo) Controller {
	return &controller{}
}
