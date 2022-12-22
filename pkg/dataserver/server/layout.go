package server

import "go-dfs-server/pkg/config"

var GlobalServerDesc *config.DataserverDesc

const DataserverAPIVersion = "v1"

type DataserverLayoutV1 struct {
	Self string
	Blob string
	Sys  string
}

type DataserverLayoutRoot struct {
	Self string
	Ping string
	Info string
	V1   DataserverLayoutV1
}

type DataserverLayoutBlob struct {
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

var APILayout = DataserverLayoutRoot{
	Self: "/",
	Ping: "/ping",
	V1: DataserverLayoutV1{
		Self: "/v1",
		Blob: "/v1/blob",
		Sys:  "/v1/sys",
	},
}
