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
	"time"
)

const ClusterDefaultDomain = "dfs.local"
const ServerDefaultConfigSearchPath0 = "/etc/go-dfs-server"
const ServerDefaultConfigSearchPath1 = "./"

var userHomeDir, _ = os.UserHomeDir()
var ServerDefaultConfigSearchPath2 = path.Join(userHomeDir, ".config/go-dfs-server")
var ClientDefaultConfigSearchPath2 = path.Join(userHomeDir, ".config/go-dfs-server")
var NameserverDefaultConfig = path.Join(userHomeDir, ".config/go-dfs-server/"+NameserverDefaultConfigName+".yaml")
var DataserverDefaultConfig = path.Join(userHomeDir, ".config/go-dfs-server/"+DataserverDefaultConfigName+".yaml")
var ClientDefaultConfig = path.Join(userHomeDir, ".config/go-dfs-server/"+ClientDefaultConfigName+".yaml")

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

const ClientDefaultConfigName = "client"
const ClientDefaultConfigSearchPath0 = "/etc/go-dfs-server"
const ClientDefaultConfigSearchPath1 = "./"

type NameserverNetworkOpt struct {
	Port      int
	Interface string
}

type DataserverNetworkOpt struct {
	Port     int
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

type NameserverDesc struct {
	Opt   NameserverOpt
	Viper *viper.Viper
}

func (o *NameserverOpt) AuthIsEnabled() bool {
	return o.Auth.AccessKey != "" && o.Auth.SecretKey != ""
}

type DataserverOpt struct {
	Role    SeverRoleType
	Network DataserverNetworkOpt
	Volume  string
	Debug   bool
	Log     LogOpt
}

type DataserverDesc struct {
	Opt   DataserverOpt
	Viper *viper.Viper
}

type ClientOpt struct {
	Token   string
	Expire  time.Time
	Address string
	Port    int16
	UseTLS  bool
}

type ClientAuthOpt struct {
	AccessKey string
	SecretKey string
	Token     string
	Expire    time.Time
}

func (o *ClientAuthOpt) AuthIsEnabled() bool {
	return o.Token != ""
}

func (o *ClientOpt) GetHTTPUrl() string {
	if o.UseTLS {
		s := fmt.Sprintf("%s://%s:%d", "https", o.Address, o.Port)
		return s
	} else {
		s := fmt.Sprintf("%s://%s:%d", "http", o.Address, o.Port)
		return s
	}
}

func NewNameserverOpt() NameserverOpt {
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

func NewNameserverDesc() NameserverDesc {
	return NameserverDesc{
		Opt:   NewNameserverOpt(),
		Viper: nil,
	}
}

func NewDataserverOpt() DataserverOpt {
	endpoint := utils.GetEndpointURL()

	return DataserverOpt{
		Role: DataserverRole,
		Network: DataserverNetworkOpt{
			Port:     DataserverDefaultPort,
			Endpoint: endpoint,
		},
		Volume: DataserverDefaultVolume,
		Debug:  false,
		Log: LogOpt{
			Level: "info",
			Path:  "",
		},
	}

}

func NewDataserverDesc() DataserverDesc {
	return DataserverDesc{
		Opt:   NewDataserverOpt(),
		Viper: nil,
	}
}

func NewClientOpt() ClientOpt {
	return ClientOpt{}
}

func NewClientAuthOpt() ClientAuthOpt {
	return ClientAuthOpt{}
}

func InitNameserverCfg(cmd *cobra.Command, args []string) {
	printFlag, _ := cmd.Flags().GetBool("print")
	outputPath, _ := cmd.Flags().GetString("output")
	overwriteFlag, _ := cmd.Flags().GetBool("yes")

	cfg := NewNameserverOpt()
	configBuffer, _ := yaml.Marshal(cfg)

	if printFlag {
		fmt.Println(string(configBuffer))
	} else {
		utils2.DumpOption(cfg, outputPath, overwriteFlag)
	}
}

func InitDataserverCfg(cmd *cobra.Command, args []string) {
	printFlag, _ := cmd.Flags().GetBool("print")
	outputPath, _ := cmd.Flags().GetString("output")
	overwriteFlag, _ := cmd.Flags().GetBool("yes")

	cfg := NewDataserverOpt()
	configBuffer, _ := yaml.Marshal(cfg)

	if printFlag {
		fmt.Println(string(configBuffer))
	} else {
		utils2.DumpOption(cfg, outputPath, overwriteFlag)
	}
}

func (o *NameserverDesc) Parse(cmd *cobra.Command) error {
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
		return nil
	}

	if err := vipCfg.Unmarshal(&o.Opt); err != nil {
		log.Fatalln("failed to unmarshal config")
		os.Exit(1)
	}
	o.Viper = vipCfg
	return nil
}

func (o *DataserverDesc) Parse(cmd *cobra.Command) error {
	vipCfg := viper.New()
	vipCfg.SetDefault("network.port", DataserverDefaultPort)
	vipCfg.SetDefault("volume", DataserverDefaultVolume)
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
	_ = vipCfg.BindEnv("network.endpoint", "DFSAPP_ENDPOINT")
	_ = vipCfg.BindEnv("debug", "DFSAPP_DEBUG")
	_ = vipCfg.BindEnv("log.level", "DFSAPP_LOG_LEVEL")
	_ = vipCfg.BindEnv("log.path", "DFSAPP_LOG_PATH")
	vipCfg.AutomaticEnv()

	_ = vipCfg.BindPFlag("network.port", cmd.Flags().Lookup("port"))
	_ = vipCfg.BindPFlag("network.endpoint", cmd.Flags().Lookup("endpoint"))
	_ = vipCfg.BindPFlag("volume", cmd.Flags().Lookup("volume"))
	_ = vipCfg.BindPFlag("debug", cmd.Flags().Lookup("debug"))

	if err := vipCfg.ReadInConfig(); err == nil {
		log.Infoln("using config file:", vipCfg.ConfigFileUsed())
	} else {
		log.Warnln(err)
		return err
	}

	if err := vipCfg.Unmarshal(&o.Opt); err != nil {
		log.Fatalln("failed to unmarshal config")
		os.Exit(1)
	}
	o.Viper = vipCfg

	return nil
}

func (o *ClientOpt) Parse(cmd *cobra.Command) (*viper.Viper, error) {
	vipCfg := viper.New()
	vipCfg.SetDefault("_config", ClientDefaultConfig)

	if configFileCmd, err := cmd.Flags().GetString("config"); err == nil && configFileCmd != "" {
		vipCfg.SetConfigFile(configFileCmd)
		vipCfg.Set("_config", configFileCmd)
	} else {
		configFileEnv := os.Getenv("DFSAPP_CONFIG")
		if configFileEnv != "" {
			vipCfg.SetConfigFile(configFileEnv)
			vipCfg.Set("_config", configFileEnv)
		} else {
			vipCfg.SetConfigName(ClientDefaultConfigName)
			vipCfg.SetConfigType("yaml")
			vipCfg.AddConfigPath(ClientDefaultConfigSearchPath0)
			vipCfg.AddConfigPath(ClientDefaultConfigSearchPath1)
			vipCfg.AddConfigPath(ClientDefaultConfigSearchPath2)
		}
	}
	if err := vipCfg.ReadInConfig(); err == nil {
		log.Debugln("using config file:", vipCfg.ConfigFileUsed())
		vipCfg.Set("_config", vipCfg.ConfigFileUsed())

	} else {
		log.Info(err)
		return vipCfg, err
	}

	if err := vipCfg.Unmarshal(o); err != nil {
		log.Errorln("failed to unmarshal config", vipCfg.ConfigFileUsed())
		os.Exit(1)
	}

	return vipCfg, nil
}

func (o *NameserverDesc) PostParse() {
	if o.Opt.Debug || o.Opt.Log.Level == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		lvl, err := log.ParseLevel(o.Opt.Log.Level)
		if err != nil {
			log.Errorf("error parsing loglevel: %s, using INFO", err)
			lvl = log.InfoLevel
		}
		log.SetLevel(lvl)
	}
}

func (o *DataserverDesc) PostParse() {
	if o.Opt.Debug || o.Opt.Log.Level == "debug" {
		log.SetLevel(log.DebugLevel)
	} else {
		lvl, err := log.ParseLevel(o.Opt.Log.Level)
		if err != nil {
			log.Errorf("error parsing loglevel: %s, using INFO", err)
			lvl = log.InfoLevel
		}
		log.SetLevel(lvl)
	}
}
