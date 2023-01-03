package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/controller"
	v12 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/utils"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type handle struct {
	path       string
	mode       int
	time       time.Time
	opened     bool
	offset     int64
	blob       v12.BlobMetaData
	syncMutex  *sync.RWMutex
	rwMutex    *sync.RWMutex
	eventGroup *sync.WaitGroup
	client     *nameServerClient
}

type Handle interface {
	Close() error
	Flush() error
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Truncate(size int64) error
	Seek(offset int64, whence int) (int64, error)

	Blob() *v12.BlobMetaData
	Opened() *bool
	Client() *nameServerClient
	Path() *string
	Mode() *int
	Time() *time.Time
	Offset() *int64

	Wait()
	Done()
	Add(delta int)
}

func NewHandle(path string, mode int, client *nameServerClient) Handle {
	return &handle{
		path:       path,
		mode:       mode,
		time:       time.Now(),
		opened:     false,
		offset:     0,
		blob:       v12.BlobMetaData{},
		syncMutex:  new(sync.RWMutex),
		rwMutex:    new(sync.RWMutex),
		eventGroup: new(sync.WaitGroup),
		client:     client,
	}
}

func (h *handle) Blob() *v12.BlobMetaData {
	return &h.blob
}

func (h *handle) Opened() *bool {
	return &h.opened
}

func (h *handle) Client() *nameServerClient {
	return h.client
}

func (h *handle) Path() *string {
	return &h.path
}

func (h *handle) Mode() *int {
	return &h.mode
}

func (h *handle) Time() *time.Time {
	return &h.time
}

func (h *handle) Offset() *int64 {
	return &h.offset
}

func (h *handle) Wait() {
	h.eventGroup.Wait()
}

func (h *handle) Done() {
	h.eventGroup.Done()
}

func (h *handle) Add(delta int) {
	h.eventGroup.Add(delta)
}

func (h *handle) Close() error {
	h.rwMutex.Lock()
	h.opened = false
	h.rwMutex.Unlock()
	err := h.Flush()
	return err
}

func (h *handle) IsEmpty() bool {
	return h.blob.Size == int64(0)
}

func (h *handle) sync() error {
	h.syncMutex.Lock()
	defer h.syncMutex.Unlock()

	h.blob.Version += int64(1)
	targetUrl, err := h.client.GetAPIUrl(server.APILayout.V1.Blob.Self, server.APILayout.V1.Blob.File)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(v1.SyncRequest{
		Path: h.path,
		Blob: h.blob,
	})
	if err != nil {
		return err
	}
	req, _ := http.NewRequest("POST", targetUrl, bytes.NewReader(payload))
	if h.client.AuthIsEnabled() {
		req.Header.Add("Authorization", "Bearer "+h.client.opt.Token)
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
	var result v1.SyncResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	} else {
		if result.Code != http.StatusOK {
			return errors.New(result.Msg)
		} else {
			h.blob = result.Blob
			return nil
		}
	}
}
func (h *handle) Flush() error {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()
	return h.sync()
}

func (h *handle) Read(buf []byte) (int, error) {

	var numBytesRead = 0
	var numBytesLeft = len(buf)

	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()
	currOffset := h.offset
	defer func() {
		h.offset = currOffset
	}()

	for numBytesLeft > 0 {
		currChunkID := utils.GetChunkID(currOffset)
		currChunkOffset := utils.GetChunkOffset(currOffset)
		boundary := utils.MinInt64(h.blob.Size, (currChunkID+1)*v12.DefaultBlobChunkSize)
		var numBytesToRead int
		if currOffset+int64(numBytesLeft) >= boundary {
			numBytesToRead = int(boundary - currOffset)
		} else {
			numBytesToRead = numBytesLeft
		}

		reader, err := h.client.blobRead(h.path, currChunkID, currChunkOffset, int64(numBytesToRead))
		if err != nil {
			return numBytesRead, err
		}

		n, err := reader.Read(buf[numBytesRead : numBytesRead+numBytesToRead])
		if err != nil && err != io.EOF {
			return numBytesRead, err
		}
		if n <= 0 && err == io.EOF {
			return numBytesRead, io.EOF
		}
		if n != numBytesToRead {
			return numBytesRead, errors.New("read less bytes than expected")
		}
		numBytesRead += n
		numBytesLeft -= n
		currOffset += int64(numBytesToRead)
	}
	return numBytesRead, nil
}

func (h *handle) Write(buf []byte) (int, error) {

	var numBytesWritten = 0
	var numBytesLeft = len(buf)

	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	if h.mode&os.O_RDONLY != 0 {
		return 0, errors.New("file is opened in read-only mode")
	}

	currOffset := h.offset
	defer func() {
		h.offset = currOffset
		h.blob.Size = utils.MaxInt64(h.blob.Size, currOffset)
	}()

	if h.IsEmpty() {
		h.blob.ExtendTo(0)
		err := h.sync()
		if err != nil {
			return numBytesWritten, err
		}
	}

	for numBytesLeft > 0 {
		currChunkID := utils.GetChunkID(currOffset)
		currChunkOffset := utils.GetChunkOffset(currOffset)
		boundary := (currChunkID + 1) * v12.DefaultBlobChunkSize
		var numBytesToWrite int
		if currOffset+int64(numBytesLeft) >= boundary {
			numBytesToWrite = int(boundary - currOffset)
			h.blob.ExtendTo(currChunkID + 1)
			h.blob.Size = utils.MaxInt64(h.blob.Size, currOffset+int64(numBytesToWrite))
			err := h.sync()
			if err != nil {
				return numBytesWritten, err
			}
		} else {
			numBytesToWrite = numBytesLeft
		}

		h.blob.Versions[currChunkID] += 1
		checksum, n, err := h.client.blobWrite(
			h.path,
			currChunkID,
			currChunkOffset,
			int64(numBytesToWrite),
			h.blob.Versions[currChunkID],
			bytes.NewBuffer(buf[numBytesWritten:numBytesWritten+numBytesToWrite]))
		if err != nil {
			return numBytesWritten, err
		}

		if n != numBytesToWrite {
			return numBytesWritten, errors.New("write less bytes than expected")
		}
		h.blob.ChunkChecksums[currChunkID] = checksum
		numBytesWritten += n
		numBytesLeft -= n
		currOffset += int64(n)
	}
	return numBytesWritten, nil
}

func (h *handle) Truncate(size int64) error {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()
	if h.offset > size {
		h.offset = size
	}
	h.blob.Size = size
	lastChunkID := utils.GetChunkID(size)
	lastChunkOffset := utils.GetChunkOffset(size)

	h.blob.TruncateTo(lastChunkID)
	h.blob.Versions[lastChunkID] += 1
	checksum, n, err := h.client.blobWrite(h.path, lastChunkID, lastChunkOffset, 0, h.blob.Versions[lastChunkID], bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	if n != 0 {
		return errors.New("remote truncate failed")
	}
	h.blob.ChunkChecksums[lastChunkID] = checksum
	return nil

}

func (h *handle) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64
	if whence == io.SeekStart {
		newOffset = offset
	} else if whence == io.SeekCurrent {
		newOffset += offset
	} else if whence == io.SeekEnd {
		newOffset = h.blob.Size - 1 - offset
	}

	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()
	if newOffset < 0 || newOffset > h.blob.Size {
		return 0, errors.New("invalid offset")
	} else {
		h.offset = newOffset
		return newOffset, nil
	}
}
