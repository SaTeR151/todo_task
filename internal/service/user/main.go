package user

import (
	"context"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres"
)

type User interface {
	Create(ctx context.Context, userCreate entity.UserCreate) (res entity.User, err error)
	Update(ctx context.Context, userUpdate entity.UserUpdate) (res entity.User, err error)
	Delete(ctx context.Context, userID string) error
	Get(ctx context.Context, opts entity.GetUsersOpts) (users []entity.User, err error)
	GetRefreshToken(ctx context.Context, userID string) (res string, err error)
	GetPassword(ctx context.Context, userID string) (res string, err error)
	GetByID(ctx context.Context, userID string) (res entity.User, err error)
	GetByLogin(ctx context.Context, login string) (res entity.User, err error)
	ParseToken(ctx context.Context, token string) (string, entity.AppError)
	Auth(ctx context.Context, userID, password string) (accessToken, refreshToken string, appErr entity.AppError)
	RefreshToken(ctx context.Context, userID, refreshToken string) (accessToken string, appErr entity.AppError)
}

func New(repo *postgres.Repository, secretKey string) User {
	return &UserService{
		repo:      repo,
		secretKey: secretKey,
	}
}
