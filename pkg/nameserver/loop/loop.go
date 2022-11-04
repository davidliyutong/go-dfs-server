package loop

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/auth"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/nameserver/apiserver/info"
	"go-dfs-server/pkg/nameserver/apiserver/ping"
	"go-dfs-server/pkg/nameserver/server"
	"os"
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

	pingGroup := ginEngine.Group(server.APILayout.Ping)
	pingController := ping.NewController(nil)
	pingGroup.GET("/", pingController.Get)

	v1, _ := auth.CreateJWTAuthGroup(ginEngine, ginJWT, server.APILayout.V1.Self)
	infoController := info.NewController(nil)
	v1.GET(server.APILayout.V1.Info, infoController.Get)

	_ = ginEngine.Run(server.GlobalServerOpt.Network.Interface + ":" + strconv.Itoa(server.GlobalServerOpt.Network.Port))

}
