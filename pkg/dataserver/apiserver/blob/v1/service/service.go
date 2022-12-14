package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "go-dfs-server/pkg/dataserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/dataserver/apiserver/blob/v1/repo"
	"go-dfs-server/pkg/dataserver/server"
	"go-dfs-server/pkg/status"
	"go-dfs-server/pkg/utils"
	"io"
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
	LockFile(path string, id string) error
	ReadChunk(path string, id int64, offset int64, size int64, c *gin.Context) error
	ReadChunkMeta(path string, id int64) (int64, string, error)
	ReadFileLock(path string) ([]string, error)
	ReadFileMeta(path string) (v1.BlobMetaData, error)
	UnlockFile(path string) error
	WriteChunk(path string, id int64, offset int64, size int64, version int64, file *multipart.FileHeader) (string, int64, error)
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
		return status.ErrFileNotFile
	}
}

func (b *blobService) CreateDirectory(path string) error {
	directoryPath := filepath.Join(server.GlobalServerDesc.Opt.Volume, path)
	_, err := os.Stat(directoryPath)
	if err == nil {
		return status.ErrFileExists
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
		return status.ErrFileExists
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
				return status.ErrChunkCannotRemove
			} else {
				return nil
			}
		} else if os.IsNotExist(err) {
			return errors.New(fmt.Sprintf("chunk %d does not exist", id))
		} else {
			return err
		}
	} else {
		return status.ErrFileNotFile
	}

}

func (b *blobService) DeleteDirectory(path string) error {
	directoryPath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return err
	}

	_, err = os.Stat(directoryPath)
	if err == nil {
		if utils.IsSameDirectory(server.GlobalServerDesc.Opt.Volume, directoryPath) {
			dir, err := os.ReadDir(directoryPath)
			if err != nil {
				return err
			}
			for _, d := range dir {
				err = os.RemoveAll(filepath.Join(directoryPath, d.Name()))
				if err != nil {
					return err
				}
			}
		} else {
			err = os.RemoveAll(directoryPath)
			return err
		}
	} else if os.IsNotExist(err) {
		return status.ErrFileOrDirectoryNotExist
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
				return status.ErrDirectoryCannotRemoveLocal
			} else {
				return nil
			}
		} else {
			return status.ErrFileNotFile
		}
	} else if os.IsNotExist(err) {
		return status.ErrFileOrDirectoryNotExist
	} else {
		return err
	}
}

func (b *blobService) getChunkMD5(path string, id int64) (string, error) {
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
			return "", status.ErrFileNotFile
		}

	} else if os.IsNotExist(err) {
		return "", status.ErrFileOrDirectoryNotExist
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
				return status.ErrFileAlreadyLocked
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
		return status.ErrFileOrDirectoryNotExist
	}
	return err
}

func (b *blobService) ReadChunk(path string, id int64, offset int64, size int64, c *gin.Context) error {
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
				c.Header("ChunkChecksums-Disposition", fmt.Sprintf("attachment; filename=%d.dat", id))
				c.Header("ChunkChecksums-Type", "application/octet-stream")
				c.Header("ChunkChecksums-Transfer-Encoding", "binary")
				c.Header("Cache-Control", "no-cache")

				if offset == 0 && size < 0 {
					c.File(chunkPath)
					return nil
				} else {
					f, _ := os.OpenFile(chunkPath, os.O_RDONLY, 0775)
					_, err = f.Seek(offset, io.SeekStart)
					if err != nil {
						return err
					}
					if size >= 0 {
						_, err = io.CopyN(c.Writer, f, size)
					} else {
						_, err = io.Copy(c.Writer, f)
					}
					return nil
				}
			} else if os.IsNotExist(err) {
				return errors.New(fmt.Sprintf("chunk %d dose not exists", id))
			} else {
				return err
			}
		} else {
			return status.ErrFileNotFile
		}
	} else if os.IsNotExist(err) {
		return status.ErrFileNotExist
	} else {
		return err
	}

}

func (b *blobService) ReadChunkMeta(path string, id int64) (int64, string, error) {
	meta, err := b.ReadFileMeta(path)
	if err != nil {
		return -1, "", err
	} else {
		MD5String, ok1 := meta.ChunkChecksums[id]
		version, ok2 := meta.Versions[id]
		if !ok1 || !ok2 {
			return -1, "", status.ErrMetaNotFound
		} else {
			return version, MD5String, nil
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
		return nil, status.ErrFileOrDirectoryNotExist
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
				return meta, status.ErrMetaDataCannotLoad
			} else {
				return meta, nil
			}
		} else {
			return meta, status.ErrFileNotFile
		}
	} else if os.IsNotExist(err) {
		return meta, status.ErrFileOrDirectoryNotExist
	} else {
		return meta, err
	}
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
			return status.ErrFileNotLocked
		}
	} else if os.IsNotExist(err) {
		return status.ErrFileOrDirectoryNotExist
	} else {
		return err
	}
}

func (b *blobService) updateMeta(path string, id int64, version int64, data interface{}) error {
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
				return status.ErrMetaDataCannotLoad
			} else {
				meta.Versions[id] = version
				meta.ChunkChecksums[id] = fmt.Sprintf("%v", data)
			}

			err = meta.Dump()
			if err != nil {
				return status.ErrMetaDataCannotDump
			} else {
				return nil
			}
		} else {
			return status.ErrFileNotFile
		}
	} else if os.IsNotExist(err) {
		return status.ErrFileOrDirectoryNotExist
	}
	return err

}

func (b *blobService) WriteChunk(path string, id int64, offset int64, size int64, version int64, file *multipart.FileHeader) (string, int64, error) {
	filePath, err := utils.JoinSubPathSafe(server.GlobalServerDesc.Opt.Volume, path)
	if err != nil {
		return "", 0, err
	}
	_, err = os.Stat(filePath)
	if err == nil {
		isFile := utils.GetFileState(filePath)
		if isFile {
			metaPath := utils.GetMetaPath(filePath)
			meta := v1.NewBlobMetaData(metaPath)
			err = meta.Load()
			if oldVersion, ok := meta.Versions[id]; ok {
				if oldVersion >= version {
					return "", 0, status.ErrMetaVersionConflict
				}
			}

			src, err := file.Open()
			if err != nil {
				return "", 0, err
			}

			chunkPath := utils.GetChunkPath(filePath, id)
			dst, err := os.OpenFile(chunkPath, os.O_RDWR, 0775)
			if err != nil {
				return "", 0, err
			}
			startOffset, err := dst.Seek(offset, io.SeekStart)
			if err != nil {
				return "", 0, err
			}

			var written int64
			if size <= 0 {
				written, err = io.Copy(dst, src)
			} else {
				written, err = io.CopyN(dst, src, size)
			}
			if err != nil {
				return "", written, err
			}
			err = dst.Truncate(startOffset + size)
			if err != nil {
				return "", written, err
			}
			err = dst.Close()
			if err != nil {
				return "", written, err
			}
			MD5String, err := b.getChunkMD5(path, id)
			if err != nil {
				return "", written, err
			}
			return MD5String, written, b.updateMeta(path, id, version, MD5String)

		} else {
			return "", 0, status.ErrFileNotFile
		}
	} else if os.IsNotExist(err) {
		return "", 0, status.ErrFileNotExist
	} else {
		return "", 0, err
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
