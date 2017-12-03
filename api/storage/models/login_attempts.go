package models

import (
	"time"
)

// LoginAttempt represents a Dairycart loginattempt
type LoginAttempt struct {
	ID         uint64    `json:"id"`         // id
	Username   string    `json:"username"`   // username
	Successful bool      `json:"successful"` // successful
	CreatedOn  time.Time `json:"created_on"` // created_on
}
