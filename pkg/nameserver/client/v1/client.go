package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-dfs-server/pkg/config"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/controller"
	v13 "go-dfs-server/pkg/nameserver/apiserver/sys/v1/controller"
	"go-dfs-server/pkg/nameserver/server"
	ping "go-dfs-server/pkg/ping/v1"
	"io"
	"net/http"
	"strings"
)

type nameServerClient struct {
	opt config.ClientOpt
}

type NameServerClient interface {
	GetBaseUrl() (string, error)
	GetAPIUrl(keys ...string) (string, error)
	Ping() error
	Opt() *config.ClientOpt

	AuthIsEnabled() bool
	AuthLogin(accessKey string, secretKey string) (string, error)
	MustAuthLogin(accessKey string, secretKey string) string
	AuthRefresh() (string, error)
	MustAuthRefresh() string

	Open(path string, mode int) (Handle, error)

	BlobMkdir(path string) error
	BlobLs(path string) ([]v1.LsFileInfo, error)
	BlobRm(path string, recursive bool) error

	SysInfo() (v13.InfoResponse, error)
	SysSession(sessionID string) (v13.GetSessionResponse, error)
	SysSessions() ([]string, error)
	SysServers() ([]config.RegisteredDataServer, error)
}

var _ NameServerClient = &nameServerClient{}

func NewNameServerClient(token string, hostname string, port int64, useTLS bool) NameServerClient {
	return &nameServerClient{
		opt: config.ClientOpt{
			Token:    token,
			Hostname: hostname,
			Port:     port,
			UseTLS:   useTLS,
		},
	}
}

func NewNameServerClientFromOpt(opt config.ClientOpt) NameServerClient {
	return &nameServerClient{
		opt: opt,
	}
}
func (c *nameServerClient) GetBaseUrl() (string, error) {
	if c.opt.Hostname == "" {
		return "", fmt.Errorf("hostname is required")
	}
	if c.opt.Port == 0 {
		return "", fmt.Errorf("port is required")
	}

	if c.opt.UseTLS {
		return fmt.Sprintf("https://%s:%d", c.opt.Hostname, c.opt.Port), nil
	} else {
		return fmt.Sprintf("http://%s:%d", c.opt.Hostname, c.opt.Port), nil
	}
}

func (c *nameServerClient) GetAPIUrl(keys ...string) (string, error) {
	baseUrl, err := c.GetBaseUrl()
	if err != nil {
		return "", err
	}
	for _, key := range keys {
		if strings.HasPrefix(key, "/") {
			baseUrl += key
		} else {
			baseUrl += "/" + key
		}
	}
	return baseUrl, nil
}

func (c *nameServerClient) Ping() error {
	url, err := c.GetAPIUrl(server.APILayout.Ping)
	if err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)

	var result ping.PingResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	} else {
		if result.Message == "pong" {
			return nil
		} else {
			return errors.New("response not pong")
		}
	}
}

func (c *nameServerClient) Opt() *config.ClientOpt {
	return &c.opt
}
