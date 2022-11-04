package utils

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/sethvargo/go-password/password"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"net"
	"os"
	"path"
	"strings"
)

func MustGenerateAuthKeys() (accessKey string, secretKey string) {
	netInterfaces, err := net.Interfaces()
	var macString string
	if err != nil {
		log.Errorf("fail to get net interfaces: %v", err)
		macString = "00-00-00-00-00-00"
	} else {
		for _, netInterface := range netInterfaces {
			macAddr := netInterface.HardwareAddr.String()
			if len(macAddr) == 0 {
				continue
			}
			macString = strings.ReplaceAll(macAddr, ":", "-")
		}
	}

	accessKey = func(str string) string {
		h := sha1.New()
		h.Write([]byte(str))
		return hex.EncodeToString(h.Sum(nil))[:16]
	}(macString)

	secretKeyGenerator, _ := password.NewGenerator(&password.GeneratorInput{
		LowerLetters: "abcdefghijklmnopqrstuvwxyz",
		UpperLetters: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		Digits:       "0123456789",
		Symbols:      "",
		Reader:       nil,
	})
	secretKey = secretKeyGenerator.MustGenerate(32, 16, 0, false, true)

	return accessKey, secretKey
}

func GetEndpointURL() string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		log.Panicln("net.Interfaces failed, err:", err.Error())
	}

	var ipString = "127.0.0.1"
	for _, netInterface := range netInterfaces {
		if (netInterface.Flags & net.FlagUp) != 0 {
			addrs, _ := netInterface.Addrs()

			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ipString = ipnet.IP.String()
						goto end
					}
				}
			}
		}
	}
end:
	return "dfs://" + ipString + ":27904"
}

func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func AskForConfirmationDefaultYes(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [Y/n]: ", s)

	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" || response == "yes" || response == "" {
		return true
	} else if response == "n" || response == "no" {
		return false
	} else {
		return false
	}

}

func DumpOption(opt interface{}, outputPath string, overwrite bool) {
	buffer, _ := yaml.Marshal(opt)

	parentPath := path.Dir(outputPath)
	if _, err := os.Stat(parentPath); os.IsNotExist(err) {
		err = os.MkdirAll(parentPath, 0644)
		if err != nil {
			log.Panicln("cannot create directory", parentPath)
		}
	}

	if !overwrite {
		if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
			ret := AskForConfirmationDefaultYes("configuration " + outputPath + " already exist, overwrite?")
			if !ret {
				log.Infoln("abort")
				return
			}
		}
	}

	log.Debugln("writing default configuration to", outputPath)
	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	defer func() { _ = f.Close() }()
	if err != nil {
		panic("cannot open " + outputPath + ", check permissions")
	}

	w := bufio.NewWriter(f)
	_, err = w.Write(buffer)
	if err != nil {
		log.Panicln("cannot write configuration", err)
	}
	_ = w.Flush()
	_ = f.Close()

}
