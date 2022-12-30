package utils

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sethvargo/go-password/password"
	log "github.com/sirupsen/logrus"
	v1 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	"gopkg.in/yaml.v3"
	"io"
	"math/rand"
	"net"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
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
	return ipString + ":27904"
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
	fileInfo, err := os.Stat(parentPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(parentPath, 0700)
		if err != nil {
			log.Errorln("cannot create directory", parentPath)
			log.Exit(1)
		}
	}

	if os.IsPermission(err) || fileInfo.Mode() != 0700 {
		err = os.Chmod(parentPath, 0700)
		if err != nil {
			log.Errorln("cannot read director", parentPath)
			log.Exit(1)
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
	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
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

func JoinSubPathSafe(path string, subpath string) (string, error) {
	subpath = filepath.Clean(subpath)
	if strings.HasPrefix(subpath, "..") {
		return "", errors.New(fmt.Sprintf("subpath %s is not allowed", subpath))
	} else {
		return filepath.Join(path, subpath), nil
	}
}

func GetChunkPath(path string, chunkID int64) string {
	return filepath.Join(path, strconv.FormatInt(chunkID, 10)+".dat")
}

func GetMetaPath(path string) string {
	return filepath.Join(path, "meta.json")
}

func GetLockPath(path string) string {
	return filepath.Join(path, ".lock")
}

func GetFileLockState(path string) bool {
	lockPath := GetLockPath(path)
	_, err := os.Stat(lockPath)
	return err == nil
}

func GetFileState(path string) bool {
	metaPath := GetMetaPath(path)
	_, err := os.Stat(metaPath)
	return err == nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetFileMD5(path string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string
	//Open the passed argument and check for any error
	file, err := os.Open(path)
	if err != nil {
		return returnMD5String, err
	}
	//Tell the program to call the following function when the current function returns
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	//Open a new hash interface to write to
	hash := md5.New()
	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	//Info the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func GetBufferMD5(data []byte) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string
	hash := md5.New()
	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, bytes.NewReader(data)); err != nil {
		return returnMD5String, err
	}
	//Info the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]
	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func MustGenerateUUID() string {
	ul := uuid.NewV4()
	return ul.String()
}

var (
	once sync.Once
)

func SelectRandomNFromArray(arr []interface{}, n int64) []interface{} {
	once.Do(func() {
		rand.Seed(time.Now().Unix())
	})
	ret := make([]interface{}, n)
	perm := rand.Perm(len(arr))[:3]
	for i, randIndex := range perm {
		ret[i] = arr[randIndex]
	}
	return ret
}

func HasError(arr []error) bool {
	for _, err := range arr {
		if err != nil {
			return true
		}
	}
	return false
}

func MaxInt64(x int64, y int64) int64 {
	if x >= y {
		return x
	} else {
		return y
	}
}

func MinInt64(x int64, y int64) int64 {
	if x < y {
		return x
	} else {
		return y
	}
}

func GetChunkID(offset int64) int64 {
	return offset / v1.DefaultBlobChunkSize
}

func GetChunkOffset(offset int64) int64 {
	return offset % v1.DefaultBlobChunkSize
}
