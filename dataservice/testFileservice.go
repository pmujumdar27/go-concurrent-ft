package dataservice

import (
	"fmt"
	"io"
	"os"
)

const libpath = "/home/pushkar/Desktop/"

// GetFile requests server a file and receives it
func GetFile(serverConn *Mysock, filename string) {
	filenameBuff := []byte(filename)
	Write(serverConn, filenameBuff)

	var filesize uint64
	ReadObj(serverConn, &filesize)
	fmt.Println("Size of requested file:", filesize)

	file, err := os.OpenFile("recv_"+filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0774)
	handleError(err)
	defer file.Close()

	recvBuf := make([]byte, BUFSIZ)
	recvsz := 0
	size := filesize

	for {
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

	fmt.Println("Downloaded!")
}

// SendFile receives filename and sends the file
func SendFile(clientConn *Mysock) {
	ReadInfo(clientConn, BUFSIZ)
	filename := string(GetBufData(clientConn))
	fmt.Println("Requested filename:", filename)
	filename = libpath + filename
	file, err := os.Open(filename)
	handleError(err)
	defer file.Close()
	fi, _ := file.Stat()
	filesize := uint64(fi.Size())
	Write(clientConn, filesize)

	handleError(err)

	sendBuf := make([]byte, BUFSIZ)
	sentsz := 0
	size := filesize

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
	fmt.Println("File sent!")
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
