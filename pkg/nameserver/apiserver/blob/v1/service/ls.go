package v1

import (
	model "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/utils"
	"os"
	"path/filepath"
)

func (b blobService) Ls(path string) ([]model.BlobMetaData, error) {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(filePath)
	if err != nil {
		return nil, err
	}
	isFile := utils.GetFileState(filePath)
	if isFile {
		res := make([]model.BlobMetaData, 1)
		res[0] = model.BlobMetaData{
			BaseName: filepath.Base(filePath),
			Size:     server.NewSession(path, filePath, "", os.O_RDONLY).GetBlobMetaData().Size,
			Type:     model.BlobFileTypeName,
		}
		return res, nil
	} else {
		lst, err := os.ReadDir(filePath)
		if err != nil {
			return nil, err
		}
		res := make([]model.BlobMetaData, 0)
		for _, val := range lst {
			if !utils.GetFileState(filepath.Join(filePath, val.Name())) {
				res = append(res, model.BlobMetaData{
					BaseName: val.Name(),
					Type:     model.BlobDirTypeName,
				})
			} else {
				res = append(res, model.BlobMetaData{
					BaseName: val.Name(),
					Size:     server.NewSession(path, filePath, "", os.O_RDONLY).GetBlobMetaData().Size,
					Type:     model.BlobFileTypeName,
				})
			}
		}
		return res, nil
	}
}
