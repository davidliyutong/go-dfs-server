package loop

import (
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/nameserver/server"
)

func runRegistration() {
	err := server.BlobDataServerManger.Register()
	if err != nil {
		log.Warningln(err)
	}
}

func runUUIDProbe() {
	stat, err := server.BlobDataServerManger.UUIDProbe()
	if err != nil {
		log.Infoln("not all data servers are ready")
	} else {
		if len(stat) > 0 {
			err := server.GlobalServerDesc.SaveConfig()
			if err != nil {
				log.Warningln("failed to save configuration")
			} else {
				log.Infoln("writing updated UUID info to file")
			}
		} else {
			log.Infoln("all data servers are ready")
		}
	}
}
