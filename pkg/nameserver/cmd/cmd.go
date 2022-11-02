package cmd

import (
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/nameserver"
	"go-dfs-server/pkg/nameserver/config"
)

var rootCmd = &cobra.Command{
	Use:   "go-dfs-nameserver",
	Short: "go-dfs-nameserver, name server of a distributed file system",
	Long:  "go-dfs-nameserver, name server of a distributed file system",
}

var serveCmd = &cobra.Command{
	Use: "serve",
	SuggestFor: []string{
		"ru", "ser",
	},
	Short: "serve start the nameserver using predefined configs.",
	Long: `serve start the nameserver using predefined configs, by the following order:

1. path specified in --config flag
2. path defined DFS_CONFIG environment variable
3. default location $HOME/.config/go-dfs-server/nameserver.yaml, /etc/go-dfs-server/nameserver.yaml, current directory

The parameters in the configuration file will be overwritten by the following order:

1. command line arguments
2. environment variables
`,
	Run: nameserver.MainLoop,
}

var initCmd = &cobra.Command{
	Use: "init",
	SuggestFor: []string{
		"ini",
	},
	Short: "init create a configuration template",
	Long: `init create a configuration template. This will generate uuids, secrets and etc. 

The configuration file can be used to launch the nameserver.
If --print flag is present, the configuration will be printed to stdout.
If '--output / -o flag is present, the configuration will be saved to the path specified
Otherwise init will output configuration file to $HOME/.config/go-dfs-server/nameserver.yaml
`,
	Example: `  go-dfs-nameserver init --print
  go-dfs-nameserver init --output /path/to/nameserver.yaml
  go-dfs-nameserver init -o /path/to/nameserver.yaml`,
	Run: config.Init,
}

func getRootCmd() *cobra.Command {

	serveCmd.Flags().String("config", "", "default configuration path")
	serveCmd.Flags().Int16P("port", "p", config.NameserverDefaultPort, "port that nameserver listen on")
	serveCmd.Flags().StringP("interface", "i", config.NameserverDefaultInterface, "interface that nameserver listen on, default to 0.0.0.0")
	serveCmd.Flags().String("volume", config.NameserverDefaultVolume, "default configuration path")
	rootCmd.AddCommand(serveCmd)

	initCmd.Flags().Bool("print", false, "print config to stdout")
	initCmd.Flags().StringP("output", "o", config.NameserverDefaultConfig, "specify output directory")
	rootCmd.AddCommand(initCmd)

	return rootCmd
}

func Execute() {
	rootCmd := getRootCmd()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
