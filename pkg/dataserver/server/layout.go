package server

import "go-dfs-server/pkg/config"

var GlobalServerDesc *config.DataServerDesc

var GlobalFileLocks map[string]map[string]bool

const DataServerAPIVersion = "v1"

type DataServerLayoutV1 struct {
	Self string
	Blob string
	Sys  string
}

type DataServerLayoutRoot struct {
	Self string
	Ping string
	Info string
	V1   DataServerLayoutV1
}

type DataServerLayoutBlob struct {
	Self            string
	createChunk     string
	createDirectory string
	createFile      string
	deleteChunk     string
	deleteDirectory string
	deleteFile      string
	lockFile        string
	readChunk       string
	readMeta        string
	unlockFile      string
	writeChunk      string
}

var APILayout = DataServerLayoutRoot{
	Self: "/",
	Ping: "/ping",
	V1: DataServerLayoutV1{
		Self: "/v1",
		Blob: "/v1/blob",
		Sys:  "/v1/sys",
	},
}
