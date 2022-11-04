package server

import "go-dfs-server/pkg/config"

var GlobalServerOpt *config.NameserverOpt

const NameserverLoginPath = "/auth/login/"
const NameserverTokenRefreshPath = "/auth/refresh"
const NameserverPingPath = "/ping"
const NameserverHeartBeatPath = "/heartbeat"
