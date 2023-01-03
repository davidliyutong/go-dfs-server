package v1

import (
	model "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
)

func (b blobService) Open(path string, mode int) (model.BlobMetaData, error) {

	return b.repo.BlobRepo().Open(path, mode)

}
