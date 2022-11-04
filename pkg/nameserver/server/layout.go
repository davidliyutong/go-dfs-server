package server

import "go-dfs-server/pkg/config"

var GlobalServerOpt *config.NameserverOpt

const NameserverAPIVersion = "v1"

type NameserverLayoutAuth struct {
	Self    string
	Login   string
	Refresh string
}

type NameserverLayoutV1 struct {
	Self string
	Info string
}

type NameserverLayoutRoot struct {
	Self string
	Auth NameserverLayoutAuth
	Ping string
	Info string
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
	Info: "/info",
	V1: NameserverLayoutV1{
		Self: "/v1",
		Info: "/v1/info",
	},
}
