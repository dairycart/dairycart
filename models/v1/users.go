package models

import (
	"time"
)

// User represents a Dairycart user
type User struct {
	ID                    uint64     `json:"id"`                       // id
	FirstName             string     `json:"first_name"`               // first_name
	LastName              string     `json:"last_name"`                // last_name
	Username              string     `json:"username"`                 // username
	Email                 string     `json:"email"`                    // email
	Password              string     `json:"password"`                 // password
	Salt                  []byte     `json:"salt"`                     // salt
	IsAdmin               bool       `json:"is_admin"`                 // is_admin
	PasswordLastChangedOn *Dairytime `json:"password_last_changed_on"` // password_last_changed_on
	CreatedOn             time.Time  `json:"created_on"`               // created_on
	UpdatedOn             *Dairytime `json:"updated_on"`               // updated_on
	ArchivedOn            *Dairytime `json:"archived_on"`              // archived_on
}

// UserCreationInput is a struct to use for creating Users
type UserCreationInput struct {
	FirstName string `json:"first_name,omitempty"` // first_name
	LastName  string `json:"last_name,omitempty"`  // last_name
	Username  string `json:"username,omitempty"`   // username
	Email     string `json:"email,omitempty"`      // email
	Password  string `json:"password,omitempty"`   // password
	IsAdmin   bool   `json:"is_admin,omitempty"`   // is_admin
}

// UserUpdateInput is a struct to use for updating Users
type UserUpdateInput struct {
	FirstName       string `json:"first_name,omitempty"`       // first_name
	LastName        string `json:"last_name,omitempty"`        // last_name
	Username        string `json:"username,omitempty"`         // username
	Email           string `json:"email,omitempty"`            // email
	NewPassword     string `json:"new_password,omitempty"`     // password
	CurrentPassword string `json:"current_password,omitempty"` // password
	IsAdmin         bool   `json:"is_admin,omitempty"`         // is_admin
}

type UserListResponse struct {
	ListResponse
	Users []User `json:"users"`
}
