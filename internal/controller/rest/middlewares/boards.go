package middlewares

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/service/board"
)

func CheckBoard(boardService board.Board) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID, exists := c.Get("user")
		if !exists {
			c.JSON(400, gin.H{"error": "user not found"})
			c.Abort()
			return
		}

		boardID := c.Param("board")
		board, err := boardService.GetByID(c.Request.Context(), userID.(string), boardID)
		if err != nil {
			if errors.Is(err, entity.ErrNotFound) {
				c.JSON(404, gin.H{"error": "board not found"})
				c.Abort()
				return
			}
			c.JSON(400, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set("board", board.ID)
		c.Next()
	}
}
