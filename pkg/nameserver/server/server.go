package server

import "go-dfs-server/pkg/config"

var GlobalServerOpt *config.NameserverOpt

const NameserverAPIVersion = "v1"
const NameserverAPIPrefix = "/" + NameserverAPIVersion

const NameserverLoginPath = "/auth/login"
const NameserverTokenRefreshPath = "/auth/refresh"
const NameserverPingPath = "/ping"

// Versioned

const NameserverInfoPath = "/info"
