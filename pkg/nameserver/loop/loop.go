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
	server.BlobDataServerManger = server.NewDataServerManager(server.GlobalServerDesc)
	server.BlobSessionManager = server.NewSessionManager()
	server.BlobLockManager = server.NewLockManager(server.GlobalServerDesc.Opt.Volume)

	stat, err := server.BlobDataServerManger.UUIDProbe()
	if err != nil {
		log.Infoln("not all data server is ready")
	} else {
		if len(stat) > 0 {
			err := server.GlobalServerDesc.SaveConfig()
			if err != nil {
				log.Warningln("failed to save configuration")
			} else {
				log.Infoln("writing updated UUID info to file")
			}
		} else {
			log.Infoln("all data server is ready")
		}

	}

	log.Debugln("uuid:", desc.Opt.UUID)
	log.Debugln("port:", desc.Opt.Network.Port)
	log.Debugln("interface:", desc.Opt.Network.Interface)
	log.Debugln("volume:", desc.Opt.Volume)
	log.Debugln("accessKey:", desc.Opt.Auth.AccessKey)
	log.Debugln("secretKey:", desc.Opt.Auth.SecretKey)
	log.Debugln("dataServers:", server.BlobDataServerManger.ListServers())

	/** 创建Gin Server **/
	ginEngine := gin.New()

	/** 注册认证模块 **/
	/** FIXME: timeout fixed to time.Second*86400 **/
	ginJWT, _ := auth.RegisterAuthModule(
		ginEngine,
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
		v1API, _ = auth.CreateJWTAuthGroup(ginEngine, ginJWT, server.APILayout.V1.Self)
	} else {
		v1API = ginEngine.Group(server.APILayout.V1.Self)
	}
	//v1API = ginEngine.Group(server.APILayout.V1.Self)

	/** /ping 永远是不认证的 **/
	pingGroup := ginEngine.Group(server.APILayout.Ping)
	pingController := ping.NewController(nil)
	pingGroup.GET("", pingController.Info)

	/** /v1/sys **/
	sysPath, _ := filepath.Rel(server.APILayout.V1.Self, server.APILayout.V1.Sys)
	sysGroup := v1API.Group(sysPath)
	sysController := sys.NewController(nil)
	sysGroup.GET("info", sysController.Info)
	sysGroup.GET("session", sysController.GetSession)
	sysGroup.GET("sessions", sysController.GetSessions)

	blobPath, _ := filepath.Rel(server.APILayout.V1.Self, server.APILayout.V1.Blob)
	blobGroup := v1API.Group(blobPath)
	blobController := blob2.NewController(repo.Repo(server.BlobDataServerManger, server.BlobSessionManager, server.BlobLockManager))
	blobGroup.POST("close", blobController.Close)
	blobGroup.POST("flush", blobController.Flush)
	blobGroup.POST("lock", blobController.Lock)
	blobGroup.GET("lock", blobController.GetLock)
	blobGroup.GET("ls", blobController.Ls)
	blobGroup.POST("mkdir", blobController.Mkdir)
	blobGroup.POST("open", blobController.Open)
	blobGroup.GET("read", blobController.Read)
	blobGroup.POST("rm", blobController.Rm)
	blobGroup.POST("rmdir", blobController.Rmdir)
	blobGroup.POST("seek", blobController.Seek)
	blobGroup.POST("truncate", blobController.Truncate)
	blobGroup.POST("unlock", blobController.Unlock)
	blobGroup.POST("write", blobController.Write)

	_ = ginEngine.Run(server.GlobalServerDesc.Opt.Network.Interface + ":" + strconv.Itoa(server.GlobalServerDesc.Opt.Network.Port))

}
