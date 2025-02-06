package services

import (
	"context"
	_ "crypto/sha256"
	"fmt"
	"net"
	"pow-server/pkg/pow"
)

const (
	HelloRequest = "Hello"
	Difficulty   = 2 // Количество нулей в начале хэша
)

func NewChallenger(comAdapter ITcpAdapter) *Challenger {
	return &Challenger{
		ConnAdapter: comAdapter,
		POW:         pow.NewScryptPow(Difficulty),
	}
}

type POWVerifier interface {
	GetChallenge() ([]byte, error)
	Verify(challenge []byte, response []byte) error
}

type ITcpAdapter interface {
	Write(ch net.Conn, data []byte) error
	Read(ch net.Conn) ([]byte, error)
}

type Challenger struct {
	ConnAdapter ITcpAdapter
	POW         POWVerifier
}

func (c *Challenger) Challenge(ctx context.Context, ch net.Conn) (bool, error) {

	if err := ctx.Err(); err != nil {
		return false, fmt.Errorf("timeout befor Hello request: %w", err)
	}
	err := c.WaitHello(ch)
	if err != nil {
		return false, err
	}

	if err := ctx.Err(); err != nil {
		return false, fmt.Errorf("timeout befor challange: %w", err)
	}
	challenge, err := c.SendChallenge(ch)
	if err != nil {
		return false, err
	}

	if err := ctx.Err(); err != nil {
		return false, fmt.Errorf("timeout befor verification: %w", err)
	}
	err = c.VerifyChallenge(ch, challenge)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Challenger) WaitHello(ch net.Conn) error {

	buf, err := c.ConnAdapter.Read(ch)
	if err != nil {
		return fmt.Errorf("reading Hello request: %w", err)
	}

	if string(buf) != HelloRequest {
		return fmt.Errorf("hello request check failed")
	}

	return nil
}

func (c *Challenger) SendChallenge(ch net.Conn) ([]byte, error) {
	// Генерируем случайный challenge
	challenge, err := c.POW.GetChallenge()
	if err != nil {
		return nil, fmt.Errorf("failed to generate challenge: %w", err)
	}

	err = c.ConnAdapter.Write(ch, challenge)
	if err != nil {
		return nil, fmt.Errorf("failed to send challenge: %w", err)
	}

	return challenge, nil
}

func (c *Challenger) VerifyChallenge(ch net.Conn, challenge []byte) error {

	response, err := c.ConnAdapter.Read(ch)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	err = c.POW.Verify(challenge, response)
	if err != nil {
		return fmt.Errorf("verification error: %w", err)
	}

	return nil
}
