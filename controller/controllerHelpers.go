package main

import (
	"fmt"
	"net"
	"os"

	"github.com/pmujumdar27/go-concurrent-ft/dataservice"
)

// MAXSERVERS is the number of servers connecting to the controller
const MAXSERVERS = 8

// BUFSIZ is the buffersize
const BUFSIZ = 1024

var controllerPort string
var connID, numConns int64
var numServers int64
var nameHashMap map[string]uint64
var sizeMap map[uint64]uint64
var serverMap map[int64]*serverInfo

type serverInfo struct {
	serverAddr string
	serverID   int64
	serverConn *dataservice.Mysock
	fileSizes  map[uint64]uint64
}

type responseToClient struct {
	fileName     string
	fileHash     uint64
	availServers map[string]bool
	fileSize     uint64
}

func runController(listener net.Listener, controllerAddr string, numServers int64) {
	serverMap = make(map[int64]*serverInfo)
	nameHashMap = make(map[string]uint64)
	sizeMap = make(map[uint64]uint64)

	for cnt := int64(1); cnt <= numServers; cnt++ {
		server, err := listener.Accept()
		handleError(err, "Error in Accepting servers")

		serverConn := dataservice.CreateMysock(server)

		serverAddr := getListenAddr(serverConn)
		numConns++
		fmt.Println("Server connected with ID:", numConns)
		fmt.Println("Server listens file requests at:", serverAddr)

		si := newServerInfo(serverAddr, numConns, serverConn)

		si.fileSizes = make(map[uint64]uint64)
		si.serverAddr = serverAddr
		serverMap[numConns] = si

		initServerFiles(serverConn, si)
		fmt.Println(si.fileSizes)
	}
	fmt.Println(nameHashMap)
	fmt.Println(serverMap)
	fmt.Printf("\nWaiting for clients\n")
	for {
		client, err := listener.Accept()
		handleError(err, "Error in Accepting clients")
		fmt.Println("[+] Client Connected!")
		clientConn := dataservice.CreateMysock(client)
		handleClient(clientConn)
	}
}

func handleClient(clientConn *dataservice.Mysock) {
	filename := dataservice.SendFileSize(clientConn, sizeMap, nameHashMap)
	serverList := getServerData(filename)
	fmt.Println(serverList)
	dataservice.Write(clientConn, serverList)
	dataservice.Write(clientConn, nameHashMap[filename])
}

func getServerData(filename string) map[string]bool {
	fileHash := nameHashMap[filename]
	retmap := make(map[string]bool)
	for _, si := range serverMap {
		fmt.Println(si.fileSizes[fileHash])
		if _, present := si.fileSizes[fileHash]; present {
			retmap[si.serverAddr] = true
		}
	}
	return retmap
}

func initServerFiles(serverConn *dataservice.Mysock, si *serverInfo) {
	// Receive hashval to filename map
	tmp1 := make(map[uint64]string)
	dataservice.ReadObj(serverConn, &tmp1)

	for hashval, filename := range tmp1 {
		nameHashMap[filename] = hashval
	}

	// Receive hashval to filesize map
	m := make(map[uint64]uint64)
	dataservice.ReadObj(serverConn, &m)

	for hashval, filesize := range m {
		si.fileSizes[hashval] = filesize
		sizeMap[hashval] = filesize
	}
}

func getListenAddr(serverConn *dataservice.Mysock) string {
	dataservice.ReadInfo(serverConn, BUFSIZ)
	return string(dataservice.GetBufData(serverConn))
}

func newResponseToClient(fileName string, fileHash, fileSize uint64, availServers map[string]bool) *responseToClient {
	as := make(map[string]bool)
	rtc := responseToClient{fileName: fileName, fileHash: fileHash, fileSize: fileSize, availServers: as}
	return &rtc
}

func newServerInfo(serverAddr string, serverID int64, serverConn *dataservice.Mysock) *serverInfo {
	fileSizes := make(map[uint64]uint64)
	si := serverInfo{serverAddr: serverAddr, serverID: serverID, serverConn: serverConn, fileSizes: fileSizes}
	return &si
}

func handleError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}
