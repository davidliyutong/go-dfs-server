package loop

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/auth"
	"go-dfs-server/pkg/config"
	ping "go-dfs-server/pkg/nameserver/apiserver/ping"
	sys "go-dfs-server/pkg/nameserver/apiserver/sys/v1/controller"
	"go-dfs-server/pkg/nameserver/server"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func MainLoop(cmd *cobra.Command, args []string) {
	/** 创建NameServerOption **/
	desc := config.NewNameServerDesc()
	err := desc.Parse(cmd) // 解析参数
	if err != nil {
		log.Fatalln("failed to parse configuration", err)
		os.Exit(1)
	}
	desc.PostParse()
	server.GlobalServerDesc = &desc //  设定全局Option

	log.Debugln("uuid:", desc.Opt.UUID)
	log.Debugln("port:", desc.Opt.Network.Port)
	log.Debugln("interface:", desc.Opt.Network.Interface)
	log.Debugln("volume:", desc.Opt.Volume)
	log.Debugln("accessKey:", desc.Opt.Auth.AccessKey)
	log.Debugln("secretKey:", desc.Opt.Auth.SecretKey)

	/** 创建Gin Server **/
	ginEngine := gin.New()

	/** 注册认证模块 **/
	/** FIXME: timeout fixed to time.Second*86400 **/
	ginJWT, _ := auth.RegisterAuthModule(ginEngine, server.APILayout.Auth.Login, server.APILayout.Auth.Refresh, time.Second*86400, auth.RepoAuthnBasic, auth.RepoAuthzBasic)

	/** 路由组 **/
	var v1API *gin.RouterGroup

	/** 如果开启认证，则创建认证路由组，否则创建普通路由组 **/
	if server.GlobalServerDesc.Opt.AuthIsEnabled() {
		v1API, _ = auth.CreateJWTAuthGroup(ginEngine, ginJWT, server.APILayout.V1.Self)
	} else {
		v1API = ginEngine.Group(server.APILayout.V1.Self)
	}

	/** /ping 永远是不认证的 **/
	pingGroup := ginEngine.Group(server.APILayout.Ping)
	pingController := ping.NewController(nil)
	pingGroup.GET("", pingController.Info)

	/** /v1/sys **/
	sysController := sys.NewController(nil)
	sysPath, _ := filepath.Rel(server.APILayout.V1.Self, server.APILayout.V1.Sys)
	v1API.GET(sysPath, sysController.Info)

	_ = ginEngine.Run(server.GlobalServerDesc.Opt.Network.Interface + ":" + strconv.Itoa(server.GlobalServerDesc.Opt.Network.Port))

}
