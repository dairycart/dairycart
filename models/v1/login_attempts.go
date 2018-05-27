package models

import (
	"time"
)

// LoginAttempt represents a Dairycart login attempt
type LoginAttempt struct {
	ID         uint64    `json:"id"`         // id
	Username   string    `json:"username"`   // username
	Successful bool      `json:"successful"` // successful
	CreatedOn  time.Time `json:"created_on"` // created_on
}

// LoginAttemptCreationInput is a struct to use for creating LoginAttempts
type LoginAttemptCreationInput struct {
	Username   string `json:"username,omitempty"`   // username
	Successful bool   `json:"successful,omitempty"` // successful
}

// LoginAttemptUpdateInput is a struct to use for updating LoginAttempts
type LoginAttemptUpdateInput struct {
	Username   string `json:"username,omitempty"`   // username
	Successful bool   `json:"successful,omitempty"` // successful
}

type LoginAttemptListResponse struct {
	ListResponse
	LoginAttempts []LoginAttempt `json:"login_attempts"`
}
