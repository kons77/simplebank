package token

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"aidanwoods.dev/go-paseto"
)

// PasetoMaker implements paseto v4 Maker interface
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
	implicit     []byte
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey []byte, implicit []byte) (Maker, error) {
	if len(symmetricKey) != 32 { // paseto v4 использует 32-байтовый ключ
		return nil, fmt.Errorf("invalid key size: must be exactly 32 bytes")
	}

	sk, err := paseto.V4SymmetricKeyFromBytes(symmetricKey)
	if err != nil {
		return nil, err
	}

	maker := &PasetoMaker{
		symmetricKey: sk,
		implicit:     implicit,
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	// create paseto token
	token := paseto.NewToken()

	//create uuid for token id
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return "", nil, err
	}

	// add data to the token
	token.Set("id", tokenID.String())
	token.Set("username", username)
	token.SetIssuedAt(time.Now())
	token.SetExpiration(time.Now().Add(duration))

	issAt, err := token.GetIssuedAt()
	if err != nil {
		return "", nil, err
	}
	expr, err := token.GetExpiration()
	if err != nil {
		return "", nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  issAt,
		ExpiresAt: expr,
	}

	return token.V4Encrypt(maker.symmetricKey, maker.implicit), payload, nil
}

// VerifyToken checks if the token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())

	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, token, maker.implicit)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// construct payload from token
	payload, err := getPayloadFromToken(parsedToken)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return payload, nil

}

func getPayloadFromToken(t *paseto.Token) (*Payload, error) {
	id, err := t.GetString("id")
	if err != nil {
		return nil, ErrInvalidToken
	}
	username, err := t.GetString("username")
	if err != nil {
		return nil, ErrInvalidToken
	}
	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return nil, ErrInvalidToken
	}
	expiredAt, err := t.GetExpiration()
	if err != nil {
		return nil, ErrInvalidToken
	}

	return &Payload{
		ID:        uuid.MustParse(id),
		Username:  username,
		IssuedAt:  issuedAt,
		ExpiresAt: expiredAt,
	}, nil
}

var _ Maker = (*PasetoMaker)(nil) // compile-time interface satisfaction check
