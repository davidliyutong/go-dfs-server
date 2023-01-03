package v1

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"io"
)

func (b blobService) Read(path string, chunkID int64, chunkOffset int64, size int64, c *gin.Context) (int64, error) {

	if chunkOffset < 0 || chunkOffset >= v1.DefaultBlobChunkSize {
		return 0, errors.New("invalid chunk offset")
	} else if size < 0 {
		size = v1.DefaultBlobChunkSize - chunkOffset
	}
	if size+chunkOffset > v1.DefaultBlobChunkSize {
		return 0, errors.New("chunk size is too large")
	}

	buf := bytes.NewBuffer(nil)
	n, err := b.repo.BlobRepo().Read(buf, path, chunkID, chunkOffset, size)
	if n != size {
		return n, errors.New("read size is not equal to size")
	}
	if err != nil {
		return n, err
	}

	_, _ = io.Copy(c.Writer, buf)
	return n, nil

}
