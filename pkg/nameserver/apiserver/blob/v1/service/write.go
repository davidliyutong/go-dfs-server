package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"io"
	"mime/multipart"
	"os"
)

func (b blobService) Write(sessionID string, c *gin.Context, file *multipart.FileHeader) (int64, error) {
	session, err := b.repo.BlobRepo().SessionManager().Get(sessionID)
	if err != nil {
		return -1, err
	}
	mode := *session.GetMode()
	switch mode {
	case os.O_RDONLY:
		return 0, errors.New("file is read only")
	case os.O_RDWR:
		src, err := file.Open()
		if err != nil {
			return 0, err
		}
		buf := make([]byte, v1.DefaultBlobChunkSize)
		var numBytes int64 = 0
		for {
			n, err := src.Read(buf)
			if err != nil && err != io.EOF {
				return numBytes, err
			}
			if err == io.EOF {
				break
			}
			if err = session.Write(buf, int64(n)); err == nil {
				numBytes += int64(n)
			} else {
				break
			}
		}
		if file.Size != numBytes {
			return numBytes, errors.New("file size not match")
		} else {
			return numBytes, nil
		}
	default:
		return 0, errors.New("invalid mode")
	}
}
