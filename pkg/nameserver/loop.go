package nameserver

import (
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/nameserver/config"
	"log"
	"time"
)

func MainLoop(cmd *cobra.Command, args []string) {
	_, cfg, _ := config.Parse(cmd)
	log.Println("port:", cfg.Network.Port)
	log.Println("interface:", cfg.Network.Interface)

	log.Println("volume:", cfg.Volume)

	log.Println("accessKey:", cfg.Auth.AccessKey)
	log.Println("secretKey:", cfg.Auth.SecretKey)

	for true {
		time.Sleep(1000)
	}

}
