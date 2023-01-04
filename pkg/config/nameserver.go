package config

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-dfs-server/pkg/utils"
	"gopkg.in/yaml.v2"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
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
	DataServers []RegisteredDataServer
}

type NameServerDesc struct {
	Opt   NameServerOpt
	Viper *viper.Viper
	UUID  string
}

func (o *NameServerOpt) AuthIsEnabled() bool {
	return o.Auth.AccessKey != "" && o.Auth.SecretKey != ""
}

type RegisteredDataServer struct {
	UUID    string
	Address string
	Port    int64
	UseTLS  bool
}

func getRegisteredDataServerFromEnv() []RegisteredDataServer {
	defaultValue := []RegisteredDataServer{
		{UUID: "", Address: "example.com", Port: 27904, UseTLS: false},
		{UUID: "", Address: "example.com", Port: 27905, UseTLS: false},
		{UUID: "", Address: "127.0.0.1", Port: 27906, UseTLS: true},
	}
	if os.Getenv("DFSAPP_DATA_SERVERS") == "" {
		return defaultValue
	} else {
		parsedValue := make([]RegisteredDataServer, 0)
		servers := strings.Split(os.Getenv("DFSAPP_DATA_SERVERS"), ",")
		for _, server := range servers {
			addr, port, err := net.SplitHostPort(server)
			if err != nil {
				continue
			} else {
				portDigit, _ := strconv.Atoi(port)
				parsedValue = append(parsedValue, RegisteredDataServer{UUID: "", Address: addr, Port: int64(portDigit), UseTLS: false})
			}
		}
		if len(parsedValue) <= 0 {
			return defaultValue
		} else {
			return parsedValue
		}
	}
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
		DataServers: getRegisteredDataServerFromEnv(),
	}
}

func NewNameServerDesc() NameServerDesc {
	return NameServerDesc{
		Opt:   NewNameServerOpt(),
		Viper: nil,
		UUID:  "",
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

func OutputServerCredential(cmd *cobra.Command, args []string) {
	/** 创建NameServerOption **/
	desc := NewNameServerDesc()
	if err := desc.Parse(cmd); err != nil {
		log.Fatalln("failed to parse configuration", err)
		os.Exit(1)
	} else {
		fmt.Printf("export DFSAPP_ACCESSKEY=%v;DFSAPP_SECRETKEY=%v;\n", desc.Opt.Auth.AccessKey, desc.Opt.Auth.SecretKey)
		return
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
		log.Debugln("using config file:", vipCfg.ConfigFileUsed())
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
	o.UUID = utils.MustGenerateUUID()
}

func (o *NameServerDesc) SaveConfig() error {
	f, err := os.OpenFile(o.Viper.ConfigFileUsed(), os.O_CREATE|os.O_RDWR, 0644)
	defer func() { _ = f.Close() }()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	s, _ := yaml.Marshal(o.Opt)
	_, err = w.Write(s)
	if err != nil {
		return err
	}
	_ = w.Flush()
	return nil
}
