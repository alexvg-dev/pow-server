package services_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"pow-server/internal/services"
	"testing"
	"time"
)

type mockConn struct {
	buffer   *bytes.Buffer
	writeErr error
}

func (m *mockConn) Read(p []byte) (n int, err error) {
	return m.buffer.Read(p)
}

func (m *mockConn) Write(p []byte) (n int, err error) {
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	return m.buffer.Write(p)
}

type mockPOW struct {
	challenge []byte
	verifyErr error
}

func (m *mockPOW) GetChallenge() ([]byte, error) {
	return m.challenge, nil
}

func (m *mockPOW) Verify(challenge, response []byte) error {
	return m.verifyErr
}

func encodeMessage(msg []byte) []byte {
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.BigEndian, uint64(len(msg)))
	buf.Write(msg)
	return buf.Bytes()
}

func TestChallenger_Challenge(t *testing.T) {
	tests := []struct {
		name      string
		conn      *mockConn
		pow       *mockPOW
		expectErr bool
	}{
		// 1. Hello request received and challenge sent
		{
			name: "success case",
			conn: &mockConn{
				buffer: bytes.NewBuffer(append(encodeMessage([]byte("Hello")), encodeMessage([]byte("response"))...)),
			},
			pow: &mockPOW{
				challenge: []byte("challenge"),
			},
			expectErr: false,
		},
		// 2. Received not Hello request
		{
			name: "invalid hello request",
			conn: &mockConn{
				buffer: bytes.NewBuffer(encodeMessage([]byte("Invalid"))),
			},
			pow:       &mockPOW{},
			expectErr: true,
		},
		// 3. Hello request received, challenge sent, but challenge verification failed
		{
			name: "challenge verification failure",
			conn: &mockConn{
				buffer: bytes.NewBuffer(append(encodeMessage([]byte("Hello")), encodeMessage([]byte("wrong response"))...)),
			},
			pow: &mockPOW{
				challenge: []byte("challenge"),
				verifyErr: errors.New("verification failed"),
			},
			expectErr: true,
		},
		// ...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			challenger := services.Challenger{POW: tt.pow}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			_, err := challenger.Challenge(ctx, tt.conn)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
		})
	}
}
