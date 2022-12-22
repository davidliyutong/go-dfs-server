package auth

import "go-dfs-server/pkg/nameserver/server"

func RepoAuthnBasic(accessKey string, secretKey string) bool {
	return accessKey == server.GlobalServerDesc.Opt.Auth.AccessKey && secretKey == server.GlobalServerDesc.Opt.Auth.SecretKey
}

func RepoAuthzBasic(accessKey string) bool {
	return accessKey == server.GlobalServerDesc.Opt.Auth.AccessKey
}
