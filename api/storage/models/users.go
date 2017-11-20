package models

import (
	"time"
)

// User represents a Diarycart user
type User struct {
	ID                    uint64    `json:"id"`                       // id
	FirstName             string    `json:"first_name"`               // first_name
	LastName              string    `json:"last_name"`                // last_name
	Username              string    `json:"username"`                 // username
	Email                 string    `json:"email"`                    // email
	Password              string    `json:"password"`                 // password
	Salt                  []byte    `json:"salt"`                     // salt
	IsAdmin               bool      `json:"is_admin"`                 // is_admin
	PasswordLastChangedOn NullTime  `json:"password_last_changed_on"` // password_last_changed_on
	CreatedOn             time.Time `json:"created_on"`               // created_on
	UpdatedOn             NullTime  `json:"updated_on"`               // updated_on
	ArchivedOn            NullTime  `json:"archived_on"`              // archived_on
}
