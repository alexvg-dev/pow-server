package quotes_client

import (
	"errors"
	"net"
	"pow-server/internal/infrastructure"
)

var (
	ErrConnect       = errors.New("can`t connect to server")
	ErrSendHello     = errors.New("can`t send Hello request")
	ErrReadChallenge = errors.New("can`t read challenge")
	ErrSendSolution  = errors.New("can`t send solution")
	ErrReadQuote     = errors.New("can`t read quote from server")
)

func NewClient(addr string, powSolver PowSolver, adapter infrastructure.TcpAdapter) *Client {
	return &Client{
		Addr:             addr,
		PowSolver:        powSolver,
		TransportAdapter: adapter,
	}
}

type PowSolver interface {
	SolveChallenge(challenge []byte) ([]byte, error)
}

type Client struct {
	Addr             string
	PowSolver        PowSolver
	TransportAdapter infrastructure.TcpAdapter
}

func (c *Client) GetQuote() (string, error) {

	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return "", ErrConnect
	}
	defer conn.Close()

	err = c.TransportAdapter.Write(conn, []byte("Hello"))
	if err != nil {
		return "", ErrSendHello
	}

	// Читаем challenge от сервера
	buf, err := c.TransportAdapter.Read(conn)
	if err != nil {
		return "", ErrReadChallenge
	}

	// Генерируем ответ на challenge
	response, err := c.PowSolver.SolveChallenge(buf)

	// Отправляем ответ на challenge
	err = c.TransportAdapter.Write(conn, response)
	if err != nil {
		return "", ErrSendSolution
	}

	// Читаем ответ от сервера
	buf, err = c.TransportAdapter.Read(conn)
	if err != nil {
		return "", ErrReadQuote
	}

	return string(buf), nil
}
