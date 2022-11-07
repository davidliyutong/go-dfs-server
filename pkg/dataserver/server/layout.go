package server

import "go-dfs-server/pkg/config"

var GlobalServerOpt *config.NameserverOpt

const DataserverAPIVersion = "v1"

type DataserverLayoutAuth struct {
	Self    string
	Login   string
	Refresh string
}

type DataserverLayoutV1 struct {
	Self string
	Info string
}

type DataserverLayoutRoot struct {
	Self string
	Auth DataserverLayoutAuth
	Ping string
	Info string
	V1   DataserverLayoutV1
}

var APILayout = DataserverLayoutRoot{
	Self: "/",
	Auth: DataserverLayoutAuth{
		Self:    "/auth",
		Login:   "/auth/login",
		Refresh: "/auth/refresh",
	},
	Ping: "/ping",
	Info: "/info",
	V1: DataserverLayoutV1{
		Self: "/v1",
		Info: "/v1/info",
	},
}
