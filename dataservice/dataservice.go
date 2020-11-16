package dataservice

import (
	"encoding/gob"
	"net"
)

type mysock struct {
	socket net.Conn
	buf    []byte
	reader *gob.Decoder
	writer *gob.Encoder
	err    error
}
