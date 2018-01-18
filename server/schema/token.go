package schema

import "time"

//go:generate reform

// UserToken is a user unique token
//reform:user_token
type UserToken struct {
	ID        uint64    `json:"id" db:"id" unique:"true" reform:"id,pk"`
	Token     string    `json:"token" db:"token" unique:"true" reform:"token"`
	UserID    uint64    `json:"user_id" db:"user_id" reform:"user_id"`
	ExpiredAt time.Time `json:"expired_at" db:"expired_at" reform:"expired_at"`
}
