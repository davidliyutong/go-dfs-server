package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-dfs-server/pkg/dataserver/server"
	ping "go-dfs-server/pkg/ping/v1"
	"io"
	"net/http"
	"strings"
)

type dataServerClient struct {
	UUID     string
	Hostname string
	Port     int64
	UseTLS   bool
}

type DataServerClient interface {
	GetBaseUrl() (string, error)
	GetAPIUrl(keys ...string) (string, error)
	GetUUID() string
	SetUUID(string)
	Ping() error
	BlobCreateChunk(path string, id int64) error
	BlobCreateFile(path string) error
	BlobCreateDirectory(path string) error
	BlobDeleteChunk(path string, id int64) error
	BlobDeleteFile(path string) error
	BlobDeleteDirectory(path string) error
	BlobLockFile(path string, session string) error
	BlobReadChunk(path string, id int64, offset int64, size int64) (io.ReadCloser, error)
	BlobReadFileLock(path string) ([]string, error)
	BlobReadFileMeta(path string) (map[int64]int64, map[int64]string, error)
	BlobReadChunkMeta(path string, id int64) (int64, string, error)
	BlobUnlockFile(path string) error
	BlobWriteChunk(path string, id int64, offset int64, size int64, version int64, data io.Reader) (string, int64, error)
	SysRole() (string, error)
	SysVolume() (string, error)
	SysUUID() (string, error)
	SysRegister(uuid string) (string, error)
}

var _ DataServerClient = &dataServerClient{}

func NewDataServerClient(uuid string, hostname string, port int64, useTLS bool) DataServerClient {
	return &dataServerClient{
		UUID:     uuid,
		Hostname: hostname,
		Port:     port,
		UseTLS:   useTLS,
	}
}

func (c *dataServerClient) GetBaseUrl() (string, error) {
	if c.Hostname == "" {
		return "", fmt.Errorf("Hostname is required")
	}
	if c.Port == 0 {
		return "", fmt.Errorf("Port is required")
	}

	if c.UseTLS {
		return fmt.Sprintf("https://%s:%d", c.Hostname, c.Port), nil
	} else {
		return fmt.Sprintf("http://%s:%d", c.Hostname, c.Port), nil
	}
}

func (c *dataServerClient) GetAPIUrl(keys ...string) (string, error) {
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

func (c *dataServerClient) GetUUID() string {
	return c.UUID
}

func (c *dataServerClient) Ping() error {
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

func (c *dataServerClient) SetUUID(uuid string) {
	c.UUID = uuid
}
