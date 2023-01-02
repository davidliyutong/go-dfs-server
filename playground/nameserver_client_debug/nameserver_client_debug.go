package main

import (
	"fmt"
	v1 "go-dfs-server/pkg/nameserver/client/v1"
	"io"
	"os"
)

func main() {
	cli := v1.NewNameServerClient("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzI3NTI3NjEsIm9yaWdfaWF0IjoxNjcyNjY2MzYxLCJ1aWQiOiJiZTdiYmNlODdiNzVhZGQ0In0.Ev_-voBx6NWNaytNBYy4b5b0H65oPONaSYkRKZst_FA", "192.168.105.131", 27903, false)
	//fmt.Println(cli.Ping())
	//_ = cli.BlobLockFile("test.file3", "123")
	//meta, err := cli.BlobReadFileMeta("test.file")

	//reader, err := os.Open("./dataserver_client_debug2.dat")
	//if err != nil {
	//	fmt.Printf(err.Error())
	//}
	//meta, err := cli.BlobWriteChunk("test.file", 2, reader)
	//txt := "Hello,world!"

	writer, err := os.OpenFile("./dataserver_client_debug2.dat", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf(err.Error())
	}
	//buf, err := cli.BlobReadChunk("test.file", 2)
	//_, _ = io.Copy(writer, buf)
	var res string

	res, _ = cli.BlobOpen("test.file", os.O_RDONLY)
	//err := cli.BlobClose("082c9724-16b7-4366-816d-ece301b62cdb")
	//err := cli.BlobWrite(res, 0, bytes.NewBufferString(txt), true)
	_ = cli.BlobSeek(res, 10, 0)

	reader, err := cli.BlobRead(res, 10)
	_, _ = io.Copy(writer, reader)
	_ = cli.BlobSeek(res, 0, 0)
	reader, err = cli.BlobRead(res, 10)
	_, _ = io.Copy(writer, reader)

	_ = cli.BlobClose(res)

	if err == nil {
		fmt.Println(res)
	} else {
		fmt.Printf(err.Error())
	}

}
