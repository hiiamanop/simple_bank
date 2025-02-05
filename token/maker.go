package token

import "time"

// Maker is an interface for making managing tokens
type Maker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)

	// VerifyToken check if the token valid or not
	VerifyToken(token string) (*Payload, error)
}
