package main

import (
	"fmt"
	"net"
)

func main() {
	listenAddr := "localhost:12345"
	server, err := net.Listen("tcp", listenAddr)
	handleError(err)

	client, err := net.Dial("tcp", listenAddr)
	handleError(err)

	conn, err := server.Accept()
	handleError(err)

	fmt.Println(conn.RemoteAddr().String(), client.LocalAddr())

}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
