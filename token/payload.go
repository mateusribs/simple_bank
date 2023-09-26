package token

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("token has expired")
	ErrInvalidToken = errors.New("token is invalid")
)

// contains the payload data of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// GetAudience implements jwt.Claims.
func (*Payload) GetAudience() (jwt.ClaimStrings, error) {
	panic("unimplemented")
}

// GetExpirationTime implements jwt.Claims.
func (*Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	panic("unimplemented")
}

// GetIssuedAt implements jwt.Claims.
func (*Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	panic("unimplemented")
}

// GetIssuer implements jwt.Claims.
func (*Payload) GetIssuer() (string, error) {
	panic("unimplemented")
}

// GetNotBefore implements jwt.Claims.
func (*Payload) GetNotBefore() (*jwt.NumericDate, error) {
	panic("unimplemented")
}

// GetSubject implements jwt.Claims.
func (*Payload) GetSubject() (string, error) {
	panic("unimplemented")
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()

	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}
