package quotes_client

import (
	"errors"
	"net"
	"pow-server/pkg/tcp_codec"
)

var (
	ErrConnect       = errors.New("can`t connect to server")
	ErrSendHello     = errors.New("can`t send Hello request")
	ErrReadChallenge = errors.New("can`t read challenge")
	ErrSendSolution  = errors.New("can`t send solution")
	ErrReadQuote     = errors.New("can`t read quote from server")
)

func NewClient(addr string, powSolver PowSolver) *Client {
	return &Client{
		Addr:      addr,
		PowSolver: powSolver,
	}
}

type PowSolver interface {
	SolveChallenge(challenge []byte) ([]byte, error)
}

type Client struct {
	Addr      string
	PowSolver PowSolver
}

func (c *Client) GetQuote() (string, error) {

	//
	// Connecting to server
	//
	conn, err := net.Dial("tcp", c.Addr)
	if err != nil {
		return "", ErrConnect
	}
	defer conn.Close()

	//
	// Sending Hello request
	//
	err = tcp_codec.Write(conn, []byte("Hello"))
	if err != nil {
		return "", ErrSendHello
	}

	//
	// Read challenge
	//
	buf, err := tcp_codec.Read(conn)
	if err != nil {
		return "", ErrReadChallenge
	}

	//
	// Solving challenge
	//
	response, err := c.PowSolver.SolveChallenge(buf)

	//
	// Send solution to server
	//
	err = tcp_codec.Write(conn, response)
	if err != nil {
		return "", ErrSendSolution
	}

	//
	// Reading quote
	//
	buf, err = tcp_codec.Read(conn)
	if err != nil {
		return "", ErrReadQuote
	}

	return string(buf), nil
}
