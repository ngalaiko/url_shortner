package schema

import "time"

// UserToken is a user unique token
type UserToken struct {
	ID        uint64    `json:"id" db:"id" unique:"true"`
	Token     string    `json:"token" db:"token" unique:"true"`
	UserID    uint64    `json:"user_id" db:"user_id"`
	ExpiredAt time.Time `json:"expired_at" db:"expired_at"`
}
