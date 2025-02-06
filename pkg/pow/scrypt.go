package pow

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/scrypt"
	"strings"
)

const (
	ChallengeSize = 32
)

var (
	ErrCantGenerateChallenge = errors.New("can`t generate challenge")
	ErrCantGenerateScryptKey = errors.New("can`t generate scrypt key")
	ErrVerificationFailed    = errors.New("solution verification failed")
)

type ScryptPow struct {
	Difficulty uint
}

func NewScryptPow(difficulty uint) *ScryptPow {
	return &ScryptPow{
		Difficulty: difficulty,
	}
}

func (pow *ScryptPow) GetChallenge() ([]byte, error) {
	challenge := make([]byte, ChallengeSize)
	_, err := rand.Read(challenge)
	if err != nil {
		return nil, ErrCantGenerateChallenge
	}

	return challenge, nil
}

func (pow *ScryptPow) SolveChallenge(challenge []byte) ([]byte, error) {

	nonce := 0 // Начинаем с нуля
	for {
		// Преобразуем nonce в байты
		nonceBytes := []byte(fmt.Sprintf("%d", nonce))

		// Вычисляем хэш с использовоНогом
		hash, err := scrypt.Key(nonceBytes, challenge, 16384, 8, 1, 32)
		if err != nil {
			return nil, ErrCantGenerateScryptKey
		}

		hashStr := hex.EncodeToString(hash)

		// Проверяем, соответствует ли хэш условиям
		if strings.HasPrefix(hashStr, strings.Repeat("0", int(pow.Difficulty))) {
			return nonceBytes, nil
		}

		nonce++ // Увеличиваем nonce для следующей итерации
	}
}

func (pow *ScryptPow) Verify(challenge, response []byte) error {

	hash, err := scrypt.Key(response, challenge, 16384, 8, 1, 32)
	if err != nil {
		return ErrCantGenerateScryptKey
	}
	hashStr := hex.EncodeToString(hash)

	if !strings.HasPrefix(hashStr, strings.Repeat("0", int(pow.Difficulty))) {
		return ErrVerificationFailed
	}

	return nil
}
