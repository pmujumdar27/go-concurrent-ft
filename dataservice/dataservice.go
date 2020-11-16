package dataservice

import (
	"bytes"
	"encoding/gob"
	"net"
)

// BUFSIZ is the buffersize
const BUFSIZ = uint64(1024)

// Mysock is the wrapper structure for socket
type Mysock struct {
	socket net.Conn
	buf    []byte
	reader *gob.Decoder
	writer *gob.Encoder
	err    error
}

// ChunkReq is the structure containing a chunk request
type ChunkReq struct {
	FileHash   uint64
	LeftOffset uint64
	ChunkSize  uint64
}

// CreateMysock creates a mysock object and returns a pointer to it
func CreateMysock(connection net.Conn) *Mysock {
	sock := Mysock{
		socket: connection,
		buf:    make([]byte, BUFSIZ),
		reader: gob.NewDecoder(connection),
		writer: gob.NewEncoder(connection),
		err:    nil,
	}
	return &sock
}

// CreateChunkReq creates a chunk and returns a pointer to it
func CreateChunkReq(fileHash uint64, leftOffset uint64, chunkSize uint64) *ChunkReq {
	chunk := ChunkReq{
		FileHash:   fileHash,
		LeftOffset: leftOffset,
		ChunkSize:  chunkSize,
	}
	return &chunk
}

// Close closes the connection and Mysock.err captures the error
func Close(sock *Mysock) {
	sock.err = sock.socket.Close()
}

// ReadInfo reads lenBuf bytes of data from the socket
func ReadInfo(sock *Mysock, lenBuf uint64) {
	buf := make([]byte, lenBuf)
	sock.err = sock.reader.Decode(&buf)
	sock.buf = bytes.Trim(buf, "\x00")
}

// ReadObj reads an object from the socket
func ReadObj(sock *Mysock, x interface{}) { //interface {} is the empty interface, used to handle unkonwn types
	sock.err = sock.reader.Decode(x)
}

// Write sends the required object throught the socket
func Write(sock *Mysock, x interface{}) {
	sock.err = sock.writer.Encode(x)
}

// GetBufData returns data in the socket's buffer
func GetBufData(sock *Mysock) []byte {
	return sock.buf
}
