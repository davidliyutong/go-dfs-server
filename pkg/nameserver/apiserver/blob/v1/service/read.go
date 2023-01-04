package v1

import (
	"bytes"
	"github.com/gin-gonic/gin"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/status"
	"io"
)

func (b blobService) Read(path string, chunkID int64, chunkOffset int64, size int64, c *gin.Context) (int64, error) {

	if chunkOffset < 0 || chunkOffset >= v1.DefaultBlobChunkSize {
		return 0, status.ErrChunkOffsetInvalid
	} else if size < 0 {
		size = v1.DefaultBlobChunkSize - chunkOffset
	}
	if size+chunkOffset > v1.DefaultBlobChunkSize {
		return 0, status.ErrChunkSizeTooLarge
	}

	buf := bytes.NewBuffer(nil)
	n, err := b.repo.BlobRepo().Read(buf, path, chunkID, chunkOffset, size)
	if n != size {
		return n, status.ErrIOReadSizeMismatch
	}
	if err != nil {
		return n, err
	}

	_, _ = io.Copy(c.Writer, buf)
	return n, nil

}
