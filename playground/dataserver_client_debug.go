package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	v1 "go-dfs-server/pkg/dataserver/client/v1"
)

func main() {
	cli := v1.NewDataServerClient("", "192.168.105.131", 27904, false)
	//fmt.Println(cli.Ping())
	//_ = cli.BlobLockFile("test.file3", "123")
	//meta, err := cli.BlobReadFileMeta("test.file")

	//reader, err := os.Open("./dataserver_client_debug.go")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//meta, err := cli.BlobWriteChunk("test.file", 2, reader)

	//writer, err := os.OpenFile("./dataserver_client_debug2.dat", os.O_CREATE|os.O_WRONLY, 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//buf, err := cli.BlobReadChunk("test.file", 2)
	//_, _ = io.Copy(writer, buf)

	res, err := cli.SysVolume()

	if err == nil {
		log.Println(res)
	} else {
		fmt.Printf(err.Error())
	}

}
