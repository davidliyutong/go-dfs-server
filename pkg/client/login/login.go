package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go-dfs-server/pkg/auth"
	"go-dfs-server/pkg/config"
	"go-dfs-server/pkg/nameserver/server"
	"go-dfs-server/pkg/utils"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

type ClientOpt struct {
	AccessKey    string
	SecretKey    string
	Authenticate bool
	Token        string
	Expire       time.Time
	Address      string
	Port         int16
	UseTLS       bool
}

func GetClientOpt() ClientOpt {
	return ClientOpt{}
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

func (o *ClientOpt) BindURL(url string) error {
	protoReg := regexp.MustCompile("^(dfs|dfss)://")
	ipReg := regexp.MustCompile(`(dfs|dfss)://([0-9.]*)[:]?`)
	portReg := regexp.MustCompile(`(dfs|dfss)://::([0-9]*)|(dfs|dfss)://[^:\\n]*:([0-9]*)`)

	protoMatch := protoReg.FindAllStringSubmatch(url, -1)
	ipMatch := ipReg.FindAllStringSubmatch(url, -1)
	portMatch := portReg.FindAllStringSubmatch(url, -1)

	if len(protoMatch) < 1 {
		return errors.New("wrong url format: " + url)
	} else {
		if protoMatch[0][1] == "dfs" {
			o.UseTLS = false
		} else {
			o.UseTLS = true
		}
	}

	var ipString string
	if len(ipMatch) < 1 {
		ipString = "127.0.0.1"
	} else {
		if ipMatch[0][2] == "" {
			ipString = "127.0.0.1"
		} else {
			ipString = ipMatch[0][2]
		}
	}

	var portInt int
	if len(portMatch) < 1 {
		portInt = config.NameserverDefaultPort
	} else {
		portInt = config.NameserverDefaultPort
		for _, portStr := range portMatch[0][2:] {
			pportInt, err := strconv.Atoi(portStr)
			if err == nil {
				portInt = pportInt
				break
			}
		}
	}

	log.Debugf("remote is %s, port is %d", ipString, portInt)
	o.Address = ipString
	o.Port = int16(portInt)
	return nil
}

func (o *ClientOpt) MustBindURL(url string) {
	err := o.BindURL(url)
	if err != nil {
		log.Errorln(err)
		os.Exit(2)
	}
}

type DFSClusterInfo struct {
	AccessKey string `json:"accessKey"`
	Message   string `json:"message"`
}

func (o *ClientOpt) Info() (DFSClusterInfo, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", o.GetHTTPUrl()+server.NameserverAPIPrefix+server.NameserverInfoPath, nil)
	if err != nil {
		return DFSClusterInfo{}, err
	}

	if o.Authenticate {
		request.Header.Add("Authorization", "Bearer "+o.Token)
	}
	request.Header.Add("Content-Type", "application/json")

	respHandle, err := client.Do(request)
	if err != nil {
		return DFSClusterInfo{}, err
	}

	defer func() {
		err := respHandle.Body.Close()
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
	}()

	content, err := io.ReadAll(respHandle.Body)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	log.Debugln("dfs response:", string(content))

	response := DFSClusterInfo{}
	err = json.Unmarshal(content, &response)
	if err != nil {
		return DFSClusterInfo{}, err
	}

	return response, nil
}

func (o *ClientOpt) BindAuthentication(cmd *cobra.Command) error {
	accessKey, _ := cmd.Flags().GetString("accessKey")
	secretKey, _ := cmd.Flags().GetString("secretKey")
	if accessKey == "" {
		fmt.Printf("Input accessKey:")
		_, err := fmt.Scanf("%s", &accessKey)
		if err != nil && err.Error() != "unexpected newline" {
			return err
		}
	}
	if secretKey == "" {
		fmt.Printf("Input secretKey:")
		_, err := fmt.Scanf("%s", &secretKey)
		if err != nil && err.Error() != "unexpected newline" {
			return err
		}
	}
	log.Debugf("accesskey: %s, secretKey: %s", accessKey, secretKey)

	o.Authenticate = !(accessKey == "")
	if o.Authenticate != false {

		o.AccessKey = accessKey
		o.SecretKey = secretKey

		credentials := auth.LoginCredential{
			AccessKey: o.AccessKey,
			SecretKey: o.SecretKey,
		}
		credentialsBytes, _ := json.Marshal(credentials)

		respHandle, err := http.Post(o.GetHTTPUrl()+server.NameserverLoginPath, "application/json", bytes.NewBuffer(credentialsBytes))
		if err != nil {
			return err
		}

		defer func() {
			err := respHandle.Body.Close()
			if err != nil {
				log.Errorln(err)
				os.Exit(1)
			}
		}()

		content, err := io.ReadAll(respHandle.Body)
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
		log.Debugln("dfs response:", string(content))

		response := auth.JWTResponse{}
		err = json.Unmarshal(content, &response)
		if err != nil {
			return err
		}

		if response.Code == 200 {
			o.Token = response.Token
			o.Expire = response.Expire
		} else {
			return errors.New("login failed, access denied")
		}
	} else {
		o.Token = ""
		o.Expire = time.UnixMicro(0)
		o.AccessKey = ""
		o.SecretKey = ""
	}

	_, err := o.Info()
	if err != nil {
		return err
	}

	return nil
}

func (o *ClientOpt) MustBindAuthentication(cmd *cobra.Command) {
	err := o.BindAuthentication(cmd)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

}

func (o *ClientOpt) Refresh() error {
	client := &http.Client{}

	request, err := http.NewRequest("POST", o.GetHTTPUrl()+server.NameserverTokenRefreshPath, nil)
	if err != nil {
		return err
	}
	request.Header.Add("Authorization", "Bearer "+o.Token)
	request.Header.Add("Content-Type", "application/json")

	respHandle, err := client.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		err := respHandle.Body.Close()
		if err != nil {
			log.Errorln(err)
			os.Exit(1)
		}
	}()

	content, err := io.ReadAll(respHandle.Body)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	log.Debugln("dfs response:", string(content))

	response := auth.JWTResponse{}
	err = json.Unmarshal(content, &response)
	if err != nil {
		return err
	}

	if response.Code == 200 {
		o.Token = response.Token
		o.Expire = response.Expire
	} else {
		return errors.New("login failed, access denied")
	}

	return nil
}

func (o *ClientOpt) MustRefresh() {
	err := o.Refresh()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

}

const ClientDefaultConfigName = "client"
const ClientDefaultConfigSearchPath0 = "/etc/go-dfs-server"
const ClientDefaultConfigSearchPath1 = "./"

var userHomeDir, _ = os.UserHomeDir()
var ClientDefaultConfigSearchPath2 = path.Join(userHomeDir, ".config/go-dfs-server")
var ClientDefaultConfig = path.Join(userHomeDir, ".config/go-dfs-server/"+ClientDefaultConfigName+".yaml")

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

func (o *ClientOpt) Check() {

}

func Login(cmd *cobra.Command, args []string) {
	log.Debugln("client auth")

	cluster := GetClientOpt()
	vipCfg, err := cluster.Parse(cmd)
	if err != nil {
		if len(args) <= 0 {
			log.Errorln("no url specified")
			os.Exit(1)
		}
		cluster.MustBindURL(args[0])
		cluster.MustBindAuthentication(cmd)
		log.Println("login success")
	} else {
		log.Debugln("%s", cluster)
		if len(args) > 0 {
			cluster.MustBindURL(args[0])
			cluster.MustBindAuthentication(cmd)
		}
		cluster.MustRefresh()
		log.Println("renew token success")
	}

	utils.DumpOption(cluster, vipCfg.GetString("_config"), true)
}
