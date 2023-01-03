package loop

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/nameserver/server"
	"os"
	"strconv"
	"time"
)

func runSessionCleaner() {
	log.Info("Starting session cleaner")
	err := server.BlobSessionManager.SetTimeOut(time.Second * 60) // TODO: make this configurable
	if err != nil {
		return
	}
	go func() {
		for {
			err := server.BlobSessionManager.Clean()
			if err != nil {
				log.Warningln(err)
			}
			log.Debugln("trigger clean, active sessions: ", server.BlobSessionManager.ListSessions())

			time.Sleep(time.Second * 30) // TODO: make this configurable
		}
	}()
}

func runErrorHandler() {

}

func runAlivenessChecker() {

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

	runSessionCleaner()

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
