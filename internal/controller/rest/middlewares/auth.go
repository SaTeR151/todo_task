package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/service/user"
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
		if !err.IsEmpty() {
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
		c.Set("user", userID)
		c.Next()
	}
}
