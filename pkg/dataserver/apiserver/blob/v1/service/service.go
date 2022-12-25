package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "go-dfs-server/pkg/dataserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/dataserver/apiserver/blob/v1/repo"
	"go-dfs-server/pkg/dataserver/server"
	"go-dfs-server/pkg/utils"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type BlobService interface {
	CreateChunk(path string, id int64) error
	CreateDirectory(path string) error
	CreateFile(path string) error
	DeleteChunk(path string, id int64) error
	DeleteDirectory(path string) error
	DeleteFile(path string) error
	GetChunkMD5(path string, id int64) (string, error)
	LockFile(path string, id string) error
	ReadChunk(path string, id int64, c *gin.Context) error
	ReadChunkMeta(path string, id int64) (string, error)
	ReadFileLock(path string) ([]string, error)
	ReadFileMeta(path string) (v1.BlobMetaData, error)
	UnlockFile(path string) error
	UpdateMeta(path string, id int64, data interface{}) error
	WriteChunk(path string, id int64, c *gin.Context, file *multipart.FileHeader) error
}

type blobService struct {
	repo repo.Repo
}

func (b *blobService) CreateChunk(path string, id int64) error {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	isFile := utils.GetFileState(filePath)
	if isFile {
		chunkPath := utils.GetChunkPath(filePath, id)
		if _, err := os.Stat(chunkPath); err == nil {
			return errors.New(fmt.Sprintf("chunk %d exists", id))
		} else if os.IsNotExist(err) {
			_, err := os.Create(chunkPath)
			if err != nil {
				return err
			} else {
				err = os.Chmod(chunkPath, 0775)
				return err
			}
		} else {
			return err
		}
	} else {
		return errors.New("not a file")
	}
}

func (b *blobService) CreateDirectory(path string) error {
	directoryPath := filepath.Join(server.GlobalServerDesc.Opt.Volume, path)
	_, err := os.Stat(directoryPath)
	if err == nil {
		return errors.New("file or directory exists")
	} else if os.IsNotExist(err) {
		err := os.Mkdir(directoryPath, 0775)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
	return err
}

func (b *blobService) CreateFile(path string) error {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		return errors.New("file or directory exists")
	} else if os.IsNotExist(err) {
		err := os.Mkdir(filePath, 0775)
		if err != nil {
			return err
		} else {
			metaPath := utils.GetMetaPath(filePath)
			meta := v1.NewBlobMetaData(metaPath)
			err := meta.Dump()
			if err != nil {
				return err
			} else {
				err = os.Chmod(metaPath, 0775)
				return err
			}
		}
	} else {
		return err
	}
}

func (b *blobService) DeleteChunk(path string, id int64) error {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	isFile := utils.GetFileState(filePath)
	if isFile {
		chunkPath := utils.GetChunkPath(filePath, id)
		_, err := os.Stat(chunkPath)
		if err == nil {
			err = os.Remove(chunkPath)
			if err != nil {
				return errors.New("failed to remove chunk")
			} else {
				return nil
			}
		} else if os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("chunk %d does not exist", id))
		} else {
			return err
		}
	} else {
		return errors.New("not a file")
	}

}

func (b *blobService) DeleteDirectory(path string) error {
	directoryPath := filepath.Join(server.GlobalServerDesc.Opt.Volume, path)
	_, err := os.Stat(directoryPath)
	if err == nil {
		err = os.RemoveAll(directoryPath)
		if err != nil {
			return errors.New("failed to remove directory")
		} else {
			return nil
		}
	} else if os.IsNotExist(err) {
		return errors.New("file or directory does not exist")
	}
	return err
}

func (b *blobService) DeleteFile(path string) error {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		isFile := utils.GetFileState(filePath)
		if isFile {
			err = os.RemoveAll(filePath)
			if err != nil {
				return errors.New("failed to remove directory")
			} else {
				return nil
			}
		} else {
			return errors.New("not a file")
		}
	} else if os.IsNotExist(err) {
		return errors.New("file or directory does not exist")
	} else {
		return err
	}
}

func (b *blobService) GetChunkMD5(path string, id int64) (string, error) {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return "", err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		isFile := utils.GetFileState(filePath)
		if isFile {
			chunkPath := utils.GetChunkPath(filePath, id)
			return utils.GetFileMD5(chunkPath)
		} else {
			return "", errors.New("not a file")
		}

	} else if os.IsNotExist(err) {
		return "", errors.New("file or directory does not exist")
	} else {
		return "", err
	}
}

func (b *blobService) LockFile(path string, id string) error {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		locks, ok := server.GlobalFileLocks[path]
		if ok {
			_, ok := locks[id]
			if ok {
				return errors.New("file already locked by this session")
			} else {
				locks[id] = true
				return nil
			}
		} else {
			locks = make(map[string]bool)
			locks[id] = true
			server.GlobalFileLocks[path] = locks
			return nil
		}
	} else if os.IsNotExist(err) {
		return errors.New("file or directory does not exist")
	}
	return err
}

func (b *blobService) ReadChunk(path string, id int64, c *gin.Context) error {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		isFile := utils.GetFileState(filePath)
		if isFile {
			chunkPath := utils.GetChunkPath(filePath, id)
			_, err := os.Stat(chunkPath)
			if err == nil {
				c.Writer.WriteHeader(http.StatusOK)
				c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%d.dat", id))
				c.Header("Content-Type", "application/octet-stream")
				c.Header("Content-Transfer-Encoding", "binary")
				c.Header("Cache-Control", "no-cache")
				c.File(chunkPath)
				return nil
			} else if os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("chunk %d dose not exists", id))
			} else {
				return err
			}
		} else {
			return errors.New("destination not a file")
		}
	} else if os.IsNotExist(err) {
		return errors.New("file dose not exists")
	} else {
		return err
	}

}

func (b *blobService) ReadChunkMeta(path string, id int64) (string, error) {
	meta, err := b.ReadFileMeta(path)
	if err != nil {
		return "", err
	} else {
		MD5String, ok := meta.Content[id]
		if !ok {
			return "", errors.New("no meta data found for this chunk id")
		} else {
			return MD5String, nil
		}
	}
}

func (b *blobService) ReadFileLock(path string) ([]string, error) {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		locks, ok := server.GlobalFileLocks[path]
		if ok {
			keys := make([]string, 0, len(locks))
			for k := range locks {
				if locks[k] {
					keys = append(keys, k)
				}
			}
			return keys, nil
		} else {
			return nil, nil
		}
	} else if os.IsNotExist(err) {
		return nil, errors.New("file or directory does not exist")
	} else {
		return nil, err
	}

}

func (b *blobService) ReadFileMeta(path string) (v1.BlobMetaData, error) {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return v1.BlobMetaData{}, err
	}
	_, err = os.Stat(filePath)
	var meta v1.BlobMetaData
	if err == nil {
		isFile := utils.GetFileState(filePath)
		if isFile {
			meta.Path = utils.GetMetaPath(filePath)
			err = meta.Load()
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

func (b *blobService) UpdateMeta(path string, id int64, data interface{}) error {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		isFile := utils.GetFileState(filePath)
		if isFile {
			metaPath := utils.GetMetaPath(filePath)
			meta := v1.NewBlobMetaData(metaPath)

			err = meta.Load()
			if err != nil {
				return errors.New("cannot load metadata")
			} else {

				meta.Content[id] = fmt.Sprintf("%v", data)
			}

			err = meta.Dump()
			if err != nil {
				return errors.New("cannot save metadata")
			} else {
				return nil
			}
		} else {
			return errors.New("not a file")
		}
	} else if os.IsNotExist(err) {
		return errors.New("file or directory does not exist")
	}
	return err

}

func (b *blobService) UnlockFile(path string) error {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		_, ok := server.GlobalFileLocks[path]
		if ok {
			delete(server.GlobalFileLocks, path)
			return nil
		} else {
			return errors.New("file not locked")
		}
	} else if os.IsNotExist(err) {
		return errors.New("file or directory does not exist")
	} else {
		return err
	}
}

func (b *blobService) WriteChunk(path string, id int64, c *gin.Context, file *multipart.FileHeader) error {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		isFile := utils.GetFileState(filePath)
		if isFile {
			chunkPath := utils.GetChunkPath(filePath, id)
			if err := c.SaveUploadedFile(file, chunkPath); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			return errors.New("destination not a file")
		}
	} else if os.IsNotExist(err) {
		return errors.New("file dose not exists")
	} else {
		return err
	}
}

var _ BlobService = (*blobService)(nil)

func newBlobService(repo repo.Repo) BlobService {
	return &blobService{repo: repo}
}

type Service interface {
	NewBlobService() BlobService
}

type service struct {
	repo repo.Repo
}

var _ Service = (*service)(nil)

func NewService(repo repo.Repo) Service {
	return &service{repo}
}

func (s *service) NewBlobService() BlobService {
	return newBlobService(s.repo)
}
