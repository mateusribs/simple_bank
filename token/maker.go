package token

import "time"

// manage tokens
type Maker interface {
	//creates new token for specific username and valid duration
	CreateToken(username string, duration time.Duration) (string, error)

	// check if token is valid
	VerifyToken(token string) (*Payload, error)
}