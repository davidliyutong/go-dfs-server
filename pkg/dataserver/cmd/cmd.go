package cmd

import (
	"github.com/spf13/cobra"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/dataserver/loop"
)

var rootCmd = &cobra.Command{
	Use:   "go-dfs-dataserver",
	Short: "go-dfs-dataserver, data server of a distributed file system",
	Long:  "go-dfs-dataserver, data server of a distributed file system",
}

var serveCmd = &cobra.Command{
	Use: "serve",
	SuggestFor: []string{
		"ru", "ser",
	},
	Short: "serve start the dataserver using predefined configs.",
	Long: `serve start the dataserver using predefined configs, by the following order:

1. path specified in --config flag
2. path defined DFS_CONFIG environment variable
3. default location $HOME/.config/go-dfs-server/dataserver.yaml, /etc/go-dfs-server/dataserver.yaml, current directory

The parameters in the configuration file will be overwritten by the following order:

1. command line arguments
2. environment variables
`,
	Run: loop.MainLoop,
}

var initCmd = &cobra.Command{
	Use: "init",
	SuggestFor: []string{
		"ini",
	},
	Short: "init create a configuration template",
	Long: `init create a configuration template. This will generate uuids, secrets and etc. 

The configuration file can be used to launch the dataserver.
If --print flag is present, the configuration will be printed to stdout.
If '--output / -o flag is present, the configuration will be saved to the path specified
Otherwise init will output configuration file to $HOME/.config/go-dfs-server/dataserver.yaml
If --yes / -y flag is present, the configuration will be overwrite without confirmation
`,
	Example: `  go-dfs-dataserver init --print
  go-dfs-dataserver init --output /path/to/dataserver.yaml
  go-dfs-dataserver init -o /path/to/dataserver.yaml -y`,
	Run: config.InitDataServerCfg,
}

func getRootCmd() *cobra.Command {

	serveCmd.Flags().String("config", "", "default configuration path")
	serveCmd.Flags().Int16P("port", "p", config.DataServerDefaultPort, "port that nameserver listen on")
	serveCmd.Flags().String("remote", "", "url to reach nameserver")
	serveCmd.Flags().String("endpoint", "", "url that nameserver use to  communicate with this dataserver")
	serveCmd.Flags().String("volume", config.DataServerDefaultVolume, "default configuration path")
	serveCmd.Flags().String("domain", config.ClusterDefaultDomain, "domain of DFS cluster, default to dfs.local")
	serveCmd.Flags().String("accessKey", "", "server access key")
	serveCmd.Flags().String("secretKey", "", "server secret key")
	serveCmd.Flags().Bool("debug", false, "toggle debug logging")
	rootCmd.AddCommand(serveCmd)

	initCmd.Flags().Bool("print", false, "print config to stdout")
	initCmd.Flags().StringP("output", "o", config.DataServerDefaultConfig, "specify output directory")
	initCmd.Flags().BoolP("yes", "y", false, "print config to stdout")

	rootCmd.AddCommand(initCmd)

	return rootCmd
}

func Execute() {
	rootCmd := getRootCmd()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
