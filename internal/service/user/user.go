package user

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres"
	"github.com/sater-151/todo-list/pkg/utils"
)

type UserService struct {
	repo      *postgres.Repository
	secretKey string
}

type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *UserService) Create(ctx context.Context, userCreate entity.UserCreate) (res entity.User, err error) {
	defer utils.AddFuncLabel("[service-create-user]", err)

	newUserID, err := s.repo.User.Create(ctx, userCreate)
	if err != nil {
		return
	}

	users, err := s.repo.User.Get(ctx, entity.GetUsersOpts{ID: newUserID})
	if err != nil {
		return
	}

	return users[0], nil
}

func (s *UserService) Update(ctx context.Context, userUpdate entity.UserUpdate) (res entity.User, err error) {
	if err = s.repo.User.Update(ctx, userUpdate); err != nil {
		return
	}

	users, err := s.repo.User.Get(ctx, entity.GetUsersOpts{ID: userUpdate.ID})
	if err != nil {
		return
	}

	return users[0], nil
}

func (s *UserService) Delete(ctx context.Context, userID string) error {
	_, err := s.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	return s.repo.User.Delete(ctx, userID)
}

func (s *UserService) GetRefreshToken(ctx context.Context, userID string) (res string, err error) {
	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return
	}

	return user.RefreshToken, nil
}

func (s *UserService) Get(ctx context.Context, opts entity.GetUsersOpts) (users []entity.User, err error) {
	users, err = s.repo.User.Get(ctx, opts)
	if err != nil {
		return
	}

	if len(users) == 0 {
		return nil, entity.ErrNotFound
	}

	return
}

func (s *UserService) GetPassword(ctx context.Context, userID string) (res string, err error) {
	user, err := s.GetByID(ctx, userID)
	if err != nil {
		return
	}

	return user.Password, nil
}

func (s *UserService) GetByID(ctx context.Context, userID string) (res entity.User, err error) {
	users, err := s.Get(ctx, entity.GetUsersOpts{ID: userID})
	if err != nil {
		return
	}

	return users[0], nil
}

func (s *UserService) GetByLogin(ctx context.Context, login string) (res entity.User, err error) {
	users, err := s.Get(ctx, entity.GetUsersOpts{Login: login})
	if err != nil {
		return
	}

	return users[0], nil
}

func (s *UserService) ParseToken(ctx context.Context, token string) (string, entity.AppError) {
	tokenString := strings.TrimPrefix(token, "Bearer ")

	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil || !jwtToken.Valid {
		return "", entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrBadAuth,
		}
	}

	var userID string

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		userID = claims["user_id"].(string)
	}

	return userID, entity.NilError()
}

func (s *UserService) Auth(ctx context.Context, login, password string) (accessToken, refreshToken string, appErr entity.AppError) {
	user, err := s.GetByLogin(ctx, login)
	if err != nil {
		return "", "", entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	if !strings.EqualFold(user.Password, password) {
		return "", "", entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrBadAuth,
		}
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	accessClaims := UserClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", "", entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	expirationTime = time.Now().Add(7 * 24 * time.Hour)

	refreshClaims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", "", entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	userUpdate := entity.UserUpdate{
		ID:           user.ID,
		RefreshToken: &refreshToken,
	}

	err = s.repo.User.Update(ctx, userUpdate)
	if err != nil {
		return "", "", entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	return
}

func (s *UserService) RefreshToken(ctx context.Context, userID, refreshToken string) (accessToken string, appErr entity.AppError) {
	userRefreshToken, err := s.GetRefreshToken(ctx, userID)
	if err != nil {
		return "", entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	if !strings.EqualFold(userRefreshToken, refreshToken) {
		return "", entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrBadAuth,
		}
	}

	jwtRefreshToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil || !jwtRefreshToken.Valid {
		return "", entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrBadAuth,
		}
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	accessClaims := UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	return
}
