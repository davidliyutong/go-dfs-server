package v1

import (
	"errors"
	"fmt"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/utils"
	"io"
	"mime/multipart"
	"os"
	"sync"
)

func (b blobService) Write(sessionID string, syncWrite bool, file *multipart.FileHeader) (int64, error) {
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

		wg := new(sync.WaitGroup)
		errChan := make(chan error, 16)
		errs := make([]error, 0)

		var numBytes int64 = 0

		wg.Add(1)
		go func() {
			for {
				err, ok := <-errChan
				if !ok {
					wg.Done()
					return
				} else {
					errs = append(errs, err)
				}
			}

		}()

		for {
			buf := make([]byte, v1.DefaultBlobChunkSize)
			n, err := src.Read(buf)
			if err != nil && err != io.EOF {
				errChan <- err
				close(errChan)
				break
			}
			if err == io.EOF {
				close(errChan)
				break
			}

			wg.Add(1)
			if err = session.Write(buf, int64(n), wg, errChan); err == nil {
				numBytes += int64(n)
			} else {
				errChan <- err
				close(errChan)
				break
			}
		}

		if syncWrite {
			wg.Wait()
			if utils.HasError(errs) {
				return numBytes, errors.New(fmt.Sprintf("write command failed, %v", utils.GetFirstError(errs)))
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
