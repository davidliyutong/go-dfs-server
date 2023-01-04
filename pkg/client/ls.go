package client

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	v12 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	v1 "go-dfs-server/pkg/nameserver/client/v1"
	"os"
)

func Ls(cmd *cobra.Command, args []string) {
	opt := config.NewClientOpt()
	vipCfg, err := opt.Parse(cmd)
	if err != nil {
		log.Println("cannot find credential, run login first")
	} else {
		cli := v1.NewNameServerClient(opt.Token, opt.Hostname, opt.Port, opt.UseTLS)
		defer refreshToken(cli, vipCfg)
		isDir, res, err := cli.BlobLs(args[0])
		if err != nil {
			log.Errorln(err)
		} else {
			if !isDir {
				h, err := cli.Open(args[0], os.O_RDONLY)
				if err != nil {
					log.Errorln(err)
					return
				}
				buf, err := json.MarshalIndent(h.Blob(), "", "    ")
				if err != nil {
					log.Errorln(err)
					return
				}
				fmt.Println(string(buf))
			} else {
				fmt.Printf("%-5s %-32s %-12s\n", "type", "name", "size")
				for _, r := range res {
					if r.Type == v12.BlobFileTypeName {
						fmt.Printf("%-5s %-32s %-12d\n", r.Type, r.BaseName, r.Size)
					} else {
						fmt.Printf("%-5s %-32s %-12s\n", r.Type, r.BaseName, "-")
					}
				}
			}
			//log.Println(res)
		}
	}
}
