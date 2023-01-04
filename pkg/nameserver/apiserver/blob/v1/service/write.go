package v1

import (
	"errors"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"mime/multipart"
)

func (b blobService) Write(path string, chunkID int64, chunkOffset int64, size int64, version int64, file *multipart.FileHeader) ([]string, int64, error) {

	if chunkOffset < 0 || chunkOffset >= v1.DefaultBlobChunkSize {
		return nil, 0, errors.New("invalid chunk offset")
	} else if size <= 0 {
		size = v1.DefaultBlobChunkSize - chunkOffset
	}
	if size+chunkOffset > v1.DefaultBlobChunkSize {
		return nil, 0, errors.New("chunk size is too large")
	}
	src, err := file.Open()
	if err != nil {
		return nil, 0, err
	}
	MD5String, written, err := b.repo.BlobRepo().Write(path, chunkID, chunkOffset, size, version, src)
	if err != nil {
		return MD5String, written, err
	}
	if written != file.Size {
		return MD5String, written, errors.New("file size not match")
	}

	return MD5String, written, nil
}
