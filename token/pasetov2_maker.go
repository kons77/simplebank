package token // !! V2 OUTDATED VERSION OF PASETO

import (
	"fmt"
	"time"

	"golang.org/x/crypto/chacha20poly1305"

	"github.com/o1egl/paseto"
)

/* From paseto.io , it looks like o1egl 's paseto package stop the support to v2
and now there are v3 and v4 go-paseto supports v3/v4 */

// PasetoMaker is a PASETO v2 (outdated) token maker
type PasetoMaker2 struct {
	paseto       *paseto.V2
	simmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker2(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker2{
		paseto:       paseto.NewV2(),
		simmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker2) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.simmetricKey, payload, nil)
	return token, payload, err

}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker2) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.simmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
