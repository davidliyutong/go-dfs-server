package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	blob "go-dfs-server/pkg/dataserver/apiserver/blob/v1/controller"
	"go-dfs-server/pkg/dataserver/server"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
)

func (c *dataServerClient) BlobCreateChunk(path string, id int64) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.CreateChunk)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(blob.CreateChunkRequest{
		Path: path,
		ID:   id,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
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
	var result blob.CreateChunkResponse
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

func (c *dataServerClient) BlobCreateFile(path string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.CreateFile)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(blob.CreateFileRequest{
		Path: path,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
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
	var result blob.CreateFileResponse
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

func (c *dataServerClient) BlobCreateDirectory(path string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.CreateDirectory)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(blob.CreateDirectoryRequest{
		Path: path,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
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
	var result blob.CreateDirectoryResponse
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

func (c *dataServerClient) BlobDeleteChunk(path string, id int64) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.DeleteChunk)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(blob.DeleteChunkRequest{
		Path: path,
		ID:   id,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
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
	var result blob.DeleteChunkResponse
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

func (c *dataServerClient) BlobDeleteFile(path string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.DeleteFile)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(blob.DeleteFileRequest{
		Path: path,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
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
	var result blob.DeleteFileResponse
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

func (c *dataServerClient) BlobDeleteDirectory(path string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.DeleteDirectory)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(blob.DeleteDirectoryRequest{
		Path: path,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
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
	var result blob.DeleteDirectoryResponse
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

func (c *dataServerClient) BlobLockFile(path string, session string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.LockFile)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(blob.LockFileRequest{
		Path:    path,
		Session: session,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
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
	var result blob.LockFileResponse
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

func (c *dataServerClient) BlobReadChunk(path string, id int64) (io.ReadCloser, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.ReadChunk)
	if err != nil {
		return nil, err
	}

	payload := url.Values{
		"path": {path},
		"id":   {strconv.FormatInt(id, 10)},
	}.Encode()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", targetUrl, payload), nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (c *dataServerClient) BlobReadFileLock(path string) ([]string, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.ReadFileLock)
	if err != nil {
		return nil, err
	}

	payload := url.Values{
		"path": {path},
	}.Encode()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", targetUrl, payload), nil)

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
	var result blob.ReadFileLockResponse
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

func (c *dataServerClient) BlobReadFileMeta(path string) (map[int64]int64, map[int64]string, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.ReadFileMeta)
	if err != nil {
		return nil, nil, err
	}

	payload := url.Values{
		"path": {path},
	}.Encode()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", targetUrl, payload), nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result blob.ReadFileMetaResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, nil, err
	} else {
		if result.Code != http.StatusOK {
			return nil, nil, errors.New(result.Msg)
		} else {
			return result.Versions, result.Checksums, nil
		}
	}
}

func (c *dataServerClient) BlobReadChunkMeta(path string, id int64) (int64, string, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.ReadChunkMeta)
	if err != nil {
		return -1, "", err
	}

	payload := url.Values{
		"path": {path},
		"id":   {strconv.FormatInt(id, 10)},
	}.Encode()

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s?%s", targetUrl, payload), nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return -1, "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, _ := io.ReadAll(resp.Body)
	var result blob.ReadChunkMetaResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return -1, "", err
	} else {
		if result.Code != http.StatusOK {
			return -1, "", errors.New(result.Msg)
		} else {
			return result.Version, result.Checksum, nil
		}
	}
}

func (c *dataServerClient) BlobUnlockFile(path string) error {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.UnlockFile)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(blob.UnlockFileRequest{
		Path: path,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
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
	var result blob.UnlockFileResponse
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

func (c *dataServerClient) BlobWriteChunk(path string, id int64, version int64, data io.Reader) (string, error) {
	targetUrl, err := c.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.WriteChunk)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	bw := multipart.NewWriter(buf)
	w1, _ := bw.CreateFormField("path")
	_, err = w1.Write([]byte(path))
	if err != nil {
		return "", err
	}

	w2, _ := bw.CreateFormField("id")
	_, err = w2.Write([]byte(strconv.FormatInt(id, 10)))
	if err != nil {
		return "", err
	}

	w3, _ := bw.CreateFormField("version")
	_, err = w3.Write([]byte(strconv.FormatInt(version, 10)))
	if err != nil {
		return "", err
	}

	w4, _ := bw.CreateFormFile("file", fmt.Sprintf("%d.dat", id))
	_, _ = io.Copy(w4, data)

	err = bw.Close()
	if err != nil {
		return "", err
	}

	req, _ := http.NewRequest("PUT", targetUrl, buf)
	req.Header.Add("Content-Type", bw.FormDataContentType())

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
	var result blob.WriteChunkResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	} else {
		if result.Code != http.StatusOK {
			return "", errors.New(result.Msg)
		} else {
			return result.Checksum, nil
		}
	}
}
