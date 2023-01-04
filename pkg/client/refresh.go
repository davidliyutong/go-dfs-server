package client

import (
	"github.com/spf13/viper"
	v1 "go-dfs-server/pkg/nameserver/client/v1"
	"go-dfs-server/pkg/utils"
)

func refreshToken(cli v1.NameServerClient, vipCfg *viper.Viper) {
	_, err := cli.AuthRefresh()
	if err != nil {
		utils.DumpOption(cli.Opt(), vipCfg.GetString("_config"), true)
	}
}
