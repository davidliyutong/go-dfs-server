package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"go-dfs-server/pkg/auth"
	"go-dfs-server/pkg/nameserver/server"
	"io"
	"net/http"
	"os"
)

func (c *nameServerClient) AuthLogin(accessKey string, secretKey string) (string, error) {

	credentials := auth.ClientLoginCredential{
		AccessKey: accessKey,
		SecretKey: secretKey,
	}
	credentialsBytes, _ := json.Marshal(credentials)

	targetUrl, err := c.GetAPIUrl(server.APILayout.Auth.Self, server.APILayout.Auth.Login)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(targetUrl, "application/json", bytes.NewBuffer(credentialsBytes))
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	log.Debugln("dfs response:", string(body))

	var result auth.JWTResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if result.Code == 200 {
		c.opt.Token = result.Token
		c.opt.Expire = result.Expire
	} else {
		return "", errors.New("login failed, access denied, try logout then login")
	}

	info, err := c.SysInfo()
	if err != nil {
		return "", err
	} else {
		log.Info(info)
		return c.opt.Token, nil
	}

}

func (c *nameServerClient) AuthRefresh() (string, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.Auth.Self, server.APILayout.Auth.Refresh)
	req, _ := http.NewRequest("POST", targetUrl, nil)

	req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	log.Debugln("dfs response:", string(body))

	var result auth.JWTResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if result.Code == 200 {
		c.opt.Token = result.Token
		c.opt.Expire = result.Expire
	} else {
		return "", errors.New("login failed, access denied, try logout then login")
	}

	return c.opt.Token, nil
}

func (c *nameServerClient) MustAuthLogin(accessKey string, secretKey string) string {
	token, err := c.AuthLogin(accessKey, secretKey)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	} else {
		return token
	}
	return ""
}

func (c *nameServerClient) MustAuthRefresh() string {
	token, err := c.AuthRefresh()
	if err != nil {
		log.Warning(err)
		os.Exit(1)
	} else {
		return token
	}
	return ""
}

func (c *nameServerClient) AuthIsEnabled() bool {
	return c.opt.AuthIsEnabled()
}
