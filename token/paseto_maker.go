package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMater creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	fmt.Printf("Creating token with duration: %v\n", duration)

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	fmt.Printf("Payload created with IssuedAt: %v, ExpiredAt: %v\n", payload.IssuedAt, payload.ExpiredAt)

	return maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	fmt.Printf("Verifying token at time: %v\n", time.Now())
	fmt.Printf("Token payload: IssuedAt: %v, ExpiredAt: %v\n", payload.IssuedAt, payload.ExpiredAt)

	err = payload.Valid()
	if err != nil {
		fmt.Printf("Token validation failed with error: %v\n", err)
		return nil, err
	}

	return payload, nil
}
