package services_test

import (
	"errors"
	"net"
	"pow-server/internal/services"
	"testing"
)

type MockTcpAdapter struct {
	WriteFunc func(ch net.Conn, data []byte) error
	ReadFunc  func(ch net.Conn) ([]byte, error)
}

func (m *MockTcpAdapter) Write(ch net.Conn, data []byte) error {
	return m.WriteFunc(ch, data)
}

func (m *MockTcpAdapter) Read(ch net.Conn) ([]byte, error) {
	return m.ReadFunc(ch)
}

type MockPOWVerifier struct {
	GetChallengeFunc func() ([]byte, error)
	VerifyFunc       func(challenge []byte, response []byte) error
}

func (m *MockPOWVerifier) GetChallenge() ([]byte, error) {
	return m.GetChallengeFunc()
}

func (m *MockPOWVerifier) Verify(challenge []byte, response []byte) error {
	return m.VerifyFunc(challenge, response)
}

func TestChallenger_Challenge(t *testing.T) {
	tests := []struct {
		name        string
		adapter     *MockTcpAdapter
		pow         *MockPOWVerifier
		expectError bool
	}{
		{
			name: "successful challenge",
			adapter: &MockTcpAdapter{
				ReadFunc: func(ch net.Conn) ([]byte, error) {
					return []byte(services.HelloRequest), nil
				},
				WriteFunc: func(ch net.Conn, data []byte) error {
					return nil
				},
			},
			pow: &MockPOWVerifier{
				GetChallengeFunc: func() ([]byte, error) {
					return []byte("challenge"), nil
				},
				VerifyFunc: func(challenge []byte, response []byte) error {
					return nil
				},
			},
			expectError: false,
		},
		{
			name: "invalid hello request",
			adapter: &MockTcpAdapter{
				ReadFunc: func(ch net.Conn) ([]byte, error) {
					return []byte("Not hello"), nil
				},
			},
			pow:         &MockPOWVerifier{},
			expectError: true,
		},
		{
			name: "challenge generation failure",
			adapter: &MockTcpAdapter{
				ReadFunc: func(ch net.Conn) ([]byte, error) {
					return []byte(services.HelloRequest), nil
				},
			},
			pow: &MockPOWVerifier{
				GetChallengeFunc: func() ([]byte, error) {
					return nil, errors.New("challenge error")
				},
			},
			expectError: true,
		},
		{
			name: "challenge verification failure",
			adapter: &MockTcpAdapter{
				ReadFunc: func(ch net.Conn) ([]byte, error) {
					return []byte(services.HelloRequest), nil
				},
				WriteFunc: func(ch net.Conn, data []byte) error {
					return nil
				},
			},
			pow: &MockPOWVerifier{
				GetChallengeFunc: func() ([]byte, error) {
					return []byte("challenge"), nil
				},
				VerifyFunc: func(challenge []byte, response []byte) error {
					return errors.New("verification failed")
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			challenger := services.NewChallenger(tt.adapter)
			challenger.POW = tt.pow

			result, err := challenger.Challenge(nil)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError && !result {
				t.Errorf("expected success but got failure")
			}
		})
	}
}
