package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/service/user"
	"github.com/sirupsen/logrus"
)

func CheckAuth(userService user.User) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		userID, err := userService.ParseToken(c.Request.Context(), token)
		if err != nil {
			if err.IsBadAuth() {
				c.JSON(401, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			c.JSON(500, gin.H{
				"error": "Internal Server Error",
				"debug": err.Error(),
			})
			c.Abort()

			return
		}

		_, userErr := userService.GetByID(c.Request.Context(), userID)
		if userErr != nil {
			if err == entity.ErrNotFound {
				c.JSON(401, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			c.JSON(500, gin.H{
				"error": "Internal Server Error",
				"debug": userErr.Error(),
			})
			c.Abort()
			return
		}

		logrus.Debugf("Request user ID: %s", userID)
		c.Set("user", userID)
		c.Next()
	}
}
