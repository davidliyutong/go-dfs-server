package loop

import (
	"github.com/gin-gonic/gin"
	"go-dfs-server/pkg/auth"
	blob2 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/controller"
	repo "go-dfs-server/pkg/nameserver/apiserver/blob/v1/repo/memory"
	sys "go-dfs-server/pkg/nameserver/apiserver/sys/v1/controller"
	"go-dfs-server/pkg/nameserver/server"
	ping "go-dfs-server/pkg/ping/v1"
	"path/filepath"
	"time"
)

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
	blobController := blob2.NewController(repo.Repo(server.GlobalServerDesc.Opt.Volume, server.BlobDataServerManger, server.BlobSessionManager))

	blobGroup.GET(server.APILayout.V1.Blob.Path, blobController.Ls)
	blobGroup.POST(server.APILayout.V1.Blob.Path, blobController.Mkdir)
	blobGroup.DELETE(server.APILayout.V1.Blob.Path, blobController.Rm)

	blobGroup.GET(server.APILayout.V1.Blob.File, blobController.Open)
	blobGroup.POST(server.APILayout.V1.Blob.File, blobController.Sync)

	blobGroup.GET(server.APILayout.V1.Blob.IO, blobController.Read)
	blobGroup.POST(server.APILayout.V1.Blob.IO, blobController.Write)

}

func createServer() *gin.Engine {
	router := gin.New()
	registerPingGroup(router)
	registerV1Group(router)
	return router
}
