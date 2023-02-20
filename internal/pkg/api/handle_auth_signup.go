package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zltl/nydus-auth/internal/pkg/db"
	"github.com/zltl/nydus-auth/internal/pkg/m"
	"golang.org/x/crypto/bcrypt"
)

// response signup page
func (s *State) handleGetAuthSignup(c *gin.Context) {
	// jwt as CSRF token

	csrfToken, err := s.genCsrfJwtToken([]string{"/auth/signup"})
	if err != nil {
		c.JSON(http.StatusMovedPermanently,
			m.ErrRes{
				Error: "nydus_auth.csrf_token",
				Msg:   "could not generate csrf token:" + err.Error(),
			})
		return
	}

	e := c.Query("e")

	c.HTML(200, "auth/signup", gin.H{
		"csrf_token": csrfToken,
		"error":      e,
	})
}

// handle post request from signup page
func (s *State) handlePostAuthSignup(c *gin.Context) {
	var body m.SignUpReq
	err := c.Bind(&body)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently,
			"/auth/signup?e="+"could not parse request body:"+err.Error())
		return
	}

	logrus.Debugf("--- %+v", body)
	_, err = s.validateCsrfJwtToken(body.CSRFToken, "/auth/signup")
	if err != nil {
		c.Redirect(http.StatusMovedPermanently,
			"/auth/signup?e="+"could not validate csrf token:"+err.Error())
		return
	}

	if len(body.Username) == 0 {
		c.Redirect(http.StatusMovedPermanently,
			"/auth/signup?e="+"username is required")
		return
	}
	if len(body.Password) == 0 {
		c.Redirect(http.StatusMovedPermanently,
			"/auth/signup?e="+"password is required")
		return
	}

	passwordBcrypt, err := bcrypt.GenerateFromPassword(
		[]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently,
			"/auth/signup?e="+"could not bcrypt password:"+err.Error())
		return
	}

	newId := s.idGenerator.Next()

	logrus.WithField("user.id", newId).Infof("new user: %+v", body)

	if err := db.Ctx.AddUser(newId, body.Username,
		string(passwordBcrypt)); err != nil {
		c.Redirect(http.StatusMovedPermanently,
			"/auth/signup?e="+"could not add user:"+err.Error())
		return
	}

	newToken, err := s.genCookieToken(newId)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently,
			"/auth/signup?e="+"could not generate cookie token:"+err.Error())
		return
	}

	// TODO: redirect to main page
	c.SetCookie("_t", newToken, s.sessionExpire, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"msg": "signup success",
	})
}
