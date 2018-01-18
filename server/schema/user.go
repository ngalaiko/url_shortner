package schema

import (
	"fmt"
	"time"
)

const (
	maxFirstNameLength = 255
	maxLastNameLength  = 255
)

//go:generate reform

// User is a user db object
//easyjson:json
//reform:users
type User struct {
	ID         uint64     `json:"id" db:"id" unique:"true" reform:"id,pk"`
	FirstName  string     `json:"first_name" db:"first_name" reform:"first_name"`
	LastName   string     `json:"last_name" db:"last_name" reform:"last_name"`
	FacebookID string     `json:"facebook_id" db:"facebook_id" reform:"facebook_id"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at" reform:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at" db:"deleted_at" reform:"deleted_at"`
}

// Validate validates user structure fields
func (u *User) Validate() error {
	switch {
	case len(u.FirstName) > maxFirstNameLength:
		return fmt.Errorf("first name %s is longer than %d", u.FirstName, maxFirstNameLength)
	case len(u.LastName) > maxLastNameLength:
		return fmt.Errorf("last name %s is longer than %d", u.FirstName, maxFirstNameLength)
	default:
		return nil
	}
}
