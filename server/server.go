package main

import (
	"fmt"
	"net"
	"os"

	"github.com/pmujumdar27/go-concurrent-ft/dataservice"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("USAGE %s <serverPort> <controllerAddr>\n", os.Args[0])
		os.Exit(1)
	}
	fmt.Println("Hello")
	serverPort = os.Args[1]
	controllerAddr = os.Args[2]
	controller, err := net.Dial("tcp", controllerAddr)
	handleError(err, "Can't reach controller")

	fmt.Println("Connected to the controller", controllerAddr)

	hashData = initializeLibrary(LIBPATH)
	sizeData = initializeFilesizes(LIBPATH, hashData)

	controllerConn := dataservice.CreateMysock(controller)

	listenAddr := "localhost:" + serverPort

	handleController(controllerConn, listenAddr)

	// listen for clients
	listener, err := net.Listen("tcp", listenAddr)
	handleError(err, "Error in listening for clients")
	defer listener.Close()

	for {
		fmt.Println("Listening for clients")
		client, err := listener.Accept()
		handleError(err, "Error in accepting client connection")
		clientConn := dataservice.CreateMysock(client)
		fmt.Println("[+] Client connected", client.RemoteAddr().String())
		handleClient(clientConn)
	}
}
