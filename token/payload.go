package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// different types of error returned by the VerifyToken function
var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// Payload contains the payload data of the token
type Payload struct {
	jwt.RegisteredClaims // embedding

	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	// IssuedAt  time.Time `json:"issued_at"` // when the token is created
	// ExpiredAt time.Time // the time at which the token expires

}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	payload := &Payload{
		ID:       tokenID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}

	return payload, nil
}

/*
// Valid checks if the token payload is valid or not
func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
*/
