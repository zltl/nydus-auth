package m

import (
	"time"

	"github.com/zltl/nydus-auth/pkg/id"
)

type TokenReq struct {
	GrantType    string `form:"grant_type"`
	ClientId     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	Code         string `form:"code"`
	RefreshToken string `form:"refresh_token"`
}

type Token struct {
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	ExpiresIn             int    `json:"expires_in"`
	Scope                 string `json:"scope"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

type RefreshToken struct {
	UserID       id.ID
	RefreshToken string
	ExpireAt     time.Time
	ClientID     string
	Scope        string
}
