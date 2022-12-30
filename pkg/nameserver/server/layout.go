package server

import (
	"go-dfs-server/pkg/config"
)

var GlobalServerDesc *config.NameServerDesc
var BlobDataServerManger DataServerManager
var BlobSessionManager SessionManager
var BlobLockManager LockManager

const NameServerAPIVersion = "v1"
const NameServerNumOfReplicas = 3

type NameServerLayoutAuth struct {
	Self    string
	Login   string
	Refresh string
}

type NameServerLayoutV1 struct {
	Self string
	Blob string
	Sys  string
}

type NameServerLayoutRoot struct {
	Self string
	Auth NameServerLayoutAuth
	Ping string
	V1   NameServerLayoutV1
}

type NameServerLayoutBlob struct {
}

type NameServerLayoutSys struct {
	Self string
	Info string
}

var APILayout = NameServerLayoutRoot{
	Self: "/",
	Auth: NameServerLayoutAuth{
		Self:    "/auth",
		Login:   "login",
		Refresh: "refresh",
	},
	Ping: "/ping",
	V1: NameServerLayoutV1{
		Self: "/v1",
		Blob: "/v1/blob",
		Sys:  "/v1/sys",
	},
}
