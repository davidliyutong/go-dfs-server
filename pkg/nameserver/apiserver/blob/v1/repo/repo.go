package repo

import (
	model "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"go-dfs-server/pkg/nameserver/server"
	"io"
)

type BlobRepo interface {
	Open(path string, mode int) (model.BlobMetaData, error)                // Sync
	Sync(path string, blob model.BlobMetaData) (model.BlobMetaData, error) // Sync

	Read(buffer io.Writer, path string, chunkID int64, chunkOffset int64, size int64) (int64, error)
	Write(path string, chunkID int64, chunkOffset int64, size int64, version int64, data io.ReadCloser) ([]string, int64, error)
	Rm(path string, recursive bool) error

	Mkdir(path string) error
	Ls(path string) ([]model.BlobMetaData, error)

	SessionManager() server.SessionManager
	DataServerManager() server.DataServerManager
}

type Repo interface {
	BlobRepo() BlobRepo
	Close() error
}
