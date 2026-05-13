package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/pkg/utils"
)

type UserService struct {
	users     Repository
	types     TypeCreator
	secretKey string
}

type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *UserService) Create(ctx context.Context, userCreate entity.UserCreate) (res entity.User, err error) {
	defer utils.AddFuncLabel("[service-create-user]", err)

	newUserID, err := s.users.Create(ctx, userCreate)
	if err != nil {
		return
	}

	user, err := s.GetByID(ctx, newUserID)
	if err != nil {
		return
	}

	nullType := entity.TypeCreate{
		UserID: user.ID,
		Name:   "null",
		Color:  "#FFFFFF",
	}

	_, err = s.types.Create(ctx, nullType)
	if err != nil {
		return
	}

	return user, nil
}

func (s *UserService) Update(ctx context.Context, userUpdate entity.UserUpdate) (res entity.User, err error) {
	if userUpdate.Login != nil || userUpdate.RefreshToken != nil {
		if err = s.users.Update(ctx, userUpdate); err != nil {
			return
		}
	}

	if userUpdate.Password != nil {
		err = s.users.UpdatePassword(ctx, userUpdate.ID, *userUpdate.Password)
		if err != nil {
			return
		}
	}

	users, err := s.Get(ctx, entity.GetUsersOpts{ID: userUpdate.ID})
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

	return s.users.Delete(ctx, userID)
}

func (s *UserService) GetRefreshToken(ctx context.Context, userID string) (res string, err error) {
	return s.users.GetRefreshToken(ctx, userID)
}

func (s *UserService) Get(ctx context.Context, opts entity.GetUsersOpts) (users []entity.User, err error) {
	users, err = s.users.Get(ctx, opts)
	if err != nil {
		return
	}

	if len(users) == 0 {
		return nil, entity.ErrNotFound
	}

	return
}

func (s *UserService) GetPassword(ctx context.Context, userID string) (res string, err error) {
	userPassword, err := s.users.GetPassword(ctx, userID)
	if err != nil {
		return
	}

	return userPassword, nil
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

func (s *UserService) ParseToken(ctx context.Context, token string) (string, *entity.AppError) {
	tokenString := strings.TrimPrefix(token, "Bearer ")

	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil || !jwtToken.Valid {
		return "", &entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrBadAuth,
		}
	}

	var userID string

	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok {
		claimUserID, ok := claims["user_id"].(string)
		if !ok || claimUserID == "" {
			return "", &entity.AppError{
				Err:       fmt.Errorf("invalid-token-claims"),
				ErrStatus: entity.ErrBadAuth,
			}
		}
		userID = claimUserID
	}

	return userID, nil
}

func (s *UserService) Auth(ctx context.Context, login, password string) (string, string, *entity.AppError) {
	user, err := s.GetByLogin(ctx, login)
	if err != nil {
		return "", "", &entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	userPassword, err := s.GetPassword(ctx, user.ID)
	if err != nil {
		return "", "", &entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	if !strings.EqualFold(userPassword, password) {
		return "", "", &entity.AppError{
			Err:       fmt.Errorf("invalid-password"),
			ErrStatus: entity.ErrBadAuth,
		}
	}

	expirationTime := time.Now().Add(600 * time.Minute)

	accessClaims := UserClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", "", &entity.AppError{
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
	refreshToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", "", &entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	userUpdate := entity.UserUpdate{
		ID:           user.ID,
		RefreshToken: &refreshToken,
	}

	err = s.users.Update(ctx, userUpdate)
	if err != nil {
		return "", "", &entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	return accessToken, refreshToken, nil
}

func (s *UserService) RefreshToken(ctx context.Context, userID, refreshToken string) (accessToken string, appErr *entity.AppError) {
	userRefreshToken, err := s.GetRefreshToken(ctx, userID)
	if err != nil {
		return "", &entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	if !strings.EqualFold(userRefreshToken, refreshToken) {
		return "", &entity.AppError{
			Err:       fmt.Errorf("invalid-refresh-token"),
			ErrStatus: entity.ErrBadAuth,
		}
	}

	jwtRefreshToken, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil || !jwtRefreshToken.Valid {
		return "", &entity.AppError{
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
		return "", &entity.AppError{
			Err:       err,
			ErrStatus: entity.ErrInternal,
		}
	}

	return
}
