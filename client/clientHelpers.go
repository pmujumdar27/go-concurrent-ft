package main

import (
	"fmt"
	"io"
	"os"

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

// func downloadFile(addressChunksMap map[string][]*dataservice.ChunkReq, filename string) {
// 	downloadedQueue := list.New()
// 	downloadChan := make(chan *dataservice.ChunkReq)
// 	totChunks := 0
// 	for sa, crs := range addressChunksMap {
// 		totChunks += len(crs)
// 		server, err := net.Dial("tcp", sa)
// 		handleError(err, "Dailing Error")
// 		serverConn := dataservice.CreateMysock(server)
// 		go getChunksFromServer(serverConn, crs, downloadedQueue)
// 	}
// 	writtenChunks := 0
// 	for true {
// 		if writtenChunks == totChunks {
// 			break
// 		}
// 		currChunk := downloadedQueue.Front()
// 	}
// }

// func getChunksFromServer(serverConn *dataservice.Mysock, crs []*dataservice.ChunkReq, downloadedQueue *list.List) {
// 	for cr := range crs {
// 		realCR := crs[cr]
// 		filename := fmt.Sprintf("%d", realCR.FileHash) + fmt.Sprintf("%d", realCR.LeftOffset) + ".dwn"
// 		dataservice.GetChunk(serverConn, realCR, filename)
// 		downloadedQueue.PushBack(realCR)
// 	}
// }

func downloadFromServer(serverConn *dataservice.Mysock, jobs <-chan *dataservice.ChunkReq, downloaded chan<- *dataservice.ChunkReq) {
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
}

func writerToMain(filename string, downloaded <-chan *dataservice.ChunkReq, fileSize uint64) {
	finalFile, err := os.OpenFile("recv_"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	handleError(err, "Open Error")
	defer finalFile.Close()
	// finalFile.Seek(int64(fileSize)-int64(1), 0)
	// _, err = finalFile.Write([]byte{0})
	// handleError(err, "Write error")
	cnt := int64(0)
	for req := range downloaded {
		cnt++
		chunkName := fmt.Sprintf("recv_%d_%d.tmp", req.FileHash, req.LeftOffset)
		fmt.Println("Writting chunk", chunkName, cnt, len(downloaded))

		chunkFile, err := os.Open(chunkName)
		handleError(err, "Chunk open error")

		off := int64(req.LeftOffset)
		finalFile.Seek(off, 0)

		readBuf := make([]byte, BUFSIZ)
		readSz := 0
		size := req.ChunkSize

		for {
			if size < BUFSIZ {
				readBuf = make([]byte, size)
			}
			if size == 0 {
				break
			}

			readSz, err = chunkFile.Read(readBuf)
			if err == io.EOF || size == 0 {
				break
			}

			if size == CHUNKSIZE && off == int64(0) {
				fmt.Println(string(readBuf))
			}

			bytesWritten, err := finalFile.Write(readBuf[:readSz])
			handleError(err, "Write error")

			if bytesWritten != readSz {
				fmt.Println("Something is fishy", bytesWritten, readSz)
			}

			readBuf = make([]byte, size)

			// if size == CHUNKSIZE && off == int64(0) {
			// 	finalFile.Seek(int64(0), 0)
			// 	tmp, tmperr := finalFile.Read(readBuf)
			// 	handleError(tmperr, "Verify read error")
			// 	fmt.Println("Bytes read", tmp)
			// 	fmt.Println(string(readBuf))
			// }

			size = size - uint64(readSz)
		}
		chunkFile.Close()
		os.Remove(chunkName)
	}
	fmt.Println("Chunks written", cnt)
}

// func writeToMainDummy(filename string, fileSize uint64, fileHash uint64) {
// 	finalFile, err := os.OpenFile("recv_"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
// 	handleError(err, "open dest err")
// 	lo := uint64(0)
// 	for {
// 		chunkName := fmt.Sprintf("recv_%d_%d.dwn", fileHash, lo)
// 		chunkSize := CHUNKSIZE
// 		if fileSize-uint64(lo) < CHUNKSIZE {
// 			chunkSize = fileSize - lo
// 		}
// 		chunkFile, err := os.Open(chunkName)
// 		handleError(err, "chunk open")
// 		sendbuf := make([]byte, chunkSize)
// 		readsz, _ := chunkFile.Read(sendbuf)
// 		writesz, _ := finalFile.Write(sendbuf)
// 		if readsz != writesz {
// 			fmt.Println("Something fishy", readsz, writesz)
// 		}
// 		lo += uint64(writesz)
// 		if lo >= fileSize {
// 			break
// 		}
// 	}
// }

func handleError(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		panic(err)
	}
}
