package client

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	v12 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	v1 "go-dfs-server/pkg/nameserver/client/v1"
	"os"
)

func Pipe(cmd *cobra.Command, args []string) {
	opt := config.NewClientOpt()
	_, err := opt.Parse(cmd)
	if err != nil {
		log.Println("cannot find credential, run login first")
	} else {
		cli := v1.NewNameServerClient(opt.Token, opt.Hostname, opt.Port, opt.UseTLS)
		h, err := cli.Open(args[0], os.O_RDWR)
		if err != nil {
			log.Errorln(err)
			return
		}

		input := os.Stdin

		var total = 0
		for {
			buf := make([]byte, v12.DefaultBlobChunkSize)
			read, err := input.Read(buf)
			if err != nil {
				if err.Error() != "EOF" {
					log.Errorln(err)
				}
				break
			}
			written, err := h.Write(buf[:read])
			if err != nil {
				log.Errorln(err)
				break
			}
			total += written
		}

		log.Debugf("written %d bytes", total)
		err = h.Close()
		if err != nil {
			log.Errorln(err)
		}
	}
}
