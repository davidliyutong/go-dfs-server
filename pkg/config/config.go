package config

import (
	"os"
	"path"
)

const ClusterDefaultDomain = "dfs.local"
const ServerDefaultConfigSearchPath0 = "/etc/go-dfs-server"
const ServerDefaultConfigSearchPath1 = "./"

var userHomeDir, _ = os.UserHomeDir()
var ServerDefaultConfigSearchPath2 = path.Join(userHomeDir, ".config/go-dfs-server")

type SeverRoleType string

type AuthOpt struct {
	Domain    string
	AccessKey string
	SecretKey string
}

type LogOpt struct {
	Level string
	Path  string
}
