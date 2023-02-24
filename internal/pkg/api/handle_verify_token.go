package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zltl/nydus-auth/internal/pkg/m"
	"github.com/zltl/nydus-auth/pkg/id"
)

func (s *State) handleVerifyToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		var head m.AuthorizeHeader

		err := c.ShouldBindHeader(&head)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":             "invalid_request",
				"error_description": "could not parse request header:" + err.Error(),
			})
			return
		}

		idTokenHeader := strings.Split(head.Authorization, "Bearer ")
		if len(idTokenHeader) != 2 {
			logrus.Error("auth header invalid: %s", head.Authorization)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":             "invalid_request",
				"error_description": "could not parse request header: Bearer",
			})
			return
		}
		token := idTokenHeader[1]

		urireq := s.externFullUri(c.Request.URL.Path)
		claim, err := s.validateToken(token, urireq)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":             "invalid_request",
				"error_description": "wrong access_token",
			})
			return
		}

		uid, err := id.Parse(claim.Subject)
		if err != nil {
			logrus.Error(err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":             "invalid_request",
				"error_description": "invalid token",
			})
			return
		}

		logrus.Debugf("uid=%s", uid.String())
		c.Set("uid", uid)

		c.Next()
	}

}
