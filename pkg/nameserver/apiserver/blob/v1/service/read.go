package v1

import (
	"github.com/gin-gonic/gin"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/utils"
	"io"
	"math"
)

func (b blobService) Read(sessionID string, size int64, c *gin.Context) (int64, error) {
	session, err := b.repo.BlobRepo().SessionManager().Get(sessionID)
	if err != nil {
		return -1, err
	}
	var bytesToRead int64
	var bytesRead int64 = 0
	if size <= 0 {
		bytesToRead = math.MaxInt64
	} else {
		bytesToRead = size
	}

	var batchSize int64
	var buf = make([]byte, v1.DefaultBlobChunkSize)
	for bytesToRead > 0 {
		batchSize = utils.MinInt64(bytesToRead, v1.DefaultBlobChunkSize)
		n, err := session.Read(buf, batchSize)
		if err != nil && err != io.EOF {
			return bytesRead, err
		} else if err == io.EOF {
			return bytesRead, nil
		}

		_, err = c.Writer.Write(buf[:batchSize])
		if err != nil {
			return bytesRead, err
		}

		bytesRead += n
		bytesToRead -= n

	}
	return bytesRead, nil
}
