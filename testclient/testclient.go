package main

import (
	"fmt"
	"net"

	"github.com/pmujumdar27/go-concurrent-ft/dataservice"
)

func main() {
	serverAddr := "localhost:12345"
	server, err := net.Dial("tcp", serverAddr)
	handleError(err)

	serverConn := dataservice.CreateMysock(server)

	defer dataservice.Close(serverConn)

	randomObj := make(map[string]int64)

	randomObj["hello"] = int64(1)
	randomObj["pushkar"] = int64(2)
	randomObj["mujumdar"] = int64(3)
	randomObj["IITGN"] = int64(12345678912345678)

	fmt.Println(randomObj)

	dataservice.Write(serverConn, randomObj)
	dataservice.ReadInfo(serverConn, dataservice.BUFSIZ)
	fmt.Println("The message from server was:", string(dataservice.GetBufData(serverConn)))

	// fmt.Println("Enter filename")
	var filename string
	// var timepass string
	fmt.Scanln(&filename)
	fmt.Println(filename)
	// filename = "testserver.go"

	filesize := dataservice.GetFileSize(serverConn, filename)
	chunkOne := dataservice.CreateChunkReq(uint64(0), uint64(0), uint64(1024*32))
	chunkTwo := dataservice.CreateChunkReq(uint64(0), uint64(1024*32), filesize-uint64(1024*32))

	dataservice.GetChunk(serverConn, chunkOne, filename)
	dataservice.GetChunk(serverConn, chunkTwo, filename)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
