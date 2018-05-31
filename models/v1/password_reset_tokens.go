package models

import (
	"time"
)

// PasswordResetToken represents a Dairycart password reset token
type PasswordResetToken struct {
	ID              uint64     `json:"id"`                // id
	UserID          uint64     `json:"user_id"`           // user_id
	Token           string     `json:"token"`             // token
	CreatedOn       time.Time  `json:"created_on"`        // created_on
	ExpiresOn       time.Time  `json:"expires_on"`        // expires_on
	PasswordResetOn *Dairytime `json:"password_reset_on"` // password_reset_on
}

type PasswordResetTokenListResponse struct {
	ListResponse
	PasswordResetTokens []PasswordResetToken `json:"password_reset_tokens"`
}
