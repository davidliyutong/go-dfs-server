package v1

import (
	"encoding/json"
	"errors"
	"go-dfs-server/pkg/config"
	v13 "go-dfs-server/pkg/nameserver/apiserver/sys/v1/controller"
	"go-dfs-server/pkg/nameserver/server"
	"io"
	"net/http"
)

func (c *nameServerClient) SysInfo() (v13.InfoResponse, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Sys.Self, server.APILayout.V1.Sys.Info)
	if err != nil {
		return v13.InfoResponse{}, err
	}

	req, _ := http.NewRequest("GET", targetUrl, nil)

	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return v13.InfoResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v13.InfoResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return v13.InfoResponse{}, err
	} else {
		if result.Code != http.StatusOK {
			return v13.InfoResponse{}, errors.New(result.Msg)
		} else {
			return result, nil
		}
	}

}

func (c *nameServerClient) SysSession(sessionID string) (v13.GetSessionResponse, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Session)
	if err != nil {
		return v13.GetSessionResponse{}, err
	}

	req, _ := http.NewRequest("GET", targetUrl, nil)
	q := req.URL.Query()

	q.Add("session", sessionID)

	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return v13.GetSessionResponse{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v13.GetSessionResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return v13.GetSessionResponse{}, err
	} else {
		if result.Code != http.StatusOK {
			return v13.GetSessionResponse{}, errors.New(result.Msg)
		} else {
			return result, nil
		}
	}

}

func (c *nameServerClient) SysSessions() ([]string, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Sys.Self, server.APILayout.V1.Sys.Sessions)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", targetUrl, nil)

	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v13.GetSessionsResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	} else {
		if result.Code != http.StatusOK {
			return nil, errors.New(result.Msg)
		} else {
			return result.Sessions, nil
		}
	}

}

func (c *nameServerClient) SysServers() ([]config.RegisteredDataServer, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Sys.Self, server.APILayout.V1.Sys.Servers)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", targetUrl, nil)

	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v13.GetServersResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	} else {
		if result.Code != http.StatusOK {
			return nil, errors.New(result.Msg)
		} else {
			return result.Servers, nil
		}
	}

}
