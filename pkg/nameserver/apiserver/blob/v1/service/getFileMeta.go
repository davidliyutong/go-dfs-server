package v1

import (
	"errors"
	model "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/utils"
	"os"
)

func (b blobService) GetFileMeta(path string) (model.BlobMetaData, error) {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return model.BlobMetaData{}, err
	}
	_, err = os.Stat(filePath)
	var meta v1.BlobMetaData
	if err == nil {
		isFile := utils.GetFileState(filePath)
		if isFile {
			metaPath := utils.GetMetaPath(filePath)
			err = meta.Load(metaPath)
			if err != nil {
				return meta, errors.New("cannot load metadata")
			} else {
				return meta, nil
			}
		} else {
			return meta, errors.New("not a file")
		}
	} else if os.IsNotExist(err) {
		return meta, errors.New("file or directory does not exist")
	} else {
		return meta, err
	}
}
