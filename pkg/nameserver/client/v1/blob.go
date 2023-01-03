package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/controller"
	"go-dfs-server/pkg/nameserver/server"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

func (c *nameServerClient) Open(path string, mode int) (Handle, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.File)
	if err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
	} else {
		if result.Code != http.StatusOK {
			return nil, errors.New(result.Msg)
		} else {
			h := NewHandle(path, mode, c)
			*h.Blob() = result.Blob
			*h.Opened() = true
			return h, nil
		}
	}

}

func (c *nameServerClient) blobRead(path string, chunkID int64, chunkOffset int64, size int64) (io.ReadCloser, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.IO)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", targetUrl, nil)
	q := req.URL.Query()

	q.Add("path", path)
	q.Add("chunk_id", strconv.FormatInt(chunkID, 10))
	q.Add("chunk_offset", strconv.FormatInt(chunkOffset, 10))
	q.Add("size", strconv.FormatInt(size, 10))

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

func (c *nameServerClient) blobWrite(path string, chunkID int64, chunkOffset int64, size int64, version int64, data io.Reader) (string, int, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.IO)
	if err != nil {
		return "", 0, err
	}

	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)
	w1, _ := bw.CreateFormField("path")
	_, err = w1.Write([]byte(path))
	if err != nil {
		return "", 0, err
	}

	w2, _ := bw.CreateFormField("chunk_id")
	_, err = w2.Write([]byte(strconv.FormatInt(chunkID, 10)))
	if err != nil {
		return "", 0, err
	}

	w3, _ := bw.CreateFormField("chunk_offset")
	_, err = w3.Write([]byte(strconv.FormatInt(chunkOffset, 10)))
	if err != nil {
		return "", 0, err
	}

	w4, _ := bw.CreateFormField("size")
	_, err = w4.Write([]byte(strconv.FormatInt(size, 10)))
	if err != nil {
		return "", 0, err
	}

	w5, _ := bw.CreateFormField("version")
	_, err = w5.Write([]byte(strconv.FormatInt(version, 10)))
	if err != nil {
		return "", 0, err
	}

	if err != nil {
		return "", 0, err
	}

	w6, _ := bw.CreateFormFile("file", fmt.Sprintf("%v.%v.%v.%v.%v", path, chunkID, chunkOffset, size, version))
	if size <= 0 {
		_, _ = io.Copy(w6, data)
	} else {
		_, _ = io.CopyN(w6, data, size)
	}

	err = bw.Close()
	if err != nil {
		return "", 0, err
	}

	req, _ := http.NewRequest("POST", targetUrl, buf)
	if c.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+c.opt.Token)
	}
	req.Header.Add("Content-Type", bw.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", 0, err
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
		return "", 0, err
	} else {
		if result.Code != http.StatusOK {
			return result.Checksum, result.Written, errors.New(result.Msg)
		} else {
			return result.Checksum, result.Written, nil
		}
	}
}

func (c *nameServerClient) BlobMkdir(path string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.Path)
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
