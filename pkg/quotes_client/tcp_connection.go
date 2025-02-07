package quotes_client

import (
	"fmt"
	"net"
	"pow-server/pkg/tcp_codec"
	"time"
)

type TCPConnection struct {
	conn net.Conn
}

func NewTCPConnection(addr string, retries int) (*TCPConnection, error) {

	var (
		conn net.Conn
		err  error
	)

	//
	// Retry policy for short-term network problems
	//
	for attempt := 1; attempt <= retries; attempt++ {
		conn, err = net.Dial("tcp", addr)
		if err == nil {
			break
		}
		if attempt == retries {
			return nil, fmt.Errorf("%w after %d attempts: %v", ErrConnect, attempt, err)
		}
		time.Sleep(time.Second * time.Duration(attempt))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &TCPConnection{conn: conn}, nil

}

func (c *TCPConnection) Write(data []byte) error {
	return tcp_codec.Write(c.conn, data)
}

func (c *TCPConnection) Read() ([]byte, error) {
	return tcp_codec.Read(c.conn)
}

func (c *TCPConnection) Close() error {
	return c.conn.Close()
}
