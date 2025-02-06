package token

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Payload contains the playground fata of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewPayload creates a now token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	fmt.Printf("NewPayload called with duration: %v\n", duration) // Add this

	now := time.Now()
	expiredAt := now.Add(duration) // Make sure this line exists

	payload := &Payload{
		Username:  username,
		IssuedAt:  now,
		ExpiredAt: expiredAt, // Not the same as IssuedAt!
	}

	fmt.Printf("Payload created with duration %v:\n", duration)
	fmt.Printf("IssuedAt:  %v\n", payload.IssuedAt)
	fmt.Printf("ExpiredAt: %v\n", payload.ExpiredAt)

	return payload, nil
}

// valid check
// Add this method to implement jwt.Claims interface
func (payload *Payload) Valid() error {
	now := time.Now()
	fmt.Printf("Validating token at: %v\n", now)
	fmt.Printf("Token expires at: %v\n", payload.ExpiredAt)

	if now.After(payload.ExpiredAt) {
		fmt.Printf("Token is expired by %v\n", now.Sub(payload.ExpiredAt))
		return ErrExpiredToken
	}
	return nil
}
