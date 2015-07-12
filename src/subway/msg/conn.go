package msg

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

func readMsgShared(c net.Conn) (buffer []byte, err error) {
	var sz int64
	err = binary.Read(c, binary.LittleEndian, &sz)
	if err != nil {
		return
	}

	buffer = make([]byte, sz)
	n, err := c.Read(buffer)

	if err != nil {
		return
	}

	if int64(n) != sz {
		err = errors.New(fmt.Sprintf("Expected to read %d bytes, but only read %d", sz, n))
		return
	}

	return
}

func ReadMsg(c net.Conn) (msg Message, err error) {
	buffer, err := readMsgShared(c)
	if err != nil {
		return
	}

	return Unpack(buffer)
}

func ReadMsgInto(c net.Conn, msg Message) (err error) {
	buffer, err := readMsgShared(c)
	if err != nil {
		return
	}
	return UnpackInto(buffer, msg)
}

func WriteMsg(c net.Conn, msg interface{}) (err error) {
	buffer, err := Pack(msg)
	if err != nil {
		return
	}

	// c.Debug("Writing message: %s", string(buffer))
	err = binary.Write(c, binary.LittleEndian, int64(len(buffer)))

	if err != nil {
		return
	}

	if _, err = c.Write(buffer); err != nil {
		return
	}

	return nil
}
