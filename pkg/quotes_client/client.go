package quotes_client

import (
	"errors"
	"fmt"
)

const (
	ConnectionRetriesCount = 3
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
	conn, err := NewTCPConnection(c.Addr, ConnectionRetriesCount)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	//
	// Sending Hello request
	//
	err = conn.Write([]byte("Hello"))
	if err != nil {
		return "", fmt.Errorf("send hello: %w", ErrSendHello)
	}

	//
	// Read challenge
	//
	challenge, err := conn.Read()
	if err != nil {
		return "", fmt.Errorf("read challenge: %w", ErrReadChallenge)
	}

	//
	// Solving challenge
	//
	solution, err := c.PowSolver.SolveChallenge(challenge)
	if err != nil {
		return "", fmt.Errorf("solve challenge: %w", err)
	}

	//
	// Send solution to server
	//
	err = conn.Write(solution)
	if err != nil {
		return "", fmt.Errorf("send solution: %w", ErrSendSolution)
	}

	//
	// Reading quote
	//
	quote, err := conn.Read()
	if err != nil {
		return "", fmt.Errorf("read quote: %w", ErrReadQuote)
	}

	return string(quote), nil
}
