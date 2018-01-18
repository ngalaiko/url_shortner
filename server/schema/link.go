package schema

import (
	"time"
)

//go:generate reform

// Link is a link db object
//easyjson:json
//reform:links
type Link struct {
	ID         uint64     `json:"id" db:"id" unique:"true" reform:"id,pk"`
	UserID     uint64     `json:"user_id" db:"user_id" reform:"user_id"`
	URL        string     `json:"url" db:"url" reform:"url"`
	ShortURL   string     `json:"short_url" db:"short_url" reform:"short_url"`
	ViewsLimit uint64     `json:"views_limit" db:"views_limit" reform:"views_limit"`
	Views      uint64     `json:"views" db:"views" reform:"views"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at" reform:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at" db:"deleted_at" reform:"deleted_at"`
}

// Valid returns true if link is valid
func (l *Link) Valid() bool {
	switch {
	case l.DeletedAt != nil:
		return false
	case l.ViewsLimit > 0 && l.Views >= l.ViewsLimit:
		return false
	default:
		return true
	}
}

// Anonim returns if link crated by not authorizated user
func (l *Link) Anonim() bool {
	return l.UserID == uint64(0)
}
