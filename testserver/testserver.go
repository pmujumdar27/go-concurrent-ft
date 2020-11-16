package main

import (
	"fmt"
	"net"

	"github.com/pmujumdar27/go-concurrent-ft/dataservice"
)

func main() {
	serverAddr := "localhost:12345"

	server, err := net.Listen("tcp", serverAddr)
	handleError(err)
	fmt.Println("Listening on ", serverAddr)
	fmt.Println("Buffersize is:", dataservice.BUFSIZ)

	for {
		conn, err := server.Accept()
		handleError(err)
		fmt.Println("[+] Client connected:", conn.RemoteAddr().String())
		clientConn := dataservice.CreateMysock(conn)
		go handleClient(clientConn)
	}
}

func handleClient(clientConn *dataservice.Mysock) {
	defer dataservice.Close(clientConn)
	m := make(map[string]int64)
	dataservice.ReadObj(clientConn, &m)
	fmt.Println("Object from client:", m)
	dataservice.Write(clientConn, []byte("Got your object!"))

	filename := dataservice.SendFileSize(clientConn)

	dataservice.SendChunk(clientConn, filename)
	dataservice.SendChunk(clientConn, filename)

	fmt.Println("Done!")
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
