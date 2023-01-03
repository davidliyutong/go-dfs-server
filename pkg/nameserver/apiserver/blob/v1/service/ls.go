package v1

import (
	model "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
)

func (b blobService) Ls(path string) ([]model.BlobMetaData, error) {
	res, err := b.repo.BlobRepo().Ls(path)
	if err != nil {
		// TODO: handle error
	}
	return res, err
}
