package tasktype

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/controller/rest/type/validation"
	"github.com/sater-151/todo-list/internal/entity"
)

func (c *TypeController) GET(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	var req dto.TypeGetUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	types, err := c.types.GetByUserID(ctx, userID)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(200, entity.Types{})
			return
		}
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, types)
}

func (c *TypeController) LIST(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	types, err := c.types.GetByUserID(ctx, userID)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(200, entity.Types{})
			return
		}
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, types)
}

func (c *TypeController) POST(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	var req dto.TypePOST
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	typeCreate := entity.TypeCreate{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}

	types, err := c.types.GetByUserID(ctx, userID)
	if err != nil && err != entity.ErrNotFound {
		ctx.JSON(500, err.Error())
		return
	}

	if err := validation.ValidateTypeCreate(typeCreate, types); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	taskType, err := c.types.Create(ctx, typeCreate)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, taskType)
}

func (c *TypeController) PATCH(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	var req dto.TypePATCH
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	var uri dto.TypePATCHUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	typeUpdate := entity.TypeUpdate{
		ID:    uri.TypeID,
		Name:  req.Name,
		Color: req.Color,
	}

	types, err := c.types.GetByUserID(ctx, userID)
	if err != nil && err != entity.ErrNotFound {
		ctx.JSON(500, err.Error())
		return
	}

	var otherTypes entity.Types
	for _, t := range types {
		if t.ID != typeUpdate.ID {
			otherTypes = append(otherTypes, t)
		}
	}

	if err := validation.ValidateTypeUpdate(typeUpdate, otherTypes); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	updatedType, err := c.types.Update(ctx, userID, typeUpdate)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, updatedType)
}

func (c *TypeController) DELETE(ctx *gin.Context) {
	var uri dto.TypeDELETEUri
	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	err := c.types.Delete(ctx, uri.TypeID)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.Status(204)
}
