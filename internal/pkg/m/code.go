package m

import (
	"time"

	"github.com/zltl/nydus-auth/pkg/id"
)

type AuthCode struct {
	ClientID string
	UserID   id.ID
	Scope    string
	ExpireAt time.Time
	Code     string
}
