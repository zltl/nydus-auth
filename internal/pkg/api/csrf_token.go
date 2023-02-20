package api

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type CsrfClaim struct {
	jwt.RegisteredClaims
}

func (s *State) genCsrfJwtToken(aud []string) (string, error) {
	now := time.Now()

	claims := CsrfClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Second * time.Duration(s.csrfJwtTTL))),
			IssuedAt: jwt.NewNumericDate(now),
			Issuer:   "nydus-auth",
			Audience: aud,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(s.csrfJwtAlgo), claims)
	token.Header["kid"] = s.csrfJwtKeyId
	token.Header["alg"] = s.csrfJwtAlgo
	token.Header["typ"] = "JWT"

	tokenString, err := token.SignedString(s.csrfJwtSecret)
	if err != nil {
		logrus.Error(err)
	}

	return tokenString, err
}

func (s *State) validateCsrfJwtToken(tokenString string, aud string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CsrfClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return s.csrfJwtSecret, nil
		})
	if err != nil {
		logrus.Error(err)
		return false, err
	}
	claim := token.Claims.(*CsrfClaim)
	if len(claim.Audience) == 0 {
		logrus.Error("require audience")
		return false, err
	}
	// claim.Audience
	for _, a := range claim.Audience {
		if a == aud {
			return true, nil
		}
	}
	logrus.Errorf("got audience %v, want %v", claim.Audience, aud)

	return false, errors.New("invalid audience")
}
