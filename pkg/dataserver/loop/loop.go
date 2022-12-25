package loop

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	blob "go-dfs-server/pkg/dataserver/apiserver/blob/v1/controller"
	sys "go-dfs-server/pkg/dataserver/apiserver/sys/v1/controller"
	"go-dfs-server/pkg/dataserver/ping"
	"go-dfs-server/pkg/dataserver/server"
	"os"
)

func MainLoop(cmd *cobra.Command, args []string) {
	desc := config.NewDataServerDesc()
	err := desc.Parse(cmd)
	if err != nil {
		log.Infoln("failed to parse configuration", err)
		os.Exit(1)
	}

	if desc.Opt.UUID == "" {
		log.Infoln("uuid is empty")
		os.Exit(1)
	}

	desc.PostParse()
	server.GlobalServerDesc = &desc                           //  设定全局Option
	server.GlobalFileLocks = make(map[string]map[string]bool) // 设定全局文件锁数据库

	/** End of server init */

	log.Debugln("uuid:", desc.Opt.UUID)
	log.Debugln("port:", desc.Opt.Network.Port)
	log.Debugln("endpoint:", desc.Opt.Network.Endpoint)
	log.Debugln("volume:", desc.Opt.Volume)

	ginEngine := gin.New()

	pingGroup := ginEngine.Group(server.APILayout.Ping)
	pingController := ping.NewController(nil)
	pingGroup.GET("", pingController.Info)

	blobGroup := ginEngine.Group(server.APILayout.V1.Blob)
	blobController := blob.NewController(nil)
	blobGroup.PUT("createChunk", blobController.CreateChunk)
	blobGroup.PUT("createDirectory", blobController.CreateDirectory)
	blobGroup.PUT("createFile", blobController.CreateFile)
	blobGroup.PUT("deleteChunk", blobController.DeleteChunk)
	blobGroup.PUT("deleteDirectory", blobController.DeleteDirectory)
	blobGroup.PUT("deleteFile", blobController.DeleteFile)
	blobGroup.POST("lockFile", blobController.LockFile)
	blobGroup.GET("readChunk", blobController.ReadChunk)
	blobGroup.GET("readChunkMeta", blobController.ReadChunkMeta)
	blobGroup.GET("readFileMeta", blobController.ReadFileMeta)
	blobGroup.GET("readFileLock", blobController.ReadFileLock)
	blobGroup.POST("unlockFile", blobController.UnlockFile)
	blobGroup.PUT("writeChunk", blobController.WriteChunk)

	sysGroup := ginEngine.Group(server.APILayout.V1.Sys)
	sysController := sys.NewController(nil)
	sysGroup.GET("info", sysController.Info)
	sysGroup.GET("uuid", sysController.UUID)
	sysGroup.GET("config", sysController.Config)

	log.Debugln()
	_ = ginEngine.Run(server.GlobalServerDesc.Opt.Network.Endpoint)

}
