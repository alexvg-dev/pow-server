package tcp_codec

import (
	"encoding/binary"
	"errors"
)

var (
	ErrSendMsgSize = errors.New("failed to send message size")
	ErrSendMsg     = errors.New("failed to send message")
	ErrReadMsgSize = errors.New("failed to read message size")
	ErrReadMsg     = errors.New("failed to read message")
)

type ReadWriteConn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

func Write(ch ReadWriteConn, data []byte) error {

	// write len
	err := binary.Write(ch, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return ErrSendMsgSize
	}

	// write msg
	_, err = ch.Write(data)
	if err != nil {
		return ErrSendMsg
	}

	return nil
}

func Read(ch ReadWriteConn) ([]byte, error) {

	// read msg len
	var length uint64
	err := binary.Read(ch, binary.BigEndian, &length)
	if err != nil {
		return nil, ErrReadMsgSize
	}

	// read msg
	data := make([]byte, length)
	_, err = ch.Read(data)
	if err != nil {
		return nil, ErrReadMsg
	}

	return data, nil
}
