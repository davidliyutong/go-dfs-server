package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	client "go-dfs-server/pkg/client"
	"os"
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

If --config /path/to/config is specified, client will first try to exam token stored in /path/to/config and refresh the 
token.

The token obtained from server will be stored in /path/to/config, which is ~/.config/go-dfs-server/client.yaml by default
`,
	Example: `  go-dfs-client login dfs://127.0.0.1
  go-dfs-client login dfs://127.0.0.1 --accessKey=12345678 --secretKey=xxxxxxxx
  go-dfs-client login dfs://127.0.0.1:27903 --accessKey=12345678 --secretKey=xxxxxxxx`,
	Args: cobra.MaximumNArgs(1),
	Run:  client.Login,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: `logout clears server credentials.`,
	Args:  cobra.ExactArgs(0),
	Run:   client.Logout,
}

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: `ls lists a remote directory or file`,
	Long: `list lists a remote directory or file, the path or file must exists. If path is a file, it will display the 
file info in json format. If path is a directory, it will display the all files in directory in a list.
`,
	Args:    cobra.ExactArgs(1),
	Example: `  go-dfs-client ls <dir|file>`,
	Run:     client.Ls,
}

var catCmd = &cobra.Command{
	Use:   "cat",
	Short: `cat a remote file`,
	Long: `cat a remote file to stdout, the file must exists.
`,
	Args:    cobra.ExactArgs(1),
	Example: `  go-dfs-client cat <file>`,
	Run:     client.Cat,
}

var pipeCmd = &cobra.Command{
	Use:   "pipe",
	Short: `send pipe to remote file`,
	Long: `send pipe to remote file, the file must exists. For example cat test.txt | go-dfs-client pipe test.txt
`,
	Args:    cobra.ExactArgs(1),
	Example: `  | go-dfs-client pipe <file>`,
	Run:     client.Pipe,
}

var rmCmd = &cobra.Command{
	Use:     "rm",
	Short:   `remove remote file or directory, set -r for recursive`,
	Args:    cobra.ExactArgs(1),
	Example: `  go-dfs-client rm [-r] <path>`,
	Run:     client.Rm,
}

var mkdirCmd = &cobra.Command{
	Use:     "mkdir",
	Short:   `create a remote directory, parent directory must exist`,
	Args:    cobra.ExactArgs(1),
	Example: `  go-dfs-client mkdir <path>`,
	Run:     client.Mkdir,
}

var touchCmd = &cobra.Command{
	Use:     "touch",
	Short:   `create an empty remote file, parent directory must exist`,
	Args:    cobra.ExactArgs(1),
	Example: `  go-dfs-client touch <path>`,
	Run:     client.Touch,
}

var getCmd = &cobra.Command{
	Use:     "get",
	Short:   `download file|directory from remote`,
	Args:    cobra.ExactArgs(2),
	Example: `  go-dfs-client get [-r] <remote_path> <local_path>`,
	Run:     client.Get,
}

var putCmd = &cobra.Command{
	Use:     "put",
	Short:   `upload file|directory to remote`,
	Args:    cobra.ExactArgs(2),
	Example: `  go-dfs-client put [-r] <local_path> <remote_path>`,
	Run:     client.Put,
}

func getRootCmd() *cobra.Command {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")

	loginCmd.Flags().String("accessKey", "", "server access key")
	loginCmd.Flags().String("secretKey", "", "server secret key")
	loginCmd.Flags().String("config", "", "default configuration path")

	rootCmd.AddCommand(loginCmd)

	rootCmd.AddCommand(logoutCmd)

	rootCmd.AddCommand(lsCmd)

	rootCmd.AddCommand(catCmd)

	rootCmd.AddCommand(pipeCmd)

	rmCmd.Flags().BoolP("recursive", "r", false, "recursive remove")
	rmCmd.Flags().BoolP("force", "f", false, "force remove")

	rootCmd.AddCommand(rmCmd)

	rootCmd.AddCommand(mkdirCmd)

	rootCmd.AddCommand(touchCmd)

	getCmd.Flags().BoolP("recursive", "r", false, "recursive put")
	getCmd.Flags().BoolP("force", "f", false, "recursive put")
	rootCmd.AddCommand(getCmd)

	putCmd.Flags().BoolP("recursive", "r", false, "recursive put")
	putCmd.Flags().BoolP("force", "f", false, "force put")

	rootCmd.AddCommand(putCmd)

	return rootCmd
}

func Execute() {
	rootCmd := getRootCmd()
	verboseFlag, _ := rootCmd.PersistentFlags().GetBool("verbose")
	verboseFlag = verboseFlag || os.Getenv("DFSAPP_VERBOSE") != ""
	if verboseFlag {
		log.SetLevel(log.DebugLevel)
	}
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
