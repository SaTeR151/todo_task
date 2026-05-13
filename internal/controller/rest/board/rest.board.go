package board

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/board/validation"
	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/pkg/utils"
)

// @Summary Create board
// @Tags board
// @Accept json
// @Produce json
// @Param board body entity.BoardCreate true "Board"
// @Success 200 {object} entity.Board
// @Router /boards [post]
func (c *BoardController) POST(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	var boardCreateDTO dto.BoardPOST
	if err := ctx.ShouldBindJSON(&boardCreateDTO); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := validation.ValidateBoardCreate(boardCreateDTO); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	boardCreate := entity.BoardCreate{
		UserID: userID,
		Name:   boardCreateDTO.Name,
	}

	board, err := c.boards.Create(ctx, boardCreate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, board)
}

// @Summary Get board
// @Tags board
// @Accept json
// @Produce json
// @Param board_id query string true "Board ID"
// @Success 200 {object} entity.Board
// @Router /boards [get]
func (c *BoardController) GET(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	var query dto.BoadGETUri

	if err := ctx.ShouldBindUri(&query); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	board, err := c.boards.GetByID(ctx, userID, query.BoardID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, board)
}

// @Summary Get boards
// @Tags board
// @Accept json
// @Produce json
// @Success 200 {object} entity.Boards
// @Router /boards [get]
func (c *BoardController) LIST(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	boards, err := c.boards.GetByUserID(ctx, userID)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(200, entity.Boards{})
			return
		}
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, boards)
}

func (c *BoardController) PATCH(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	var uri dto.BoardPATCHUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var req dto.BoardPATCH
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	boards, err := c.boards.GetByUserID(ctx, userID)
	if err != nil && err != entity.ErrNotFound {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	otherBoards := utils.BoardsExcept(boards, uri.BoardID)

	if err := validation.ValidateBoardUpdate(req, otherBoards); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	boardUpdate := entity.BoardUpdate{
		ID:   uri.BoardID,
		Name: req.Name,
	}

	board, err := c.boards.Update(ctx, userID, boardUpdate)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, board)
}

func (c *BoardController) DELETE(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	var boardDELETEQuery dto.BoardDELETEUri
	if err := ctx.ShouldBindUri(&boardDELETEQuery); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := c.boards.Delete(ctx, userID, boardDELETEQuery.BoardID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, nil)
}
