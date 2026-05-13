package user

import (
	"context"
	"errors"
	"testing"

	"github.com/sater-151/todo-list/internal/entity"
)

type fakeUserRepo struct {
	users        entity.Users
	passwords    map[string]string
	refreshToken map[string]string
	created      []entity.UserCreate
	updated      []entity.UserUpdate
	deleted      []string
	err          error
}

func (r *fakeUserRepo) Create(_ context.Context, userCreate entity.UserCreate) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	id := "user-new"
	r.created = append(r.created, userCreate)
	r.users = append(r.users, entity.User{ID: id, Login: userCreate.Login})
	if r.passwords == nil {
		r.passwords = map[string]string{}
	}
	r.passwords[id] = userCreate.Password
	return id, nil
}

func (r *fakeUserRepo) Update(_ context.Context, userUpdate entity.UserUpdate) error {
	if r.err != nil {
		return r.err
	}
	r.updated = append(r.updated, userUpdate)
	for i := range r.users {
		if r.users[i].ID != userUpdate.ID {
			continue
		}
		if userUpdate.Login != nil {
			r.users[i].Login = *userUpdate.Login
		}
		if userUpdate.RefreshToken != nil {
			if r.refreshToken == nil {
				r.refreshToken = map[string]string{}
			}
			r.refreshToken[userUpdate.ID] = *userUpdate.RefreshToken
			r.users[i].RefreshToken = *userUpdate.RefreshToken
		}
	}
	return nil
}

func (r *fakeUserRepo) UpdatePassword(_ context.Context, userID, newPassword string) error {
	if r.err != nil {
		return r.err
	}
	if r.passwords == nil {
		r.passwords = map[string]string{}
	}
	r.passwords[userID] = newPassword
	return nil
}

func (r *fakeUserRepo) Delete(_ context.Context, userID string) error {
	if r.err != nil {
		return r.err
	}
	r.deleted = append(r.deleted, userID)
	return nil
}

func (r *fakeUserRepo) Get(_ context.Context, opts entity.GetUsersOpts) (entity.Users, error) {
	if r.err != nil {
		return nil, r.err
	}
	var res entity.Users
	for _, user := range r.users {
		if opts.ID != "" && user.ID != opts.ID {
			continue
		}
		if opts.Login != "" && user.Login != opts.Login {
			continue
		}
		res = append(res, user)
	}
	return res, nil
}

func (r *fakeUserRepo) GetRefreshToken(_ context.Context, userID string) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return r.refreshToken[userID], nil
}

func (r *fakeUserRepo) GetPassword(_ context.Context, userID string) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	return r.passwords[userID], nil
}

type fakeTypeCreator struct {
	created []entity.TypeCreate
	err     error
}

func (r *fakeTypeCreator) Create(_ context.Context, typeCreate entity.TypeCreate) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	r.created = append(r.created, typeCreate)
	return "type-new", nil
}

func TestUserServiceGetReturnsNotFoundWhenEmpty(t *testing.T) {
	service := New(&fakeUserRepo{}, &fakeTypeCreator{}, "secret")

	_, err := service.Get(context.Background(), entity.GetUsersOpts{Login: "missing"})
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUserServiceCreateAlsoCreatesNullType(t *testing.T) {
	users := &fakeUserRepo{}
	types := &fakeTypeCreator{}
	service := New(users, types, "secret")

	created, err := service.Create(context.Background(), entity.UserCreate{Login: "alice", Password: "pwd"})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID != "user-new" || created.Login != "alice" {
		t.Fatalf("unexpected created user: %+v", created)
	}
	if len(types.created) != 1 || types.created[0].UserID != "user-new" || types.created[0].Name != "null" {
		t.Fatalf("unexpected created types: %+v", types.created)
	}
}

func TestUserServiceUpdateDeleteAuthAndRefresh(t *testing.T) {
	users := &fakeUserRepo{
		users:        entity.Users{{ID: "user-1", Login: "alice"}},
		passwords:    map[string]string{"user-1": "pwd"},
		refreshToken: map[string]string{},
	}
	service := New(users, &fakeTypeCreator{}, "secret")

	login := "bob"
	password := "new-pwd"
	updated, err := service.Update(context.Background(), entity.UserUpdate{ID: "user-1", Login: &login, Password: &password})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Login != login || users.passwords["user-1"] != password {
		t.Fatalf("unexpected update result: user=%+v passwords=%+v", updated, users.passwords)
	}

	accessToken, refreshToken, appErr := service.Auth(context.Background(), "bob", "new-pwd")
	if appErr != nil {
		t.Fatalf("Auth returned app error: %v", appErr)
	}
	if accessToken == "" || refreshToken == "" {
		t.Fatalf("expected access and refresh tokens")
	}
	parsedUserID, appErr := service.ParseToken(context.Background(), "Bearer "+accessToken)
	if appErr != nil || parsedUserID != "user-1" {
		t.Fatalf("unexpected ParseToken result: user=%q err=%v", parsedUserID, appErr)
	}

	newAccessToken, appErr := service.RefreshToken(context.Background(), "user-1", refreshToken)
	if appErr != nil {
		t.Fatalf("RefreshToken returned app error: %v", appErr)
	}
	if newAccessToken == "" {
		t.Fatal("expected refreshed access token")
	}

	if err := service.Delete(context.Background(), "user-1"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if len(users.deleted) != 1 || users.deleted[0] != "user-1" {
		t.Fatalf("unexpected deletes: %+v", users.deleted)
	}
}

func TestUserServiceAuthRejectsBadPassword(t *testing.T) {
	users := &fakeUserRepo{
		users:     entity.Users{{ID: "user-1", Login: "alice"}},
		passwords: map[string]string{"user-1": "pwd"},
	}
	service := New(users, &fakeTypeCreator{}, "secret")

	_, _, appErr := service.Auth(context.Background(), "alice", "wrong")
	if appErr == nil || !appErr.IsBadAuth() {
		t.Fatalf("expected bad auth error, got %v", appErr)
	}
}
