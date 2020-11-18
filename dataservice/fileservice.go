package dataservice

import (
	"fmt"
	"io"
	"os"
)

// GetFileSize sends filename and gets and returns filesize
func GetFileSize(serverConn *Mysock, filename string) uint64 {
	filenameBuff := []byte(filename)
	Write(serverConn, filenameBuff)

	var filesize uint64
	ReadObj(serverConn, &filesize)
	fmt.Println("The filesize is:", filesize)
	return filesize
}

// SendFileSize gets filename and sends filesize to client
func SendFileSize(clientConn *Mysock, sizeData map[uint64]uint64, hashData map[string]uint64) string {
	ReadInfo(clientConn, BUFSIZ)
	filename := string(GetBufData(clientConn))
	fmt.Println("Requested filename:", filename)
	filesize := sizeData[hashData[filename]]
	Write(clientConn, filesize)

	return filename
}

// GetChunk sends chunk request and gets that chunk
func GetChunk(serverConn *Mysock, chunkInfo *ChunkReq, filename string) {
	Write(serverConn, chunkInfo)
	fmt.Println("Error is:", serverConn.err)

	// -----------------------------------------
	file, err := os.OpenFile("recv_"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0774)
	handleError(err)
	defer file.Close()

	recvBuf := make([]byte, BUFSIZ)
	recvsz := 0
	size := chunkInfo.ChunkSize

	offset := int64(chunkInfo.LeftOffset)
	fmt.Println("Offset is:", offset)
	file.Seek(offset, 0)

	for {
		fmt.Println("Bytes remaining:", size)
		if size < BUFSIZ {
			recvBuf = make([]byte, size)
		}
		if size == 0 {
			break
		}

		ReadObj(serverConn, &recvBuf)
		recvsz, err = file.Write(recvBuf)
		handleError(err)
		size = size - uint64(recvsz)
	}

	fmt.Println("Chunk Downloaded!")
}

// SendChunk gets chunk request and sends that chunk
func SendChunk(clientConn *Mysock, hashData map[uint64]string, libpath string) {
	var chunkInfo ChunkReq
	ReadObj(clientConn, &chunkInfo)
	fmt.Println("Got the chunkinfo")
	fmt.Println("The requested chunkInfo is:", chunkInfo)

	filename := hashData[chunkInfo.FileHash]

	filename = libpath + filename
	fmt.Println("Opening", filename)
	file, err := os.Open(filename)
	handleError(err)
	offset := int64(chunkInfo.LeftOffset)
	fmt.Println("Offset is", offset)
	file.Seek(offset, 0)

	// -------------------------------------
	sendBuf := make([]byte, BUFSIZ)
	sentsz := 0
	size := chunkInfo.ChunkSize

	for {
		if size < BUFSIZ {
			sendBuf = make([]byte, size)
		}
		if size == 0 {
			break
		}

		sentsz, err = file.Read(sendBuf)
		if err == io.EOF || size == 0 {
			break
		}

		Write(clientConn, sendBuf)

		size = size - uint64(sentsz)
	}
	fmt.Println("Chunk sent!")
}
