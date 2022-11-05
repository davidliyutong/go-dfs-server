package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/auth"
	"go-dfs-server/pkg/client/api/info"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/nameserver/server"
	"io"
	"net/http"
	"os"
)

type Client struct {
	*config.ClientOpt
}

type AuthClusterInfo struct {
	AccessKey string `json:"accessKey"`
	Message   string `json:"message"`
}

func (o *Client) AuthLogin(a *config.ClientAuthOpt) error {

	credentials := auth.ClientLoginCredential{
		AccessKey: a.AccessKey,
		SecretKey: a.SecretKey,
	}
	credentialsBytes, _ := json.Marshal(credentials)

	respHandle, err := http.Post(o.GetHTTPUrl()+server.APILayout.Auth.Login, "application/json", bytes.NewBuffer(credentialsBytes))
	if err != nil {
		return err
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

	response := auth.JWTResponse{}
	err = json.Unmarshal(content, &response)
	if err != nil {
		return err
	}

	if response.Code == 200 {
		o.Token = response.Token
		o.Expire = response.Expire
	} else {
		return errors.New("login failed, access denied, try logout then login")
	}

	infoClient := info.Client{ClientOpt: o.ClientOpt}
	_, err = infoClient.Info(a)
	if err != nil {
		return err
	}

	return nil
}

func (o *Client) MustAuthLogin(a *config.ClientAuthOpt) {
	err := o.AuthLogin(a)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

func (o *Client) AuthRefresh() error {
	client := &http.Client{}

	request, err := http.NewRequest("POST", o.GetHTTPUrl()+server.APILayout.Auth.Refresh, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Bearer "+o.Token)
	request.Header.Add("Content-Type", "application/json")

	respHandle, err := client.Do(request)
	if err != nil {
		return err
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

	response := auth.JWTResponse{}
	err = json.Unmarshal(content, &response)
	if err != nil {
		return err
	}

	if response.Code == 200 {
		o.Token = response.Token
		o.Expire = response.Expire
	} else {
		return errors.New("login failed, access denied, try logout then login")
	}

	return nil
}

func (o *Client) MustAuthRefresh() {
	err := o.AuthRefresh()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

}
