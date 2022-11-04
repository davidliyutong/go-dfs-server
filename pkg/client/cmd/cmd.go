package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	client "go-dfs-server/pkg/client"
	"go-dfs-server/pkg/client/login"
)

var rootCmd = &cobra.Command{
	Use:   "go-dfs-client",
	Short: "go-dfs-client, dfs client of a distributed file system",
	Long:  "go-dfs-client, dfs client of a distributed file system",
	Run:   client.MainLoop,
}

var loginCmd = &cobra.Command{
	Use: "login",
	Short: `login connects client to dfs nameserver and store connection configuration.
`,
	Example: `  go-dfs-client login dfs://127.0.0.1
  go-dfs-client login dfs://127.0.0.1 --accessKey=12345678 --secretKey=xxxxxxxx
  go-dfs-client login dfs://127.0.0.1:27903 --accessKey=12345678 --secretKey=xxxxxxxx`,
	Args: cobra.ExactArgs(1),
	Run:  login.Login,
}

var logoutCmd = &cobra.Command{
	Use: "logout",
	Short: `logout clears server credentials.
`,
	Args: cobra.ExactArgs(0),
	Run:  login.Login,
}

var lsCmd = &cobra.Command{
	Use:     "ls",
	Short:   `list a directory`,
	Args:    cobra.ExactArgs(1),
	Example: `  go-dfs-client ls /`,
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		fmt.Printf("exec ls %s", path)
	},
}

var catCmd = &cobra.Command{
	Use:     "ls",
	Short:   `list a directory`,
	Args:    cobra.ExactArgs(1),
	Example: `  go-dfs-client ls /`,
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		fmt.Printf("exec ls %s", path)
	},
}

func getRootCmd() *cobra.Command {
	loginCmd.Flags().String("accessKey", "", "server access key")
	loginCmd.Flags().String("secretKey", "", "server secret key")
	rootCmd.AddCommand(loginCmd)

	rootCmd.AddCommand(logoutCmd)

	rootCmd.AddCommand(lsCmd)

	return rootCmd
}

func Execute() {
	rootCmd := getRootCmd()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
