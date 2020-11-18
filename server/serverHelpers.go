package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/OneOfOne/xxhash"
	"github.com/pmujumdar27/go-concurrent-ft/dataservice"
)

// BUFSIZ is the buffersize
const BUFSIZ = 1024

// LIBPATH is the path of the library
const LIBPATH = "./library/"

// HASHSIZE is hashsize
const HASHSIZE = 64

// CHUNKSIZE is chunksize
const CHUNKSIZE = 1024 * 32

var serverPort, controllerAddr string
var hashData map[uint64]string
var sizeData map[uint64]uint64

func handleClient(clientConn *dataservice.Mysock) {
	dataservice.SendChunk(clientConn, hashData, LIBPATH)
}

func handleController(controllerConn *dataservice.Mysock, listenAddr string) {
	// send listenAddr
	fmt.Println("Listen addr:", listenAddr)
	dataservice.Write(controllerConn, []byte(listenAddr))

	// send hashval to filename map
	dataservice.Write(controllerConn, hashData)

	// send hashval to filesize map
	dataservice.Write(controllerConn, sizeData)
}

func initializeFilesizes(libpath string, hashData map[uint64]string) map[uint64]uint64 {
	fmt.Println("Initializing file sizes ...")
	m := make(map[uint64]uint64)
	for hashval, filename := range hashData {
		fi, err := os.Stat(LIBPATH + "/" + filename)
		handleError(err, "Cant get file details")
		m[hashval] = uint64(fi.Size())
	}
	fmt.Println("The filesize data is:")
	fmt.Println(m)
	return m
}

func initializeLibrary(libpath string) map[uint64]string {
	fmt.Println("Initializing Library ...")
	files, err := ioutil.ReadDir(libpath)
	handleError(err, "Error in initializing library")

	m := make(map[uint64]string)

	for _, f := range files {
		fName := f.Name()
		m[gethash(fName)] = fName
	}
	return m
}

func gethash(fName string) uint64 {
	h := xxhash.New64()
	r := strings.NewReader(fName)
	io.Copy(h, r)
	// fmt.Println("xxhash.Backend:", xxhash.Backend)
	checksumval := h.Sum64()
	fmt.Printf("File checksum of %s: %d\n", fName, checksumval)
	return h.Sum64()
}

func handleError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}
