package v1

import (
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/utils"
)

func (b blobService) GetLock(path string) ([]string, error) {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return nil, err
	}
	return b.repo.BlobRepo().GetLock(filePath)
}
