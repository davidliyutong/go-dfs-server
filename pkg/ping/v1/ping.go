package v1

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Controller interface {
	Info(c *gin.Context)
}

type PingResponse struct {
	Message string `json:"message"`
}

func (o *controller) Info(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	log.Debugln("the claims is:", claims)
	c.IndentedJSON(http.StatusOK, PingResponse{
		Message: "pong",
	})
}

type controller struct {
}

type repo interface {
}

func NewController(repo repo) Controller {
	return &controller{}
}
