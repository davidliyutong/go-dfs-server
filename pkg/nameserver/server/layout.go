package server

import "go-dfs-server/pkg/config"

var GlobalServerDesc *config.NameserverDesc

const NameserverAPIVersion = "v1"

type NameserverLayoutAuth struct {
	Self    string
	Login   string
	Refresh string
}

type NameserverLayoutV1 struct {
	Self string
	Blob string
	Sys  string
}

type NameserverLayoutRoot struct {
	Self string
	Auth NameserverLayoutAuth
	Ping string
	V1   NameserverLayoutV1
}

var APILayout = NameserverLayoutRoot{
	Self: "/",
	Auth: NameserverLayoutAuth{
		Self:    "/auth",
		Login:   "/auth/login",
		Refresh: "/auth/refresh",
	},
	Ping: "/ping",
	V1: NameserverLayoutV1{
		Self: "/v1",
		Blob: "/v1/blob",
		Sys:  "/v1/sys",
	},
}
