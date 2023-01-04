package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	v1 "go-dfs-server/pkg/nameserver/client/v1"
)

func Mkdir(cmd *cobra.Command, args []string) {
	opt := config.NewClientOpt()
	vipCfg, err := opt.Parse(cmd)
	if err != nil {
		log.Println("cannot find credential, run login first")
	} else {
		cli := v1.NewNameServerClient(opt.Token, opt.Hostname, opt.Port, opt.UseTLS)
		defer refreshToken(cli, vipCfg)
		err := cli.BlobMkdir(args[0])
		if err != nil {
			log.Errorln(err)
		}
	}
}
