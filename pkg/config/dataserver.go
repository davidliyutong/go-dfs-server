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

var DataServerDefaultConfig = path.Join(userHomeDir, ".config/go-dfs-server/"+DataServerDefaultConfigName+".yaml")

const DataServerDefaultConfigName = "dataserver"
const DataServerDefaultPort = 27904
const DataServerDefaultVolume = "/data"

const DataServerRole = "dataserver"

type DataServerNetworkOpt struct {
	Port     int64
	Endpoint string
}

type DataServerOpt struct {
	Role    SeverRoleType
	UUID    string
	Network DataServerNetworkOpt
	Volume  string
	Debug   bool
	Log     LogOpt
}

type DataServerDesc struct {
	Opt   DataServerOpt
	Viper *viper.Viper
}

func NewDataServerOpt() DataServerOpt {
	//endpoint := utils.GetEndpointURL()

	return DataServerOpt{
		Role: DataServerRole,
		UUID: "",
		Network: DataServerNetworkOpt{
			Port:     DataServerDefaultPort,
			Endpoint: "0.0.0.0:27904",
		},
		Volume: DataServerDefaultVolume,
		Debug:  false,
		Log: LogOpt{
			Level: "info",
			Path:  "",
		},
	}
}

func NewDataServerDesc() DataServerDesc {
	return DataServerDesc{
		Opt:   NewDataServerOpt(),
		Viper: nil,
	}
}

func InitDataServerCfg(cmd *cobra.Command, args []string) {
	printFlag, _ := cmd.Flags().GetBool("print")
	outputPath, _ := cmd.Flags().GetString("output")
	overwriteFlag, _ := cmd.Flags().GetBool("yes")

	cfg := NewDataServerOpt()
	cfg.UUID = utils.MustGenerateUUID()
	configBuffer, _ := yaml.Marshal(cfg)

	if printFlag {
		fmt.Println(string(configBuffer))
	} else {
		utils.DumpOption(cfg, outputPath, overwriteFlag)
	}
}

func (o *DataServerDesc) Parse(cmd *cobra.Command) error {
	vipCfg := viper.New()
	vipCfg.SetDefault("network.port", DataServerDefaultPort)
	vipCfg.SetDefault("volume", DataServerDefaultVolume)
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
			vipCfg.SetConfigName(DataServerDefaultConfigName)
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
	_ = vipCfg.BindEnv("log.level", "DFSAPP_LOG_LEVEL")
	_ = vipCfg.BindEnv("log.path", "DFSAPP_LOG_PATH")
	vipCfg.AutomaticEnv()

	_ = vipCfg.BindPFlag("uuid", cmd.Flags().Lookup("uuid"))
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

func (o *DataServerDesc) PostParse() {
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

	if o.Opt.UUID == "" {
		log.Infoln("uuid is empty")
		os.Exit(1)
	}
}
