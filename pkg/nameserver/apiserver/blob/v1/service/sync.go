package v1

import "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"

func (b blobService) Sync(path string, blob v1.BlobMetaData) (v1.BlobMetaData, error) {
	return b.repo.BlobRepo().Sync(path, blob)
}
