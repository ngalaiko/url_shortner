package schema

import (
	"time"
)

// Link is a link db object
//easyjson:json
type Link struct {
	ID        uint64     `json:"id" db:"id" unique:"true"`
	UserID    uint64     `json:"user_id" db:"user_id"`
	URL       string     `json:"url" db:"url"`
	ShortURL  string     `json:"short_url" db:"short_url"`
	Clicks    uint64     `json:"clicks" db:"clicks"`
	Views     uint64     `json:"views" db:"views"`
	ExpiredAt time.Time  `json:"expired_at" db:"expired_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	DeletedAt *time.Time `json:"deleted_at" db:"deleted_at"`
}
