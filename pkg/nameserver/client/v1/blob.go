package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/controller"
	v12 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/nameserver/server"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (c *nameServerClient) BlobOpen(path string, mode int) (string, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Session)
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("GET", targetUrl, nil)
	q := req.URL.Query()

	q.Add("path", path)
	q.Add("mode", strconv.FormatInt(int64(mode), 10))

	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.URL.RawQuery = q.Encode()
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
	var result v1.OpenResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	} else {
		if result.Code != http.StatusOK {
			return "", errors.New(result.Msg)
		} else {
			return result.Session, nil
		}
	}

}

func (c *nameServerClient) BlobClose(sessionID string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Session)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(v1.CloseRequest{
		Session: sessionID,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("DELETE", targetUrl, bytes.NewReader(payload))
	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v1.CloseResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	} else {
		if result.Code != http.StatusOK {
			return errors.New(result.Msg)
		} else {
			return nil
		}
	}
}

func (c *nameServerClient) BlobFlush(sessionID string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Session)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(v1.CloseRequest{
		Session: sessionID,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v1.CloseResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	} else {
		if result.Code != http.StatusOK {
			return errors.New(result.Msg)
		} else {
			return nil
		}
	}
}

func (c *nameServerClient) BlobRead(sessionID string, size int64) (io.ReadCloser, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.IO)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", targetUrl, nil)
	q := req.URL.Query()

	q.Add("session", sessionID)
	q.Add("size", strconv.FormatInt(int64(size), 10))

	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func (c *nameServerClient) BlobWrite(sessionID string, size int64, data io.Reader, sync bool) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.IO)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)
	w1, _ := bw.CreateFormField("session")
	_, err = w1.Write([]byte(sessionID))
	if err != nil {
		return err
	}

	w2, _ := bw.CreateFormField("sync")
	if sync {
		_, err = w2.Write([]byte("1"))
	} else {
		_, err = w2.Write([]byte("0"))
	}

	if err != nil {
		return err
	}

	w3, _ := bw.CreateFormFile("file", fmt.Sprintf("%v.dat", sessionID))
	if size <= 0 {
		_, _ = io.Copy(w3, data)
	} else {
		_, _ = io.CopyN(w3, data, size)
	}

	err = bw.Close()
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", targetUrl, buf)
	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", bw.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v1.WriteResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	} else {
		if result.Code != http.StatusOK {
			return errors.New(result.Msg)
		} else {
			return nil
		}
	}
}

func (c *nameServerClient) BlobTruncate(sessionID string, size int64) error {
	//TODO implement me
	return errors.New("implement me")
}

func (c *nameServerClient) BlobSeek(sessionID string, offset int64, whence int) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Seek)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(v1.SeekRequest{
		Session: sessionID,
		Offset:  offset,
		Whence:  whence,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v1.SeekResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	} else {
		if result.Code != http.StatusOK {
			return errors.New(result.Msg)
		} else {
			return nil
		}
	}
}

func (c *nameServerClient) BlobLock(sessionID string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Lock)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(v1.LockRequest{
		Session: sessionID,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v1.LockResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	} else {
		if result.Code != http.StatusOK {
			return errors.New(result.Msg)
		} else {
			return nil
		}
	}
}

func (c *nameServerClient) BlobUnlock(sessionID string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Lock)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(v1.UnlockRequest{
		Session: sessionID,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("DELETE", targetUrl, bytes.NewReader(payload))
	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v1.UnlockResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	} else {
		if result.Code != http.StatusOK {
			return errors.New(result.Msg)
		} else {
			return nil
		}
	}
}

func (c *nameServerClient) BlobGetLock(path string) ([]string, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Lock)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", targetUrl, nil)
	q := req.URL.Query()

	q.Add("path", path)

	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.URL.RawQuery = q.Encode()
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
	var result v1.GetLockResponse
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

func (c *nameServerClient) BlobMkdir(path string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Lock)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(v1.MkdirRequest{
		Path: path,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v1.MkdirResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	} else {
		if result.Code != http.StatusOK {
			return errors.New(result.Msg)
		} else {
			return nil
		}
	}
}

func (c *nameServerClient) BlobLs(path string) ([]v1.LsFileInfo, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Path)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", targetUrl, nil)
	q := req.URL.Query()

	q.Add("path", path)

	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.URL.RawQuery = q.Encode()
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
	var result v1.LsResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	} else {
		if result.Code != http.StatusOK {
			return nil, errors.New(result.Msg)
		} else {
			return result.List, nil
		}
	}

}

func (c *nameServerClient) BlobRm(path string, recursive bool) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Path)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(v1.RmRequest{
		Path:      path,
		Recursive: recursive,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("DELETE", targetUrl, bytes.NewReader(payload))
	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v1.RmResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	} else {
		if result.Code != http.StatusOK {
			return errors.New(result.Msg)
		} else {
			return nil
		}
	}
}

func (c *nameServerClient) BlobGetFileMeta(path string) (v12.BlobMetaData, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Meta)
	if err != nil {
		return v12.BlobMetaData{}, err
	}

	req, _ := http.NewRequest("GET", targetUrl, nil)
	q := req.URL.Query()

	q.Add("path", path)

	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return v12.BlobMetaData{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result v1.GetFileMetaResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return v12.BlobMetaData{}, err
	} else {
		if result.Code != http.StatusOK {
			return v12.BlobMetaData{}, errors.New(result.Msg)
		} else {
			return result.Blob, nil
		}
	}

}
