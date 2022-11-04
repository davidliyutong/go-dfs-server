package login

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/client/api/auth"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/utils"
	"os"
)

func Login(cmd *cobra.Command, args []string) {
	log.Debugln("client auth")

	opt := config.GetClientOpt()
	vipCfg, err := opt.Parse(cmd)
	if err != nil {
		if len(args) <= 0 {
			log.Errorln("no url specified")
			os.Exit(1)
		}
		opt.MustBindURL(args[0])
		opt.MustBindAuthentication(cmd)
		dfsClient := auth.Client{ClientOpt: &opt}
		dfsClient.MustAuthLogin()
		log.Println("login success")
	} else {
		log.Debugln("%s", opt)
		if len(args) > 0 {
			opt.MustBindURL(args[0])
			opt.MustBindAuthentication(cmd)
		}
		dfsClient := auth.Client{ClientOpt: &opt}
		err = dfsClient.AuthRefresh()
		if err != nil {
			dfsClient.MustAuthLogin()
		}
		log.Println("renew token success")
	}

	utils.DumpOption(opt, vipCfg.GetString("_config"), true)
}
