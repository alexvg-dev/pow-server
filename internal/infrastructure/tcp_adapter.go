package infrastructure

import (
	"encoding/binary"
	"fmt"
	"net"
)

type TcpAdapter struct {
}

func NewTcpAdapter() TcpAdapter {
	return TcpAdapter{}
}

func (tcp TcpAdapter) Write(ch net.Conn, data []byte) error {

	err := binary.Write(ch, binary.BigEndian, uint64(len(data)))
	if err != nil {
		return fmt.Errorf("failed to send message size: %w", err)
	}

	// Отправляем само сообщение
	_, err = ch.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (tcp TcpAdapter) Read(ch net.Conn) ([]byte, error) {
	var length uint64
	err := binary.Read(ch, binary.BigEndian, &length)
	if err != nil {
		return nil, fmt.Errorf("failed to read message size: %w", err)
	}

	// Читаем само сообщение
	data := make([]byte, length)
	_, err = ch.Read(data)
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return data, nil
}
