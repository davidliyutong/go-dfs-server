package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	sys "go-dfs-server/pkg/dataserver/apiserver/sys/v1/controller"
	"go-dfs-server/pkg/dataserver/server"
	"io"
	"net/http"
)

func (c *dataServerClient) SysRole() (string, error) {
	url, err := c.GetAPIUrl(server.APILayout.V1.Sys.Self, server.APILayout.V1.Sys.Info)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)

	var result sys.InfoResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	} else {
		if result.Code != http.StatusOK {
			return "", errors.New(result.Msg)
		} else {
			return result.Role, nil
		}
	}
}

func (c *dataServerClient) SysVolume() (string, error) {
	url, err := c.GetAPIUrl(server.APILayout.V1.Sys.Self, server.APILayout.V1.Sys.Config)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)

	var result sys.ConfigResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	} else {
		if result.Code != http.StatusOK {
			return "", errors.New(result.Msg)
		} else {
			return result.Volume, nil
		}
	}
}
func (c *dataServerClient) SysUUID() (string, error) {
	url, err := c.GetAPIUrl(server.APILayout.V1.Sys.Self, server.APILayout.V1.Sys.UUID)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)

	var result sys.UUIDResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	} else {
		if result.Code != http.StatusOK {
			return "", errors.New(result.Msg)
		} else {
			return result.UUID, nil
		}
	}
}

func (c *dataServerClient) SysRegister(uuid string) (string, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Sys.Self, server.APILayout.V1.Sys.Register)
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(sys.RegisterRequest{
		UUID: uuid,
	})
	if err != nil {
		return "", err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
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
	var result sys.RegisterResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	} else {
		if result.Code != http.StatusOK {
			return "", errors.New(result.Msg)
		} else {
			return result.UUID, nil
		}
	}
}
