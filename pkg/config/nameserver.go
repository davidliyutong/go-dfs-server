package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-dfs-server/pkg/utils"
	"gopkg.in/yaml.v2"
	"os"
	"path"
)

var NameServerDefaultConfig = path.Join(userHomeDir, ".config/go-dfs-server/"+NameServerDefaultConfigName+".yaml")

const NameServerDefaultConfigName = "nameserver"
const NameServerDefaultPort = 27903
const NameServerDefaultInterface = "0.0.0.0"
const NameServerDefaultVolume = "/data"

const NameServerRole = "nameserver"

type NameServerNetworkOpt struct {
	Port      int
	Interface string
}

type NameServerOpt struct {
	Role        SeverRoleType
	UUID        string
	Network     NameServerNetworkOpt
	Volume      string
	Auth        AuthOpt
	Debug       bool
	Log         LogOpt
	DataServers []string
}

type NameServerDesc struct {
	Opt   NameServerOpt
	Viper *viper.Viper
}

func (o *NameServerOpt) AuthIsEnabled() bool {
	return o.Auth.AccessKey != "" && o.Auth.SecretKey != ""
}

func NewNameServerOpt() NameServerOpt {
	accessKey, secretKey := utils.MustGenerateAuthKeys()

	return NameServerOpt{
		Role: NameServerRole,
		UUID: "",
		Network: NameServerNetworkOpt{
			Port:      NameServerDefaultPort,
			Interface: NameServerDefaultInterface,
		},
		Volume: NameServerDefaultVolume,
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

func NewNameServerDesc() NameServerDesc {
	return NameServerDesc{
		Opt:   NewNameServerOpt(),
		Viper: nil,
	}
}

func InitNameServerCfg(cmd *cobra.Command, args []string) {
	printFlag, _ := cmd.Flags().GetBool("print")
	outputPath, _ := cmd.Flags().GetString("output")
	overwriteFlag, _ := cmd.Flags().GetBool("yes")

	cfg := NewNameServerOpt()
	cfg.UUID = utils.MustGenerateUUID()
	configBuffer, _ := yaml.Marshal(cfg)

	if printFlag {
		fmt.Println(string(configBuffer))
	} else {
		utils.DumpOption(cfg, outputPath, overwriteFlag)
	}
}

func (o *NameServerDesc) Parse(cmd *cobra.Command) error {
	vipCfg := viper.New()
	vipCfg.SetDefault("network.port", NameServerDefaultPort)
	vipCfg.SetDefault("network.interface", NameServerDefaultInterface)
	vipCfg.SetDefault("volume", NameServerDefaultVolume)
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
			vipCfg.SetConfigName(NameServerDefaultConfigName)
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
	_ = vipCfg.BindEnv("log.level", "DFSAPP_LOG_LEVEL")
	_ = vipCfg.BindEnv("log.path", "DFSAPP_LOG_PATH")
	vipCfg.AutomaticEnv()

	_ = vipCfg.BindPFlag("uuid", cmd.Flags().Lookup("uuid"))
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

func (o *NameServerDesc) PostParse() {
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
