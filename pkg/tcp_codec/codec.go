package tcp_codec

import (
	"encoding/binary"
	"fmt"
	"net"
)

func Write(ch net.Conn, data []byte) error {

	// write len
	err := binary.Write(ch, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return fmt.Errorf("failed to send message size: %w", err)
	}

	// write msg
	_, err = ch.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func Read(ch net.Conn) ([]byte, error) {

	// read msg len
	var length uint64
	err := binary.Read(ch, binary.BigEndian, &length)
	if err != nil {
		return nil, fmt.Errorf("failed to read message size: %w", err)
	}

	// read msg
	data := make([]byte, length)
	_, err = ch.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return data, nil
}
