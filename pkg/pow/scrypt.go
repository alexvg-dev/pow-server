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
	ChallengeSize     = 32
	ScryptParamN      = 16384
	ScryptParamR      = 8
	ScryptParamP      = 1
	ScryptParamKeyLen = 32
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

	// Random challenge is good enough against reply-attacks
	// But it can be enhanced by extending length of challenge or adding timestamp
	_, err := rand.Read(challenge)
	if err != nil {
		return nil, ErrCantGenerateChallenge
	}

	return challenge, nil
}

func (pow *ScryptPow) SolveChallenge(challenge []byte) ([]byte, error) {

	nonce := 0
	for {
		nonceBytes := []byte(fmt.Sprintf("%d", nonce))

		hash, err := pow.genKey(nonceBytes, challenge)
		if err != nil {
			return nil, ErrCantGenerateScryptKey
		}

		hashStr := hex.EncodeToString(hash)

		if strings.HasPrefix(hashStr, strings.Repeat("0", int(pow.Difficulty))) {
			return nonceBytes, nil
		}

		nonce++
	}
}

func (pow *ScryptPow) Verify(challenge, response []byte) error {

	hash, err := pow.genKey(response, challenge)
	if err != nil {
		return ErrCantGenerateScryptKey
	}
	hashStr := hex.EncodeToString(hash)

	if !strings.HasPrefix(hashStr, strings.Repeat("0", int(pow.Difficulty))) {
		return ErrVerificationFailed
	}

	return nil
}

func (pow *ScryptPow) genKey(nonceBytes, challenge []byte) ([]byte, error) {
	hash, err := scrypt.Key(nonceBytes, challenge, ScryptParamN, ScryptParamR, ScryptParamP, ScryptParamKeyLen)
	if err != nil {
		return nil, ErrCantGenerateScryptKey
	}

	return hash, nil
}
