package loop

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/auth"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/nameserver/apiserver/ping"
	"go-dfs-server/pkg/nameserver/apiserver/sys"
	"go-dfs-server/pkg/nameserver/server"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func MainLoop(cmd *cobra.Command, args []string) {
	opt := config.GetNameserverOpt()
	_, err := opt.Parse(cmd)
	if err != nil {
		log.Infoln("failed to parse configuration", err)
		os.Exit(1)
	}
	opt.PostParse()
	server.GlobalServerOpt = &opt

	log.Debugln("port:", opt.Network.Port)
	log.Debugln("interface:", opt.Network.Interface)
	log.Debugln("volume:", opt.Volume)
	log.Debugln("accessKey:", opt.Auth.AccessKey)
	log.Debugln("secretKey:", opt.Auth.SecretKey)

	ginEngine := gin.New()
	ginJWT, _ := auth.RegisterAuthModule(ginEngine, server.APILayout.Auth.Login, server.APILayout.Auth.Refresh, time.Second*86400, auth.RepoAuthnBasic, auth.RepoAuthzBasic)

	var v1API *gin.RouterGroup
	var infoAPI *gin.RouterGroup
	if server.GlobalServerOpt.AuthIsEnabled() {
		v1API, _ = auth.CreateJWTAuthGroup(ginEngine, ginJWT, server.APILayout.V1.Self)
		infoAPI, _ = auth.CreateJWTAuthGroup(ginEngine, ginJWT, server.APILayout.Info)
	} else {
		v1API = ginEngine.Group(server.APILayout.V1.Self)
		infoAPI = ginEngine.Group(server.APILayout.Info)
	}

	pingGroup := ginEngine.Group(server.APILayout.Ping)
	pingController := ping.NewController(nil)
	pingGroup.GET("", pingController.Get)
	infoController := sys.NewController(nil)
	infoPath, _ := filepath.Rel(server.APILayout.V1.Self, server.APILayout.V1.Info)

	v1API.GET(infoPath, infoController.Get)
	infoAPI.GET("", infoController.Get)

	_ = ginEngine.Run(server.GlobalServerOpt.Network.Interface + ":" + strconv.Itoa(server.GlobalServerOpt.Network.Port))

}
