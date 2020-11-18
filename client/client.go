package main

import (
	"fmt"
	"net"
	"os"

	"github.com/pmujumdar27/go-concurrent-ft/dataservice"
)

var controllerAddr string

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE: %s <controllerIP:port>\n", os.Args[0])
		os.Exit(1)
	}
	controllerAddr = os.Args[1]

	controller, err := net.Dial("tcp", controllerAddr)
	handleError(err, "Error in Dialing to Controller")

	fmt.Println("Connected to controller at", controllerAddr)

	controllerConn := dataservice.CreateMysock(controller)

	defer controller.Close()

	for {
		fmt.Println("Enter the name of the file you want")
		var filename string
		fmt.Scanln(&filename)
		fileSize, availServers := requestServerList(controllerConn, filename)
		var fileHash uint64
		dataservice.ReadObj(controllerConn, &fileHash)
		fmt.Println("FileHash:", fileHash)

		// The following code is temporary chunked file transfer ... replace it with concurrent file transfer
		// -----------------------------------------------------------------------------------------------
		var tmpAddr string

		for as := range availServers {
			tmpAddr = as
		}

		done := uint64(0)

		for done = uint64(0); done < fileSize; done += CHUNKSIZE {
			server, err := net.Dial("tcp", tmpAddr)
			handleError(err, "Dial server error")
			serverConn := dataservice.CreateMysock(server)
			reqSize := CHUNKSIZE
			if fileSize-done < CHUNKSIZE {
				reqSize = fileSize - done
			}
			chunkReq := dataservice.CreateChunkReq(fileHash, done, reqSize)
			dataservice.GetChunk(serverConn, chunkReq, filename)
		}
		// ----------------------------------------------------------------------------------------------
	}
}
