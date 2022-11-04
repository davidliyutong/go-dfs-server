package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-dfs-server/pkg/utils"
	utils2 "go-dfs-server/pkg/utils"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

const ClusterDefaultDomain = "dfs.local"
const ServerDefaultConfigSearchPath0 = "/etc/go-dfs-server"
const ServerDefaultConfigSearchPath1 = "./"

var userHomeDir, _ = os.UserHomeDir()
var ServerDefaultConfigSearchPath2 = path.Join(userHomeDir, ".config/go-dfs-server")
var NameserverDefaultConfig = path.Join(userHomeDir, ".config/go-dfs-server/"+NameserverDefaultConfigName+".yaml")
var DataserverDefaultConfig = path.Join(userHomeDir, ".config/go-dfs-server/"+DataserverDefaultConfigName+".yaml")

const NameserverDefaultConfigName = "nameserver"
const NameserverDefaultPort = 27903
const NameserverDefaultInterface = "0.0.0.0"
const NameserverDefaultVolume = "/data"

const DataserverDefaultConfigName = "dataserver"
const DataserverDefaultPort = 27904
const DataserverDefaultVolume = "/data"

type SeverRoleType string

const NameserverRole = "nameserver"
const DataserverRole = "dataserver"

type NameserverNetworkOpt struct {
	Port      int
	Interface string
}

type DataserverNetworkOpt struct {
	Port     int
	Remote   string
	Endpoint string
}

type AuthOpt struct {
	Domain    string
	AccessKey string
	SecretKey string
}

type LogOpt struct {
	Level string
	Path  string
}

type NameserverOpt struct {
	Role    SeverRoleType
	Network NameserverNetworkOpt
	Volume  string
	Auth    AuthOpt
	Debug   bool
	Log     LogOpt
}

type DataserverOpt struct {
	Role    SeverRoleType
	Network DataserverNetworkOpt
	Volume  string
	Auth    AuthOpt
	Debug   bool
	Log     LogOpt
}

func GetNameserverOpt() NameserverOpt {
	accessKey, secretKey := utils.MustGenerateAuthKeys()

	return NameserverOpt{
		Role: NameserverRole,
		Network: NameserverNetworkOpt{
			Port:      NameserverDefaultPort,
			Interface: NameserverDefaultInterface,
		},
		Volume: NameserverDefaultVolume,
		Auth: AuthOpt{
			Domain:    ClusterDefaultDomain,
			AccessKey: accessKey,
			SecretKey: secretKey,
		},
		Debug: false,
		Log: LogOpt{
			Level: "info",
			Path:  "",
		},
	}

}

func GetDataserverOpt() DataserverOpt {
	endpoint := utils.GetEndpointURL()

	return DataserverOpt{
		Role: DataserverRole,
		Network: DataserverNetworkOpt{
			Port:     DataserverDefaultPort,
			Endpoint: endpoint,
		},
		Volume: DataserverDefaultVolume,
		Auth: AuthOpt{
			Domain:    ClusterDefaultDomain,
			AccessKey: "",
			SecretKey: "",
		},
		Debug: false,
		Log: LogOpt{
			Level: "info",
			Path:  "",
		},
	}

}

func NameserverInit(cmd *cobra.Command, args []string) {
	printFlag, _ := cmd.Flags().GetBool("print")
	outputPath, _ := cmd.Flags().GetString("output")
	overwriteFlag, _ := cmd.Flags().GetBool("yes")

	cfg := GetNameserverOpt()
	configBuffer, _ := yaml.Marshal(cfg)

	if printFlag {
		fmt.Println(string(configBuffer))
	} else {
		utils2.DumpOption(cfg, outputPath, overwriteFlag)
	}
}

func DataserverInit(cmd *cobra.Command, args []string) {
	printFlag, _ := cmd.Flags().GetBool("print")
	outputPath, _ := cmd.Flags().GetString("output")
	overwriteFlag, _ := cmd.Flags().GetBool("yes")

	cfg := GetDataserverOpt()
	configBuffer, _ := yaml.Marshal(cfg)

	if printFlag {
		fmt.Println(string(configBuffer))
	} else {
		utils2.DumpOption(cfg, outputPath, overwriteFlag)
	}
}

func (o *NameserverOpt) Parse(cmd *cobra.Command) (*viper.Viper, error) {
	vipCfg := viper.New()
	vipCfg.SetDefault("network.port", NameserverDefaultPort)
	vipCfg.SetDefault("network.interface", NameserverDefaultInterface)
	vipCfg.SetDefault("volume", NameserverDefaultVolume)
	vipCfg.SetDefault("auth.domain", ClusterDefaultDomain)
	vipCfg.SetDefault("debug", false)
	vipCfg.SetDefault("log.debug", "info")
	vipCfg.SetDefault("log.path", "")

	if configFileCmd, err := cmd.Flags().GetString("config"); err == nil && configFileCmd != "" {
		vipCfg.SetConfigFile(configFileCmd)
	} else {
		configFileEnv := os.Getenv("DFSAPP_CONFIG")
		if configFileEnv != "" {
			vipCfg.SetConfigFile(configFileEnv)
		} else {
			vipCfg.SetConfigName(NameserverDefaultConfigName)
			vipCfg.SetConfigType("yaml")
			vipCfg.AddConfigPath(ServerDefaultConfigSearchPath0)
			vipCfg.AddConfigPath(ServerDefaultConfigSearchPath1)
			vipCfg.AddConfigPath(ServerDefaultConfigSearchPath2)
		}
	}
	vipCfg.WatchConfig()

	vipCfg.SetEnvPrefix("DFSAPP")
	_ = vipCfg.BindEnv("network.port", "DFSAPP_PORT")
	_ = vipCfg.BindEnv("network.interface", "DFSAPP_INTERFACE")
	_ = vipCfg.BindEnv("auth.domain", "DFSAPP_DOMAIN")
	_ = vipCfg.BindEnv("auth.accesskey", "DFSAPP_ACCESSKEY")
	_ = vipCfg.BindEnv("auth.secretkey", "DFSAPP_SECRETKEY")
	_ = vipCfg.BindEnv("debug", "DFSAPP_DEBUG")
	_ = vipCfg.BindEnv("log.level", "DFSAPP_LOG_LEVEL")
	_ = vipCfg.BindEnv("log.path", "DFSAPP_LOG_PATH")
	vipCfg.AutomaticEnv()

	_ = vipCfg.BindPFlag("network.port", cmd.Flags().Lookup("port"))
	_ = vipCfg.BindPFlag("network.interface", cmd.Flags().Lookup("interface"))
	_ = vipCfg.BindPFlag("volume", cmd.Flags().Lookup("volume"))
	_ = vipCfg.BindPFlag("auth.domain", cmd.Flags().Lookup("domain"))
	_ = vipCfg.BindPFlag("auth.accesskey", cmd.Flags().Lookup("accessKey"))
	_ = vipCfg.BindPFlag("auth.secretkey", cmd.Flags().Lookup("secretKey"))
	_ = vipCfg.BindPFlag("debug", cmd.Flags().Lookup("debug"))

	// If a config file is found, read it in.
	if err := vipCfg.ReadInConfig(); err == nil {
		log.Infoln("using config file:", vipCfg.ConfigFileUsed())
	} else {
		log.Warnln(err)
		return vipCfg, err
	}

	if err := vipCfg.Unmarshal(o); err != nil {
		log.Panicln("failed to unmarshal config")
	}

	return vipCfg, nil
}

func (o *DataserverOpt) Parse(cmd *cobra.Command) (*viper.Viper, error) {
	vipCfg := viper.New()
	vipCfg.SetDefault("network.port", DataserverDefaultPort)
	vipCfg.SetDefault("volume", DataserverDefaultVolume)
	vipCfg.SetDefault("auth.domain", ClusterDefaultDomain)
	vipCfg.SetDefault("debug", false)
	vipCfg.SetDefault("log.debug", "info")
	vipCfg.SetDefault("log.path", "")

	if configFileCmd, err := cmd.Flags().GetString("config"); err == nil && configFileCmd != "" {
		vipCfg.SetConfigFile(configFileCmd)
	} else {
		configFileEnv := os.Getenv("DFSAPP_CONFIG")
		if configFileEnv != "" {
			vipCfg.SetConfigFile(configFileEnv)
		} else {
			vipCfg.SetConfigName(DataserverDefaultConfigName)
			vipCfg.SetConfigType("yaml")
			vipCfg.AddConfigPath(ServerDefaultConfigSearchPath0)
			vipCfg.AddConfigPath(ServerDefaultConfigSearchPath1)
			vipCfg.AddConfigPath(ServerDefaultConfigSearchPath2)
		}
	}
	vipCfg.WatchConfig()

	vipCfg.SetEnvPrefix("DFSAPP")
	_ = vipCfg.BindEnv("network.port", "DFSAPP_PORT")
	_ = vipCfg.BindEnv("network.remote", "DFSAPP_REMOTE")
	_ = vipCfg.BindEnv("network.endpoint", "DFSAPP_ENDPOINT")
	_ = vipCfg.BindEnv("auth.domain", "DFSAPP_DOMAIN")
	_ = vipCfg.BindEnv("auth.accesskey", "DFSAPP_ACCESSKEY")
	_ = vipCfg.BindEnv("auth.secretkey", "DFSAPP_SECRETKEY")
	_ = vipCfg.BindEnv("debug", "DFSAPP_DEBUG")
	_ = vipCfg.BindEnv("log.level", "DFSAPP_LOG_LEVEL")
	_ = vipCfg.BindEnv("log.path", "DFSAPP_LOG_PATH")
	vipCfg.AutomaticEnv()

	_ = vipCfg.BindPFlag("network.port", cmd.Flags().Lookup("port"))
	_ = vipCfg.BindPFlag("network.remote", cmd.Flags().Lookup("remote"))
	_ = vipCfg.BindPFlag("network.endpoint", cmd.Flags().Lookup("endpoint"))
	_ = vipCfg.BindPFlag("volume", cmd.Flags().Lookup("volume"))
	_ = vipCfg.BindPFlag("auth.domain", cmd.Flags().Lookup("domain"))
	_ = vipCfg.BindPFlag("auth.accesskey", cmd.Flags().Lookup("accessKey"))
	_ = vipCfg.BindPFlag("auth.secretkey", cmd.Flags().Lookup("secretKey"))
	_ = vipCfg.BindPFlag("debug", cmd.Flags().Lookup("debug"))

	if err := vipCfg.ReadInConfig(); err == nil {
		log.Infoln("using config file:", vipCfg.ConfigFileUsed())
	} else {
		log.Warnln(err)
		return vipCfg, err
	}

	if err := vipCfg.Unmarshal(o); err != nil {
		log.Panicln("failed to unmarshal config")
	}

	return vipCfg, nil
}

func (o *NameserverOpt) PostParse() {
	if o.Debug || o.Log.Level == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		lvl, err := log.ParseLevel(o.Log.Level)
		if err != nil {
			log.Errorf("error parsing loglevel: %s, using INFO", err)
			lvl = log.InfoLevel
		}
		log.SetLevel(lvl)
	}
}

func (o *DataserverOpt) PostParse() {
	if o.Debug || o.Log.Level == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		lvl, err := log.ParseLevel(o.Log.Level)
		if err != nil {
			log.Errorf("error parsing loglevel: %s, using INFO", err)
			lvl = log.InfoLevel
		}
		log.SetLevel(lvl)
	}
}
