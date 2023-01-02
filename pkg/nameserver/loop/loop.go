package loop

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/auth"
	"go-dfs-server/pkg/config"
	blob2 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/controller"
	repo "go-dfs-server/pkg/nameserver/apiserver/blob/v1/repo/memory"
	sys "go-dfs-server/pkg/nameserver/apiserver/sys/v1/controller"
	"go-dfs-server/pkg/nameserver/server"
	ping "go-dfs-server/pkg/ping/v1"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func runRegistration() {
	err := server.BlobDataServerManger.Register()
	if err != nil {
		log.Warningln(err)
	}
}

func runUUIDProbe() {
	stat, err := server.BlobDataServerManger.UUIDProbe()
	if err != nil {
		log.Infoln("not all data servers are ready")
	} else {
		if len(stat) > 0 {
			err := server.GlobalServerDesc.SaveConfig()
			if err != nil {
				log.Warningln("failed to save configuration")
			} else {
				log.Infoln("writing updated UUID info to file")
			}
		} else {
			log.Infoln("all data servers are ready")
		}
	}
}

func registerPingGroup(router *gin.Engine) {
	grp := router.Group(server.APILayout.Ping)
	controller := ping.NewController(nil)
	grp.GET("", controller.Info)
}

func registerV1Group(router *gin.Engine) {
	/** 注册认证模块 **/
	/** FIXME: timeout fixed to time.Second*86400 **/
	ginJWT, _ := auth.RegisterAuthModule(
		router,
		server.APILayout.Auth.Self,
		server.APILayout.Auth.Login,
		server.APILayout.Auth.Refresh,
		time.Second*86400,
		auth.RepoAuthnBasic,
		auth.RepoAuthzBasic)

	/** 路由组 **/
	var v1API *gin.RouterGroup

	/** 如果开启认证，则创建认证路由组，否则创建普通路由组 **/
	if server.GlobalServerDesc.Opt.AuthIsEnabled() {
		v1API, _ = auth.CreateJWTAuthGroup(router, ginJWT, server.APILayout.V1.Self)
	} else {
		v1API = router.Group(server.APILayout.V1.Self)
	}

	/** /v1/sys **/
	sysPath, _ := filepath.Rel(server.APILayout.V1.Self, server.APILayout.V1.Sys.Self)
	sysGroup := v1API.Group(sysPath)
	sysController := sys.NewController(nil)
	sysGroup.GET(server.APILayout.V1.Sys.Info, sysController.Info)
	sysGroup.GET(server.APILayout.V1.Sys.Session, sysController.GetSession)
	sysGroup.GET(server.APILayout.V1.Sys.Sessions, sysController.GetSessions)
	sysGroup.GET(server.APILayout.V1.Sys.Servers, sysController.GetServers)

	blobPath, _ := filepath.Rel(server.APILayout.V1.Self, server.APILayout.V1.Blob.Self)
	blobGroup := v1API.Group(blobPath)
	blobController := blob2.NewController(repo.Repo(server.BlobDataServerManger, server.BlobSessionManager, server.BlobLockManager))

	blobGroup.POST(server.APILayout.V1.Blob.Lock, blobController.Lock)
	blobGroup.GET(server.APILayout.V1.Blob.Lock, blobController.GetLock)
	blobGroup.DELETE(server.APILayout.V1.Blob.Lock, blobController.Unlock)

	blobGroup.GET(server.APILayout.V1.Blob.Meta, blobController.GetFileMeta)

	blobGroup.GET(server.APILayout.V1.Blob.Path, blobController.Ls)
	blobGroup.POST(server.APILayout.V1.Blob.Path, blobController.Mkdir)
	blobGroup.DELETE(server.APILayout.V1.Blob.Path, blobController.Rm)

	blobGroup.GET(server.APILayout.V1.Blob.Session, blobController.Open)
	blobGroup.DELETE(server.APILayout.V1.Blob.Session, blobController.Close)
	blobGroup.POST(server.APILayout.V1.Blob.Session, blobController.Flush)

	blobGroup.GET(server.APILayout.V1.Blob.IO, blobController.Read)
	blobGroup.POST(server.APILayout.V1.Blob.IO, blobController.Write)
	blobGroup.DELETE(server.APILayout.V1.Blob.IO, blobController.Truncate)

	blobGroup.POST(server.APILayout.V1.Blob.Seek, blobController.Seek)
}

func createServer() *gin.Engine {
	router := gin.New()
	registerPingGroup(router)
	registerV1Group(router)
	return router
}

func MainLoop(cmd *cobra.Command, args []string) {
	/** 创建NameServerOption **/
	desc := config.NewNameServerDesc()
	if err := desc.Parse(cmd); err != nil {
		log.Fatalln("failed to parse configuration", err)
		os.Exit(1)
	} else {
		desc.PostParse()
		server.GlobalServerDesc = &desc //  设定全局Option
		server.BlobDataServerManger = server.NewDataServerManager(server.GlobalServerDesc)
		server.BlobSessionManager = server.NewSessionManager()
		server.BlobLockManager = server.NewLockManager(server.GlobalServerDesc.Opt.Volume)
	}

	runRegistration()
	runUUIDProbe()

	log.Infoln("uuid:", desc.Opt.UUID)
	log.Infoln("port:", desc.Opt.Network.Port)
	log.Infoln("interface:", desc.Opt.Network.Interface)
	log.Infoln("volume:", desc.Opt.Volume)
	log.Infoln("accessKey:", desc.Opt.Auth.AccessKey)
	log.Infoln("secretKey:", desc.Opt.Auth.SecretKey)
	log.Infoln("dataServers:", server.BlobDataServerManger.ListServers())

	/** 创建Gin Server **/
	ginEngine := createServer()

	_ = ginEngine.Run(server.GlobalServerDesc.Opt.Network.Interface + ":" + strconv.Itoa(server.GlobalServerDesc.Opt.Network.Port))

}
