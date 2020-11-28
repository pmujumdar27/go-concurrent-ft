package main

import (
	"fmt"
	"net"
	"os"
	"sync"

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
		// var tmpAddr string

		// for as := range availServers {
		// 	tmpAddr = as
		// }

		// done := uint64(0)
		// server, err := net.Dial("tcp", tmpAddr)
		// handleError(err, "Dial server error")
		// serverConn := dataservice.CreateMysock(server)
		// contFlag := int64(1)
		// dataservice.Write(serverConn, contFlag)
		// for done = uint64(0); done < fileSize; done += CHUNKSIZE {
		// 	reqSize := CHUNKSIZE
		// 	if fileSize-done < CHUNKSIZE {
		// 		reqSize = fileSize - done
		// 		contFlag = 0
		// 	}
		// 	chunkReq := dataservice.CreateChunkReq(fileHash, done, reqSize)
		// 	dataservice.GetChunk(serverConn, chunkReq, filename)
		// 	dataservice.Write(serverConn, contFlag)
		// }
		// dataservice.Close(serverConn)
		// ----------------------------------------------------------------------------------------------

		// return
		numChunks := fileSize / CHUNKSIZE
		if (fileSize % CHUNKSIZE) != 0 {
			numChunks++
		}

		jobs := make(chan *dataservice.ChunkReq, numChunks)
		downloaded := make(chan *dataservice.ChunkReq, numChunks)

		// Create jobs (chunk requests)
		done := uint64(0)

		for done = uint64(0); done < fileSize; done += CHUNKSIZE {

			reqSize := CHUNKSIZE
			if fileSize-done < CHUNKSIZE {
				reqSize = fileSize - done
			}
			chunkReq := dataservice.CreateChunkReq(fileHash, done, reqSize)
			jobs <- chunkReq
		}
		close(jobs)

		var wg sync.WaitGroup

		for as := range availServers {
			server, err := net.Dial("tcp", as)
			handleError(err, "Dial error")
			serverConn := dataservice.CreateMysock(server)
			wg.Add(1)
			// go downloadFromServer(serverConn, jobs, downloaded)
			go func() {
				defer wg.Done()
				for req := range jobs {
					dataservice.Write(serverConn, int64(1))
					chunkName := fmt.Sprintf("%d_%d.tmp", req.FileHash, req.LeftOffset)
					dataservice.GetChunk(serverConn, req, chunkName)
					downloaded <- req
					fmt.Println("Downloaded job put to write channel")
				}
				dataservice.Write(serverConn, int64(0))
				dataservice.Close(serverConn)
				fmt.Println("Done DOwnloading the chunks")
				return
			}()
			wg.Wait()
		}

		writerToMain(filename, downloaded, fileSize)
		// writeToMainDummy(filename, fileSize, fileHash)
		fmt.Println("Num chunks:", numChunks)
		fmt.Println("Done!")
	}
}
