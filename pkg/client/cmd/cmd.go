package cmd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	client "go-dfs-server/pkg/client"
	"go-dfs-server/pkg/client/login"
)

var rootCmd = &cobra.Command{
	Use:   "go-dfs-client",
	Short: "go-dfs-client, dfs client of a distributed file system",
	Long: `go-dfs-client, dfs client of a distributed file system

Use --verbose / -v to enable debug ouput`,
	Run: client.MainLoop,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: `login connects client to dfs nameserver and store connection configuration.`,
	Long: `login connects client to dfs nameserver and store connection configuration.

If --config /path/to/config is specified, client will first try to exam token stored in /path/to/config and refresh the token.

The token obtained from server will be stored in /path/to/config, which is ~/.config/go-dfs-server/client.yaml by default
`,
	Example: `  go-dfs-client login dfs://127.0.0.1
  go-dfs-client login dfs://127.0.0.1 --accessKey=12345678 --secretKey=xxxxxxxx
  go-dfs-client login dfs://127.0.0.1:27903 --accessKey=12345678 --secretKey=xxxxxxxx`,
	Args: cobra.MaximumNArgs(1),
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
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")

	loginCmd.Flags().String("accessKey", "", "server access key")
	loginCmd.Flags().String("secretKey", "", "server secret key")
	loginCmd.Flags().String("config", "", "default configuration path")

	rootCmd.AddCommand(loginCmd)

	rootCmd.AddCommand(logoutCmd)

	rootCmd.AddCommand(lsCmd)

	return rootCmd
}

func Execute() {
	rootCmd := getRootCmd()
	verboseFlag, _ := rootCmd.PersistentFlags().GetBool("verbose")
	if verboseFlag {
		log.SetLevel(log.DebugLevel)
	}

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
