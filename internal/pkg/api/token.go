package api

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/zltl/nydus-auth/pkg/id"
)

type CookieClaim struct {
	Scope []string `json:"scope"`
	jwt.RegisteredClaims
}

func (s *State) genToken(uid id.ID, scope []string, expires int) (string, error) {
	now := time.Now()

	claims := CookieClaim{
		Scope: scope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Second * time.Duration(expires))),
			IssuedAt: jwt.NewNumericDate(now),
			Issuer:   "nydus-auth",
			Subject:  uid.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(string(s.sessionJwtAlgo)), claims)
	token.Header["kid"] = s.sessionJwtKeyId
	token.Header["alg"] = s.sessionJwtAlgo
	token.Header["typ"] = "JWT"

	tokenString, err := token.SignedString(s.sessionJwtSecret)
	if err != nil {
		logrus.Error(err)
	}

	return tokenString, err
}

func (s *State) validateToken(tokenString string, uri string) (*CookieClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CookieClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return s.sessionJwtSecret, nil
		})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	claim := token.Claims.(*CookieClaim)
	if len(claim.Scope) == 0 {
		logrus.Error("require scope")
		return nil, err
	}

	// check uri with aud
	ok := false
	for _, scope := range claim.Scope {
		if strings.HasPrefix(uri, scope) {
			ok = true
			break
		}
	}
	if !ok {
		logrus.Error("scope not contain this uri")
		return nil, err
	}

	return claim, nil
}
