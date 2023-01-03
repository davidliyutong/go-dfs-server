package loop

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/dataserver/server"
	"os"
)

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
