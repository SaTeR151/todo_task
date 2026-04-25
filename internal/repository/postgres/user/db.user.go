package user

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres/models"
	"github.com/sater-151/todo-list/pkg/pgutils"
	"github.com/sater-151/todo-list/pkg/utils"
	"github.com/sirupsen/logrus"
)

func (s *UserStorage) Get(ctx context.Context, opts entity.GetUsersOpts) (res entity.Users, err error) {
	defer utils.AddFuncLabel("[repo-get-users]", err)

	query := s.queryBuilder.Select(
		"id",
		"login",
	).From(fmt.Sprintf("%s.%s AS users", s.scheme, TABLE_USERS))

	query = pgutils.SearchEq(query, "users.id", opts.ID)
	query = pgutils.SearchEq(query, "users.login", opts.Login)

	sql, args, err := query.ToSql()

	logrus.Debugf("\n [get-users: INCOMING] \n OPTS: %+v \n", opts)
	logrus.Debugf("\n [get-users: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return nil, err
	}

	rows, err := s.client.Query(ctx, sql, args...)
	defer rows.Close()

	for rows.Next() {
		u := models.User{}

		scanValues := []any{
			&u.ID,
			&u.Login,
		}

		if err = rows.Scan(scanValues...); err != nil {
			return nil, err
		}

		userEntity := u.ToEntity()

		res = append(res, userEntity)
	}

	return
}

func (s *UserStorage) GetPassword(ctx context.Context, userID string) (res string, err error) {
	defer utils.AddFuncLabel("[repo-get-user-password]", err)

	sql, args, err := s.queryBuilder.
		Select(fmt.Sprintf("pgp_sym_decrypt(password::bytea, '%s')", s.cryptoKey)).
		From(fmt.Sprintf("%s.%s", s.scheme, TABLE_USERS)).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	logrus.Debugf("\n [get-password: INCOMING] \n user-id: %s \n", userID)
	logrus.Debugf("\n [get-password: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return res, err
	}

	err = s.client.QueryRow(ctx, sql, args...).Scan(&res)

	return
}

func (s *UserStorage) GetRefreshToken(ctx context.Context, userID string) (res string, err error) {
	defer utils.AddFuncLabel("[repo-get-user-refresh-token]", err)

	sql, args, err := s.queryBuilder.
		Select("refresh_token").
		From(fmt.Sprintf("%s.%s", s.scheme, TABLE_USERS)).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	logrus.Debugf("\n [get-refresh-token: INCOMING] \n user-id: %s \n", userID)
	logrus.Debugf("\n [get-refresh-token: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return res, err
	}

	err = s.client.QueryRow(ctx, sql, args...).Scan(&res)

	return
}

func (s *UserStorage) Create(ctx context.Context, userCreate entity.UserCreate) (res string, err error) {
	defer utils.AddFuncLabel("[repo-user-create]", err)

	sql, args, err := squirrel.Expr(
		fmt.Sprintf("INSERT INTO todo.users (login,password) VALUES ('%s', pgp_sym_encrypt('%s', '%s')) RETURNING id", userCreate.Login, userCreate.Password, s.cryptoKey),
	).
		ToSql()

	logrus.Debugf("\n [create-user: INCOMING] \n user-create: %+v \n", userCreate)
	logrus.Debugf("\n [create-user: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return res, err
	}

	err = s.client.QueryRow(ctx, sql, args...).Scan(&res)

	return
}

func (s *UserStorage) Update(ctx context.Context, userUpdate entity.UserUpdate) (err error) {
	defer utils.AddFuncLabel("[repo-user-update]", err)

	query := s.queryBuilder.
		Update(fmt.Sprintf("%s.%s", s.scheme, TABLE_USERS)).
		Where(squirrel.Eq{"id": userUpdate.ID})

	if userUpdate.Login != nil {
		query = query.Set("login", *userUpdate.Login)
	}

	if userUpdate.RefreshToken != nil {
		query = query.Set("refresh_token", *userUpdate.RefreshToken)
	}

	sql, args, err := query.ToSql()

	logrus.Debugf("\n [update-user: INCOMING] \n user-update: %+v \n", userUpdate)
	logrus.Debugf("\n [update-user: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	_, err = s.client.Exec(ctx, sql, args...)

	return
}

func (s *UserStorage) UpdatePassword(ctx context.Context, userID, newPassword string) (err error) {
	defer utils.AddFuncLabel("[repo-user-update-password]", err)

	sql, args, err := squirrel.Expr(
		fmt.Sprintf("UPDATE todo.users SET password = pgp_sym_encrypt('%s', '%s') WHERE id = '%s'", newPassword, s.cryptoKey, userID),
	).
		ToSql()

	logrus.Debugf("\n [update-password: INCOMING] \n user-id: %s \n new-password: %s \n", userID, newPassword)
	logrus.Debugf("\n [update-password: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	_, err = s.client.Exec(ctx, sql, args...)

	return
}

func (s *UserStorage) Delete(ctx context.Context, userID string) (err error) {
	defer utils.AddFuncLabel("[repo-user-delete]", err)

	sql, args, err := s.queryBuilder.
		Delete(fmt.Sprintf("%s.%s", s.scheme, TABLE_USERS)).
		Where(squirrel.Eq{"id": userID}).
		ToSql()

	logrus.Debugf("\n [delete-user: INCOMING] \n user-id: %s \n", userID)
	logrus.Debugf("\n [delete-user: SQL] \n Query: %s \n ARGS: %v \n", sql, args)

	if err != nil {
		return
	}

	_, err = s.client.Exec(ctx, sql, args...)

	return
}
