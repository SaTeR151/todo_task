package column

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/column/validation"
	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/pkg/utils"
)

func (c *ColumnController) POST(ctx *gin.Context) {
	boardID := ctx.MustGet("board").(string)

	var columnPOST dto.ColumnPOST
	if err := ctx.ShouldBindJSON(&columnPOST); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	columns, err := c.columns.GetByBoardID(ctx, boardID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := validation.ValidateColumnCreate(columns, columnPOST); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	columnCreate := entity.ColumnCreate{
		Name:       columnPOST.Name,
		BoardID:    boardID,
		OderNumber: columnPOST.OrderNumber,
	}

	column, err := c.columns.CreateColumn(ctx, columnCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, column)
}

func (c *ColumnController) GET(ctx *gin.Context) {
	boardID := ctx.MustGet("board").(string)

	var columnGETQuery dto.ColumnGETUri
	if err := ctx.ShouldBindUri(&columnGETQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	columns, err := c.columns.GetByID(ctx, boardID, columnGETQuery.ColumnID)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, columns)
}

func (c *ColumnController) LIST(ctx *gin.Context) {
	boardID := ctx.MustGet("board").(string)

	columns, err := c.columns.GetByBoardID(ctx, boardID)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(http.StatusOK, entity.Columns{})
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, columns)
}

func (c *ColumnController) DELETE(ctx *gin.Context) {
	boardID := ctx.MustGet("board").(string)

	var columnDELETEQuery dto.ColumnDELETEUri
	if err := ctx.ShouldBindUri(&columnDELETEQuery); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := c.columns.DeleteColumn(ctx, boardID, columnDELETEQuery.ColumnID)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *ColumnController) PATCH(ctx *gin.Context) {
	boardID := ctx.MustGet("board").(string)

	var uri dto.ColumnPATCHUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	var req dto.ColumnPATCH
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	columns, err := c.columns.GetByBoardID(ctx, boardID)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	otherColumns := utils.ColumnsExcept(columns, uri.ColumnID)

	if err := validation.ValidateColumnUpdate(otherColumns, req); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	columnUpdate := entity.ColumnUpdate{
		ID:          uri.ColumnID,
		Name:        req.Name,
		OrderNumber: req.OrderNumber,
	}

	column, err := c.columns.UpdateColumn(ctx, boardID, columnUpdate)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, column)
}

func (c *ColumnController) SWAP(ctx *gin.Context) {
	boardID := ctx.MustGet("board").(string)

	var columnSWAP dto.ColumnSWAP
	if err := ctx.ShouldBindJSON(&columnSWAP); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	columnA, err := c.columns.GetByID(ctx, boardID, columnSWAP.ColumnIDA)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	columnB, err := c.columns.GetByID(ctx, boardID, columnSWAP.ColumnIDB)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if err := validation.ValidateColumnSwap(columnA, columnB); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = c.columns.SwapColumns(ctx, boardID, columnSWAP.ColumnIDA, columnSWAP.ColumnIDB)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(http.StatusNotFound, err.Error())
			return
		}
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Status(http.StatusOK)
}
