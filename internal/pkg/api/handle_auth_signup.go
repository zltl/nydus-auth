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
	thisUri := s.externFullUri("/auth/signup")
	csrfToken, err := s.genCsrfJwtToken([]string{thisUri}, []string{c.ClientIP()})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{
			"error": "client not found",
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

	// TODO: delete this log
	logrus.Debugf("--- %+v", body)

	thisUri := s.externFullUri("/auth/signup")
	_, err = s.validateCsrfJwtToken(body.CSRFToken, thisUri, c.ClientIP())
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

	newToken, err := s.genToken(newId,
		[]string{s.externFullUri("/auth/"), s.externFullUri("/api/")},
		s.sessionJwtTTL)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently,
			"/auth/signup?e="+"could not generate cookie token:"+err.Error())
		return
	}

	// TODO: redirect to login page
	c.SetCookie("_t", newToken, s.sessionExpire, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"msg": "signup success",
	})
}
