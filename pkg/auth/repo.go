package auth

import "go-dfs-server/pkg/nameserver/server"

func RepoAuthnBasic(accessKey string, secretKey string) bool {
	return accessKey == server.GlobalServerOpt.Auth.AccessKey && secretKey == server.GlobalServerOpt.Auth.SecretKey
}

func RepoAuthzBasic(accessKey string) bool {
	return accessKey == server.GlobalServerOpt.Auth.AccessKey
}
