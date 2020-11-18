package main

import (
	"fmt"

	"github.com/pmujumdar27/go-concurrent-ft/dataservice"
)

// BUFSIZ is the buffersize
const BUFSIZ = 1024

// CHUNKSIZE is chunksize
const CHUNKSIZE = uint64(1024 * 32)

// LIBPATH is where downloaded files will be stored
const LIBPATH = "./library"

type responseToClient struct {
	fileName     string
	fileHash     uint64
	availServers map[string]bool
	fileSize     uint64
}

func requestServerList(controllerConn *dataservice.Mysock, filename string) (uint64, map[string]bool) {
	filesize := dataservice.GetFileSize(controllerConn, filename)
	var availServers map[string]bool
	dataservice.ReadObj(controllerConn, &availServers)
	fmt.Println(filesize, availServers)

	return filesize, availServers
}

func handleError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}
