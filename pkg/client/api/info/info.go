package info

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/nameserver/server"
	"io"
	"net/http"
	"os"
)

type Client struct {
	*config.ClientOpt
}

type ClusterInfo struct {
	AccessKey string `json:"accessKey"`
	Message   string `json:"message"`
}

func (o *Client) Info(a *config.ClientAuthOpt) (ClusterInfo, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", o.GetHTTPUrl()+server.APILayout.V1.Sys, nil)
	if err != nil {
		return ClusterInfo{}, err
	}

	if a.AuthIsEnabled() {
		request.Header.Add("Authorization", "Bearer "+o.Token)
	}
	request.Header.Add("Content-Type", "application/json")

	respHandle, err := client.Do(request)
	if err != nil {
		return ClusterInfo{}, err
	}

	defer func() {
		err := respHandle.Body.Close()
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
	}()

	content, err := io.ReadAll(respHandle.Body)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	log.Debugln("dfs response:", string(content))

	response := ClusterInfo{}
	err = json.Unmarshal(content, &response)
	if err != nil {
		return ClusterInfo{}, err
	}

	return response, nil
}
