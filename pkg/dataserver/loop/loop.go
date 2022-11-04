package loop

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	"os"
	"time"
)

func MainLoop(cmd *cobra.Command, args []string) {
	opt := config.GetDataserverOpt()
	_, err := opt.Parse(cmd)
	if err != nil {
		log.Infoln("failed to parse configuration", err)
		os.Exit(1)
	}
	opt.PostParse()

	log.Debugln("port:", opt.Network.Port)
	log.Debugln("remote:", opt.Network.Remote)
	log.Debugln("endpoint:", opt.Network.Endpoint)
	log.Debugln("volume:", opt.Volume)
	log.Debugln("accessKey:", opt.Auth.AccessKey)
	log.Debugln("secretKey:", opt.Auth.SecretKey)

	for true {
		time.Sleep(1000)
	}

}
