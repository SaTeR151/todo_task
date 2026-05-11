package validation

import (
	"strings"
	"testing"

	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/entity"
)

func TestValidateBoardCreate(t *testing.T) {
	tests := []struct {
		name    string
		req     dto.BoardPOST
		wantErr bool
	}{
		{name: "valid", req: dto.BoardPOST{Name: "board"}},
		{name: "empty", req: dto.BoardPOST{}, wantErr: true},
		{name: "too long", req: dto.BoardPOST{Name: strings.Repeat("a", 51)}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBoardCreate(tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateBoardCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateBoardUpdate(t *testing.T) {
	valid := "new"
	empty := ""
	duplicate := "old"

	tests := []struct {
		name    string
		req     dto.BoardPATCH
		wantErr bool
	}{
		{name: "valid nil"},
		{name: "valid name", req: dto.BoardPATCH{Name: &valid}},
		{name: "empty", req: dto.BoardPATCH{Name: &empty}, wantErr: true},
		{name: "too long", req: dto.BoardPATCH{Name: ptr(strings.Repeat("a", 51))}, wantErr: true},
		{name: "duplicate", req: dto.BoardPATCH{Name: &duplicate}, wantErr: true},
	}

	boards := entity.Boards{{ID: "board-1", Name: "old"}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateBoardUpdate(tt.req, boards)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateBoardUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ptr[T any](value T) *T {
	return &value
}
