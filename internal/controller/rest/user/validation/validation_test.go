package validation

import (
	"strings"
	"testing"

	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/entity"
)

func TestValidateUserCreate(t *testing.T) {
	users := entity.Users{{ID: "user-1", Login: "alice"}}
	tests := []struct {
		name    string
		req     entity.UserCreate
		wantErr bool
	}{
		{name: "valid", req: entity.UserCreate{Login: "bob", Password: "pwd"}},
		{name: "empty login", req: entity.UserCreate{Password: "pwd"}, wantErr: true},
		{name: "duplicate", req: entity.UserCreate{Login: "alice", Password: "pwd"}, wantErr: true},
		{name: "bad login", req: entity.UserCreate{Login: "bob!", Password: "pwd"}, wantErr: true},
		{name: "long login", req: entity.UserCreate{Login: strings.Repeat("a", 11), Password: "pwd"}, wantErr: true},
		{name: "empty password", req: entity.UserCreate{Login: "bob"}, wantErr: true},
		{name: "long password", req: entity.UserCreate{Login: "bob", Password: strings.Repeat("a", 31)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserCreate(tt.req, users)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateUserCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateUserPasswordChange(t *testing.T) {
	tests := []struct {
		name    string
		req     dto.UserPasswordChange
		current string
		wantErr bool
	}{
		{name: "valid", req: dto.UserPasswordChange{OldPassword: "old", NewPassword: "new"}, current: "old"},
		{name: "wrong old", req: dto.UserPasswordChange{OldPassword: "bad", NewPassword: "new"}, current: "old", wantErr: true},
		{name: "empty old", req: dto.UserPasswordChange{NewPassword: "new"}, current: "", wantErr: true},
		{name: "empty new", req: dto.UserPasswordChange{OldPassword: "old"}, current: "old", wantErr: true},
		{name: "same", req: dto.UserPasswordChange{OldPassword: "old", NewPassword: "old"}, current: "old", wantErr: true},
		{name: "long new", req: dto.UserPasswordChange{OldPassword: "old", NewPassword: strings.Repeat("a", 31)}, current: "old", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserPasswordChange(tt.req, tt.current)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateUserPasswordChange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
