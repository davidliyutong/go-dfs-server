package cmd

import (
	"github.com/spf13/cobra"
	client "go-dfs-server/pkg/client"
)

var rootCmd = &cobra.Command{
	Use:   "go-dfs-client",
	Short: "go-dfs-client, dfs client of a distributed file system",
	Long:  "go-dfs-client, dfs client of a distributed file system",
	Run:   client.MainLoop,
}

var authCmd = &cobra.Command{
	Use: "auth",
	Short: `auth connects client to dfs nameserver and store connection configuration.
If --print flag is present, the auth information will be printed to stdout.
If '--output / -o flag is present, the auth information will be saved to the path specified
Otherwise init will output configuration file to $HOME/.config/go-dfs-server/auth.yaml
`,
	Example: `  go-dfs-client auth --print
  go-dfs-client init --output /path/to/auth.yaml
  go-dfs-client init -o /path/to/auth.yaml`,
	Run: client.Authenticate,
}

func getRootCmd() *cobra.Command {
	rootCmd.AddCommand(authCmd)
	return rootCmd
}

func Execute() {
	rootCmd := getRootCmd()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
