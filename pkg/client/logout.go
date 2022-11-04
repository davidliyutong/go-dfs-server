package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/utils"
	"os"
	"time"
)

func Logout(cmd *cobra.Command, args []string) {
	log.Debugln("client auth")

	opt := config.GetClientOpt()
	vipCfg, err := opt.Parse(cmd)
	if err != nil {
		log.Infoln("found no configuration, exit")
		os.Exit(0)
	} else {
		opt.Token = ""
		opt.Expire = time.UnixMicro(0)
		log.Println("delete token")
	}

	utils.DumpOption(opt, vipCfg.GetString("_config"), true)
}
