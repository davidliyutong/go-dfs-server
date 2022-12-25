package config

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strconv"
	"time"
)

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
		portInt = NameServerDefaultPort
	} else {
		portInt = NameServerDefaultPort
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

func (o *ClientAuthOpt) BindAuthentication(cmd *cobra.Command) error {
	var accessKey, secretKey string
	if cmd != nil {
		accessKey, _ = cmd.Flags().GetString("accessKey")
		secretKey, _ = cmd.Flags().GetString("secretKey")
	} else {
		accessKey = ""
		secretKey = ""
	}

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

	keyPairIsValid := accessKey != "" && secretKey != ""
	if keyPairIsValid {
		o.AccessKey = accessKey
		o.SecretKey = secretKey
		o.Token = "<dummy_token>"
		o.Expire = time.UnixMicro(0)
	} else {
		o.AccessKey = ""
		o.SecretKey = ""
		o.Token = ""
		o.Expire = time.UnixMicro(0)
	}

	return nil
}

func (o *ClientAuthOpt) MustBindAuthentication(cmd *cobra.Command) {
	err := o.BindAuthentication(cmd)
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}

}
