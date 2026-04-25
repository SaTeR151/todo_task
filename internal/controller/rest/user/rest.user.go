package user

import (
	"github.com/gin-gonic/gin"
	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/controller/rest/user/validation"
	"github.com/sater-151/todo-list/internal/entity"
)

func (c *UserController) POST(ctx *gin.Context) {
	var req dto.UserPOST
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	userCreate := entity.UserCreate{
		Login:    req.Login,
		Password: req.Password,
	}

	users, err := c.s.UserService.Get(ctx, entity.GetUsersOpts{})
	if err != nil && err != entity.ErrNotFound {
		ctx.JSON(500, err.Error())
		return
	}

	if err := validation.ValidateUserCreate(userCreate, users); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	user, err := c.s.UserService.Create(ctx, userCreate)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, user)
}

func (c *UserController) Auth(ctx *gin.Context) {
	var req dto.UserAuth
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	accessToken, refreshToken, err := c.s.UserService.Auth(ctx, req.Login, req.Password)
	if err != nil {

		if err.IsNotFound() {
			ctx.JSON(404, err.Error())
			return
		}

		if err.IsBadAuth() {
			ctx.JSON(401, err.Error())
			return
		}
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (c *UserController) Get(ctx *gin.Context) {
	userId := ctx.MustGet("user").(string)

	user, err := c.s.UserService.GetByID(ctx, userId)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, user)
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	var req dto.UserPasswordChange
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	userCurrentPassword, err := c.s.UserService.GetPassword(ctx, userID)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	if err := validation.ValidateUserPasswordChange(req, userCurrentPassword); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	newPassword := req.NewPassword

	userUpdate := entity.UserUpdate{
		ID:       userID,
		Password: &newPassword,
	}

	_, err = c.s.UserService.Update(ctx, userUpdate)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.Status(200)
}

func (c *UserController) RefreshToken(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	var req dto.UserRefreshToken
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, err.Error())
		return
	}

	newAccessToken, err := c.s.UserService.RefreshToken(ctx, userID, req.RefreshToken)
	if err != nil {
		if err.IsBadAuth() {
			ctx.JSON(401, err.Error())
			return
		}
		ctx.JSON(500, err.Error())
		return
	}

	ctx.JSON(200, map[string]string{
		"access_token": newAccessToken,
	})
}

func (c *UserController) DELETE(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	err := c.s.UserService.Delete(ctx, userID)
	if err != nil {
		if err == entity.ErrNotFound {
			ctx.JSON(404, err.Error())
			return
		}
		ctx.JSON(500, err.Error())
		return
	}

	ctx.Status(200)
}

func (c *UserController) LogOut(ctx *gin.Context) {
	userID := ctx.MustGet("user").(string)

	nullRefreshToken := ""

	userUpdate := entity.UserUpdate{
		ID:           userID,
		RefreshToken: &nullRefreshToken,
	}

	_, err := c.s.UserService.Update(ctx, userUpdate)
	if err != nil {
		ctx.JSON(500, err.Error())
		return
	}

	ctx.Status(200)
}
