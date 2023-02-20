package api

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/zltl/nydus-auth/pkg/id"
)

type CookieClaim struct {
	jwt.RegisteredClaims
}

func (s *State) genCookieToken(uid id.ID) (string, error) {
	now := time.Now()

	claims := CookieClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Second * time.Duration(s.sessionJwtTTL))),
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

func (s *State) validateCookieToken(tokenString string) (*CookieClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CookieClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return s.sessionJwtSecret, nil
		})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	claim := token.Claims.(*CookieClaim)
	if len(claim.Audience) == 0 {
		logrus.Error("require audience")
		return nil, err
	}
	return claim, nil
}
