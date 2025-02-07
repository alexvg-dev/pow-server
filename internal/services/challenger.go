package services

import (
	"context"
	_ "crypto/sha256"
	"fmt"
	"pow-server/pkg/pow"
	"pow-server/pkg/tcp_codec"
)

const (
	HelloRequest = "Hello"
	Difficulty   = 2 // Количество нулей в начале хэша
)

func NewChallenger() *Challenger {
	return &Challenger{
		POW: pow.NewScryptPow(Difficulty),
	}
}

type ReadWriteConn interface {
	Read(b []byte) (n int, err error)
	Write(b []byte) (n int, err error)
}

type POWVerifier interface {
	GetChallenge() ([]byte, error)
	Verify(challenge []byte, response []byte) error
}

type Challenger struct {
	POW POWVerifier
}

func (c *Challenger) Challenge(ctx context.Context, ch ReadWriteConn) (bool, error) {

	err := c.WaitHello(ctx, ch)
	if err != nil {
		return false, err
	}

	challenge, err := c.SendChallenge(ctx, ch)
	if err != nil {
		return false, err
	}

	err = c.VerifyChallenge(ctx, ch, challenge)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Challenger) WaitHello(ctx context.Context, ch ReadWriteConn) error {

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("timeout befor Hello request: %w", err)
	}

	buf, err := tcp_codec.Read(ch)
	if err != nil {
		return fmt.Errorf("reading Hello request: %w", err)
	}

	if string(buf) != HelloRequest {
		return fmt.Errorf("hello request check failed")
	}

	return nil
}

func (c *Challenger) SendChallenge(ctx context.Context, ch ReadWriteConn) ([]byte, error) {

	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("timeout befor challange: %w", err)
	}

	challenge, err := c.POW.GetChallenge()
	if err != nil {
		return nil, fmt.Errorf("failed to generate challenge: %w", err)
	}

	err = tcp_codec.Write(ch, challenge)
	if err != nil {
		return nil, fmt.Errorf("failed to send challenge: %w", err)
	}

	return challenge, nil
}

func (c *Challenger) VerifyChallenge(ctx context.Context, ch ReadWriteConn, challenge []byte) error {

	if err := ctx.Err(); err != nil {
		return fmt.Errorf("timeout befor verification: %w", err)
	}

	response, err := tcp_codec.Read(ch)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	err = c.POW.Verify(challenge, response)
	if err != nil {
		return fmt.Errorf("verification error: %w", err)
	}

	return nil
}
