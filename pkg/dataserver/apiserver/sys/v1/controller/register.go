package v1

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/dataserver/server"
	"net/http"
)

type RegisterRequest struct {
	UUID string `json:"uuid"`
}

type RegisterResponse struct {
	Code int16  `json:"code"`
	Msg  string `json:"msg"`
	UUID string `json:"uuid"`
}

func (o *controller) Register(c *gin.Context) {
	var request RegisterRequest

	err := c.ShouldBind(&request)
	if err != nil {
		log.Debug(err)
		c.IndentedJSON(http.StatusBadRequest, RegisterResponse{Code: http.StatusBadRequest, Msg: "failed"})
	} else {

		lastUUID := server.GlobalServerDesc.NameServerUUID
		server.GlobalServerDesc.NameServerUUID = request.UUID
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, RegisterResponse{Code: http.StatusInternalServerError, Msg: err.Error()})
		} else {
			c.IndentedJSON(http.StatusOK, RegisterResponse{Code: http.StatusOK, Msg: "", UUID: lastUUID})
		}

		log.Debug("blob/CreateChunk ", request)
	}

}
