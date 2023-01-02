package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	v1 "go-dfs-server/pkg/nameserver/client/v1"
	"go-dfs-server/pkg/utils"
	"os"
)

func Login(cmd *cobra.Command, args []string) {
	log.Debugln("client auth")

	opt := config.NewClientOpt()
	authOpt := config.NewClientAuthOpt()
	vipCfg, err := opt.Parse(cmd)
	if err != nil {
		if len(args) <= 0 {
			log.Errorln("no url specified")
			os.Exit(1)
		}
		opt.MustBindURL(args[0])
		authOpt.MustBindAuthentication(cmd)
		dfsClient := v1.NewNameServerClientFromOpt(opt)
		dfsClient.MustAuthLogin(authOpt.AccessKey, authOpt.SecretKey)
		log.Println("login success")
	} else {
		log.Debugln("%s", opt)
		if len(args) > 0 {
			opt.MustBindURL(args[0])
			authOpt.MustBindAuthentication(cmd)
		}
		dfsClient := v1.NewNameServerClientFromOpt(opt)
		_, err := dfsClient.AuthRefresh()
		if err != nil {
			_, err = dfsClient.AuthLogin(authOpt.AccessKey, authOpt.SecretKey)
			if err != nil {
				authOpt.MustBindAuthentication(nil)
				dfsClient.MustAuthLogin(authOpt.AccessKey, authOpt.SecretKey)
				log.Println("login success")
			}
		}
		log.Println("renew token success")
	}

	utils.DumpOption(opt, vipCfg.GetString("_config"), true)
}
