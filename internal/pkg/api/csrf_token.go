package api

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

type CsrfClaim struct {
	Scope   []string `json:"scope"`
	AllowIp []string `json:"allow_ip"`
	jwt.RegisteredClaims
}

func (s *State) genCsrfJwtToken(scope []string, allowIp []string) (string, error) {
	now := time.Now()

	claims := CsrfClaim{
		AllowIp: allowIp,
		Scope:   scope,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Second * time.Duration(s.csrfJwtTTL))),
			IssuedAt: jwt.NewNumericDate(now),
			Issuer:   "nydus-auth",
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

func (s *State) validateCsrfJwtToken(tokenString string, url string, cip string) (bool, error) {
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
	// check cip
	ok := false
	for _, ip := range claim.AllowIp {
		if ip == cip {
			ok = true
			break
		}
	}
	if !ok {
		logrus.Errorf("got cip %v, want %v", cip, claim.AllowIp)
		return false, errors.New("invalid client ip")
	}

	for _, a := range claim.Audience {
		if a == url {
			return true, nil
		}
	}
	logrus.Errorf("got scope %v, want %v", claim.Audience, url)

	return false, errors.New("invalid audience")
}
