package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	v12 "go-dfs-server/pkg/nameserver/apiserver/blob/v1/model"
	v1 "go-dfs-server/pkg/nameserver/client/v1"
	"io"
	"os"
)

const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzI4NTk3ODAsIm9yaWdfaWF0IjoxNjcyNzczMzgwLCJ1aWQiOiJiZTdiYmNlODdiNzVhZGQ0In0.TsGHM4KDLEXJv75j46pVYPfBVHkYh_jzcNK_8obPCqo"

func testPint() {
	cli := v1.NewNameServerClient(token, "192.168.105.131", 27903, false)
	fmt.Println(cli.Ping())
}

func testClientOpen() {
	cli := v1.NewNameServerClient(token, "192.168.105.131", 27903, false)
	blob, err := cli.Open("test.file", os.O_RDWR)
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infoln(blob)
	log.Infoln("done")
}

func testClientWriteString() {
	cli := v1.NewNameServerClient(token, "192.168.105.131", 27903, false)
	handle, err := cli.Open("test.file", os.O_RDWR)
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infoln(handle)
	written, err := handle.Write([]byte("Hello,world!\n"))
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Debugf("written %d bytes", written)
	err = handle.Close()
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infoln("done")
}

func testClientReadString() {
	cli := v1.NewNameServerClient(token, "192.168.105.131", 27903, false)
	handle, err := cli.Open("test.file", os.O_RDWR)
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infoln(handle)
	buf := make([]byte, 1024)
	n, err := handle.Read(buf)
	if err != nil && err != io.EOF {
		log.Errorln(err)
		return
	}
	log.Debugf("read %s", string(buf[:n]))
	err = handle.Close()
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infoln("done")
}

func testClientWriteFile() {
	cli := v1.NewNameServerClient(token, "192.168.105.131", 27903, false)
	handle, err := cli.Open("test.file2", os.O_RDWR)
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infoln(handle)
	input, err := os.OpenFile("./94152426693025709.rar", os.O_RDONLY, 0664)
	if err != nil {
		log.Errorln(err)
	}

	var total = 0
	for {
		buf := make([]byte, v12.DefaultBlobChunkSize)
		read, err := input.Read(buf)
		if err != nil {
			log.Errorln(err)
			break
		}
		written, err := handle.Write(buf[:read])
		if err != nil {
			log.Errorln(err)
			break
		}
		total += written
	}

	log.Debugf("written %d bytes", total)
	err = handle.Close()
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln("done")
}

func testClientReadFile() {
	cli := v1.NewNameServerClient(token, "192.168.105.131", 27903, false)
	handle, err := cli.Open("test.file2", os.O_RDONLY)
	if err != nil {
		log.Errorln(err)
		return
	}
	log.Infoln(handle)
	output, err := os.OpenFile("./output.rar", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		log.Errorln(err)
	}
	n, _ := io.Copy(output, handle)

	log.Debugf("read %d bytes", n)
	err = handle.Close()
	if err != nil {
		log.Errorln(err)
	}
	log.Infoln("done")
}

func testClientRm() {
	cli := v1.NewNameServerClient(token, "192.168.105.131", 27903, false)
	err := cli.BlobRm("test.file2", false)
	if err != nil {
		log.Errorln(err)
	}

	log.Infoln("done")
}

func main() {
	log.SetLevel(log.DebugLevel)
	//testClientOpen()
	//testClientRm()
	//testClientWriteString()
	//testClientReadString()
	//testClientWriteFile()
	testClientReadFile()

}
