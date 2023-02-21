package api

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zltl/nydus-auth/internal/pkg/db"
	"github.com/zltl/nydus-auth/internal/pkg/m"
)

func (s *State) handleAuthToken(c *gin.Context) {
	var body m.TokenReq
	err := c.Bind(&body)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error":             "invalid_request",
			"error_description": "could not parse request body:" + err.Error(),
		})
		return
	}

	if body.GrantType != "authorization_code" &&
		body.GrantType != "refresh_token" {

		logrus.Error("grant_type must be authorization_code")
		c.JSON(400, gin.H{
			"error":             "unsupported_grant_type",
			"error_description": "grant_type must be authorization_code",
		})
		return
	}

	// verify client secret
	client, err := db.Ctx.GetAuthClient(body.ClientId)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error":             "invalid_client",
			"error_description": "client not found",
		})
		return
	}
	if client.Secret != body.ClientSecret {
		logrus.Error("client secret mismatch")
		c.JSON(400, gin.H{
			"error":             "unauthorized_client",
			"error_description": "client secret mismatch",
		})
		return
	}

	if body.GrantType == "authorization_code" {
		s.handleTokenCode(&body, c)
	} else {
		s.handleTokenRefresh(&body, c)
	}
}

func (s *State) handleTokenRefresh(body *m.TokenReq, c *gin.Context) {
	// verify refresh token
	refs, err := db.Ctx.QueryRefreshToken(body.ClientId, body.RefreshToken)
	if err != nil || len(refs) == 0 {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error":             "invalid_grant",
			"error_description": "refresh token error",
		})
		return
	}

	mref := refs[0]
	anyNotExpire := false
	for _, ref := range refs {
		if !ref.ExpireAt.Before(time.Now()) {
			anyNotExpire = true
			mref = ref
			break
		}
	}
	if !anyNotExpire {
		logrus.Error("refresh token expired")
		c.JSON(400, gin.H{
			"error":             "invalid_grant",
			"error_description": "refresh token expired",
		})
		return
	}

	scopeList := strings.Split(mref.Scope, ",")
	tokenStr, err := s.genToken(mref.UserID, scopeList, s.sessionJwtTTL)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error":             "internal_error",
			"error_description": "could not generate token",
		})
		return
	}

	c.JSON(
		200,
		m.Token{
			AccessToken:           tokenStr,
			TokenType:             "Bearer",
			ExpiresIn:             s.sessionJwtTTL,
			Scope:                 mref.Scope,
			RefreshToken:          mref.RefreshToken,
			RefreshTokenExpiresIn: int(time.Until(mref.ExpireAt).Seconds()),
		},
	)
}

func (s *State) handleTokenCode(body *m.TokenReq, c *gin.Context) {
	// verify code
	code, err := db.Ctx.QueryCode(body.ClientId, body.Code)
	if err != nil {
		logrus.Error(err)
		c.JSON(400, gin.H{
			"error":             "invalid_grant",
			"error_description": "code error",
		})
		return
	}
	// TODO: delete expires code
	scopeList := strings.Split(code.Scope, ",")
	tokenStr, err := s.genToken(code.UserID, scopeList, s.sessionJwtTTL)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error":             "internal_error",
			"error_description": "could not generate token",
		})
		return
	}

	// generate refresh
	reftk := make([]byte, 32)
	_, err = rand.Read(reftk)
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error":             "internal_error",
			"error_description": "could not generate refresh token",
		})
		return
	}
	reftkStr := hex.EncodeToString(reftk)

	// save refresh token
	err = db.Ctx.SaveRefreshToken(&m.RefreshToken{
		UserID:       code.UserID,
		ClientID:     code.ClientID,
		RefreshToken: reftkStr,
		ExpireAt: time.Now().Add(
			time.Duration(s.sessionJwtRefreshExpire) * time.Second),
		Scope: code.Scope,
	})
	if err != nil {
		logrus.Error(err)
		c.JSON(500, gin.H{
			"error":             "internal_error",
			"error_description": "could not save refresh token",
		})
		return
	}

	c.JSON(
		200,
		m.Token{
			AccessToken:           tokenStr,
			TokenType:             "Bearer",
			ExpiresIn:             s.sessionJwtTTL,
			Scope:                 code.Scope,
			RefreshToken:          reftkStr,
			RefreshTokenExpiresIn: s.sessionJwtRefreshExpire,
		},
	)

}
