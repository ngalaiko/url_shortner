package schema

import (
	"fmt"
	"time"
)

const (
	maxFirstNameLength = 255
	maxLastNameLength  = 255
)

// User is a user db object
type User struct {
	ID        uint64     `json:"id" db:"id" unique:"true"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
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
