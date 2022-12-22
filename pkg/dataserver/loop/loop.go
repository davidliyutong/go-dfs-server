package loop

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	v1 "go-dfs-server/pkg/dataserver/apiserver/blob/v1/controller"
	"go-dfs-server/pkg/dataserver/ping"
	"go-dfs-server/pkg/dataserver/server"
	"os"
)

func MainLoop(cmd *cobra.Command, args []string) {
	desc := config.NewDataserverDesc()
	err := desc.Parse(cmd)
	if err != nil {
		log.Infoln("failed to parse configuration", err)
		os.Exit(1)
	}
	desc.PostParse()
	server.GlobalServerDesc = &desc //  设定全局Option

	log.Debugln("port:", desc.Opt.Network.Port)
	log.Debugln("endpoint:", desc.Opt.Network.Endpoint)
	log.Debugln("volume:", desc.Opt.Volume)

	ginEngine := gin.New()

	pingGroup := ginEngine.Group(server.APILayout.Ping)
	pingController := ping.NewController(nil)
	pingGroup.GET("", pingController.Get)

	blobGroup := ginEngine.Group(server.APILayout.V1.Blob)
	blobController := v1.NewController(nil)
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
	//blobGroup.GET("readFileLock", blobController.ReadFileLock)
	blobGroup.POST("unlockFile", blobController.UnlockFile)
	blobGroup.PUT("writeChunk", blobController.WriteChunk)

	log.Debugln()
	_ = ginEngine.Run(server.GlobalServerDesc.Opt.Network.Endpoint)

}
