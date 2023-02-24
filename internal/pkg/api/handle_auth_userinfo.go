package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zltl/nydus-auth/internal/pkg/db"
	"github.com/zltl/nydus-auth/pkg/id"
)

func (s *State) handleAuthUserinfo(c *gin.Context) {
	uido, ok := c.Get("uid")
	if !ok {
		c.JSON(400, gin.H{
			"error":             "invalid_request",
			"error_description": "could not get user be token",
		})
		return
	}
	uid := uido.(id.ID)

	user, err := db.Ctx.GetUserInfo(uid)
	if err != nil {
		c.JSON(400, gin.H{
			"error":             "invalid_request",
			"error_description": "could not get user info:" + err.Error(),
		})
		return
	}
	c.JSON(200, user)
}
