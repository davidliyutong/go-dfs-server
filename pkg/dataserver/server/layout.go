package server

import "go-dfs-server/pkg/config"

var GlobalServerDesc *config.DataServerDesc

var GlobalFileLocks map[string]map[string]bool

const DataServerAPIVersion = "v1"

type DataServerLayoutV1 struct {
	Self string
	Blob DataServerLayoutBlob
	Sys  DataServerLayoutSys
}

type DataServerLayoutRoot struct {
	Self string
	Ping string
	V1   DataServerLayoutV1
}

type DataServerLayoutBlob struct {
	Self            string
	CreateChunk     string
	CreateDirectory string
	CreateFile      string
	DeleteChunk     string
	DeleteDirectory string
	DeleteFile      string
	LockFile        string
	ReadChunk       string
	ReadChunkMeta   string
	ReadFileMeta    string
	ReadFileLock    string
	UnlockFile      string
	WriteChunk      string
}

type DataServerLayoutSys struct {
	Self     string
	Info     string
	UUID     string
	Config   string
	Register string
}

var (
	APILayout = DataServerLayoutRoot{
		Self: "/",
		Ping: "/ping",
		V1: DataServerLayoutV1{
			Self: "/v1",
			Blob: DataServerLayoutBlob{
				Self:            "/v1/blob",
				CreateChunk:     "createChunk",
				CreateDirectory: "createDirectory",
				CreateFile:      "createFile",
				DeleteChunk:     "deleteChunk",
				DeleteDirectory: "deleteDirectory",
				DeleteFile:      "deleteFile",
				LockFile:        "lockFile",
				ReadChunk:       "readChunk",
				ReadChunkMeta:   "readChunkMeta",
				ReadFileMeta:    "readFileMeta",
				ReadFileLock:    "readFileLock",
				UnlockFile:      "unlockFile",
				WriteChunk:      "writeChunk",
			},
			Sys: DataServerLayoutSys{
				Self:     "/v1/sys",
				Info:     "info",
				UUID:     "uuid",
				Config:   "config",
				Register: "register",
			},
		},
	}
)
