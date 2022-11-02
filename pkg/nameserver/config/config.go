package config

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-dfs-server/pkg/nameserver/utils"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
)

const NameserverDefaultConfigName = "nameserver"
const NameserverDefaultPort = 27903
const NameserverDefaultInterface = "0.0.0.0"
const NameserverDefaultVolume = "/data"
const NameserverDefaultConfigSearchPath0 = "/etc/go-dfs-server"
const NameserverDefaultConfigSearchPath1 = "./"

var userHomeDir, _ = os.UserHomeDir()
var NameserverDefaultConfig = path.Join(userHomeDir, ".config/go-dfs-server/nameserver.yaml")
var NameserverDefaultConfigSearchPath2 = path.Join(userHomeDir, ".config/go-dfs-server")

type NetworkCfg struct {
	Port      int
	Interface string
}

type AuthCfg struct {
	AccessKey string
	SecretKey string
}

type NameserverCfg struct {
	Network NetworkCfg
	Volume  string
	Auth    AuthCfg
}

func Init(cmd *cobra.Command, args []string) {
	printFlag, _ := cmd.Flags().GetBool("print")
	outputPath, _ := cmd.Flags().GetString("output")

	accessKey, secretKey := utils.GenerateAuthenticationKeys()
	cfg := NameserverCfg{
		Network: NetworkCfg{
			Port:      NameserverDefaultPort,
			Interface: NameserverDefaultInterface,
		},
		Volume: NameserverDefaultVolume,
		Auth: AuthCfg{
			AccessKey: accessKey,
			SecretKey: secretKey,
		},
	}
	configBuffer, _ := yaml.Marshal(cfg)

	if printFlag {
		fmt.Println(string(configBuffer))
	} else {
		parentPath := path.Dir(outputPath)
		if _, err := os.Stat(parentPath); os.IsNotExist(err) {
			err = os.MkdirAll(parentPath, 0644)
			if err != nil {
				log.Panicln("cannot create directory", parentPath)
			}
		}

		if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
			ret := utils.AskForConfirmationDefaultYes("configuration " + outputPath + " already exist, overwrite?")
			if !ret {
				log.Println("abort")
				return
			}
		}

		log.Println("writing default configuration to", outputPath)
		f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, 0644)
		defer func() { _ = f.Close() }()
		if err != nil {
			panic("cannot open " + outputPath + ", check permissions")
		}

		w := bufio.NewWriter(f)
		_, err = w.Write(configBuffer)
		if err != nil {
			log.Panicln("cannot write configuration", err)
		}
		_ = w.Flush()
		_ = f.Close()
	}
}

func Parse(cmd *cobra.Command) (*viper.Viper, NameserverCfg, error) {
	vipCfg := viper.New()
	vipCfg.SetDefault("network.port", NameserverDefaultPort)
	vipCfg.SetDefault("network.interface", NameserverDefaultInterface)
	vipCfg.SetDefault("volume", NameserverDefaultVolume)

	if configFileCmd, err := cmd.Flags().GetString("config"); err == nil && configFileCmd != "" {
		vipCfg.SetConfigFile(configFileCmd)
	} else {
		configFileEnv := os.Getenv("DFSAPP_CONFIG")
		if configFileEnv != "" {
			vipCfg.SetConfigFile(configFileEnv)
		} else {
			vipCfg.SetConfigName(NameserverDefaultConfigName)
			vipCfg.SetConfigType("yaml")
			vipCfg.AddConfigPath(NameserverDefaultConfigSearchPath0)
			vipCfg.AddConfigPath(NameserverDefaultConfigSearchPath1)
			vipCfg.AddConfigPath(NameserverDefaultConfigSearchPath2)
		}
	}
	vipCfg.WatchConfig()

	_ = viper.BindPFlag("network.port", cmd.Flags().Lookup("port"))
	_ = viper.BindPFlag("network.interface", cmd.Flags().Lookup("interface"))
	_ = viper.BindPFlag("volume", cmd.Flags().Lookup("volume"))
	_ = viper.BindPFlag("auth.accessKey", cmd.Flags().Lookup("accessKey"))
	_ = viper.BindPFlag("auth.secretKey", cmd.Flags().Lookup("secretKey"))

	vipCfg.SetEnvPrefix("DFSAPP")
	vipCfg.AutomaticEnv()

	// If a config file is found, read it in.
	serverCfg := NameserverCfg{}
	if err := vipCfg.ReadInConfig(); err == nil {
		fmt.Println("using config file:", viper.ConfigFileUsed())
	} else {
		return vipCfg, serverCfg, err
	}

	if err := vipCfg.Unmarshal(&serverCfg); err != nil {
		log.Panicln("failed to unmarshal config")
	}

	return vipCfg, serverCfg, nil
}
