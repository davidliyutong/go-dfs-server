package ping

import (
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/nameserver/server"
	"io"
	"net/http"
)

type Client struct {
	*config.ClientOpt
}

func (o *Client) Ping() (bool, error) {
	respHandle, err := http.Get(o.GetHTTPUrl() + server.APILayout.Ping)
	if err != nil {
		return false, err
	}

	defer func() {
		err := respHandle.Body.Close()
		if err != nil {
			log.Errorln(err)
		}
	}()

	content, err := io.ReadAll(respHandle.Body)
	if err != nil {
		log.Errorln(err)
		return false, err
	}
	log.Debugln("dfs response:", string(content))

	return true, nil
}
