package loop

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	blob "go-dfs-server/pkg/dataserver/apiserver/blob/v1/controller"
	sys "go-dfs-server/pkg/dataserver/apiserver/sys/v1/controller"
	"go-dfs-server/pkg/dataserver/server"
	ping "go-dfs-server/pkg/ping/v1"
	"os"
)

func registerPingGroup(router *gin.Engine) {
	grp := router.Group(server.APILayout.Ping)
	controller := ping.NewController(nil)
	grp.GET("", controller.Info)
}

func registerBlobGroup(router *gin.Engine) {
	grp := router.Group(server.APILayout.V1.Blob.Self)
	blobController := blob.NewController(nil)
	grp.POST(server.APILayout.V1.Blob.CreateChunk, blobController.CreateChunk)
	grp.POST(server.APILayout.V1.Blob.CreateDirectory, blobController.CreateDirectory)
	grp.POST(server.APILayout.V1.Blob.CreateFile, blobController.CreateFile)
	grp.POST(server.APILayout.V1.Blob.DeleteChunk, blobController.DeleteChunk)
	grp.POST(server.APILayout.V1.Blob.DeleteDirectory, blobController.DeleteDirectory)
	grp.POST(server.APILayout.V1.Blob.DeleteFile, blobController.DeleteFile)
	grp.POST(server.APILayout.V1.Blob.LockFile, blobController.LockFile)
	grp.GET(server.APILayout.V1.Blob.ReadChunk, blobController.ReadChunk)
	grp.GET(server.APILayout.V1.Blob.ReadChunkMeta, blobController.ReadChunkMeta)
	grp.GET(server.APILayout.V1.Blob.ReadFileMeta, blobController.ReadFileMeta)
	grp.GET(server.APILayout.V1.Blob.ReadFileLock, blobController.ReadFileLock)
	grp.POST(server.APILayout.V1.Blob.UnlockFile, blobController.UnlockFile)
	grp.PUT(server.APILayout.V1.Blob.WriteChunk, blobController.WriteChunk)
}

func registerSysGroup(router *gin.Engine) {
	grp := router.Group(server.APILayout.V1.Sys.Self)
	sysController := sys.NewController(nil)
	grp.GET(server.APILayout.V1.Sys.Info, sysController.Info)
	grp.GET(server.APILayout.V1.Sys.UUID, sysController.UUID)
	grp.GET(server.APILayout.V1.Sys.Config, sysController.Config)
	grp.POST(server.APILayout.V1.Sys.Register, sysController.Register)

}

func createServer() *gin.Engine {
	router := gin.New()
	registerPingGroup(router)
	registerBlobGroup(router)
	registerSysGroup(router)

	return router
}

func MainLoop(cmd *cobra.Command, args []string) {
	desc := config.NewDataServerDesc()
	if err := desc.Parse(cmd); err != nil {
		log.Infoln("failed to parse configuration", err)
		os.Exit(1)
	} else {
		desc.PostParse()
		server.GlobalServerDesc = &desc                           //  设定全局Option
		server.GlobalFileLocks = make(map[string]map[string]bool) // 设定全局文件锁数据库
	}

	/** End of server init */
	log.Infoln("uuid:", desc.Opt.UUID)
	log.Infoln("port:", desc.Opt.Network.Port)
	log.Infoln("endpoint:", desc.Opt.Network.Endpoint)
	log.Infoln("volume:", desc.Opt.Volume)

	ginEngine := createServer()

	log.Debugln()
	_ = ginEngine.Run(server.GlobalServerDesc.Opt.Network.Endpoint)

}
