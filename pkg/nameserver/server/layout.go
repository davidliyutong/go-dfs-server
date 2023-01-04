package server

import (
	"go-dfs-server/pkg/config"
)

var GlobalServerDesc *config.NameServerDesc
var BlobDataServerManger DataServerManager
var BlobSessionManager SessionManager

const NameServerAPIVersion = "v1"
const NameServerNumOfReplicas = 3

type NameServerLayoutAuth struct {
	Self    string
	Login   string
	Refresh string
}
type NameServerLayoutBlob struct {
	Self string
	Path string
	File string
	IO   string
}

type NameServerLayoutSys struct {
	Self     string
	Info     string
	Session  string
	Sessions string
	Servers  string
}

type NameServerLayoutV1 struct {
	Self string
	Blob NameServerLayoutBlob
	Sys  NameServerLayoutSys
}

type NameServerLayoutRoot struct {
	Self string
	Auth NameServerLayoutAuth
	Ping string
	V1   NameServerLayoutV1
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
		Blob: NameServerLayoutBlob{
			Self: "/v1/blob",
			Path: "path",
			File: "file",
			IO:   "io",
		},
		Sys: NameServerLayoutSys{
			Self:     "/v1/sys",
			Info:     "info",
			Session:  "session",
			Sessions: "sessions",
			Servers:  "servers",
		},
	},
}
