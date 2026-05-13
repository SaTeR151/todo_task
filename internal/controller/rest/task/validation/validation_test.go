package validation

import (
	"strings"
	"testing"

	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/entity"
)

func TestValidateTaskCreate(t *testing.T) {
	columns := entity.Columns{{ID: "column-1"}}
	types := entity.Types{{ID: "type-1"}}
	tests := []struct {
		name    string
		req     dto.TaskPOST
		wantErr bool
	}{
		{name: "valid minimum", req: dto.TaskPOST{Label: "task"}},
		{name: "valid full", req: dto.TaskPOST{Label: "task", ColumnID: "column-1", TypeID: "type-1", Description: "desc"}},
		{name: "empty label", req: dto.TaskPOST{}, wantErr: true},
		{name: "long label", req: dto.TaskPOST{Label: strings.Repeat("a", 51)}, wantErr: true},
		{name: "bad column", req: dto.TaskPOST{Label: "task", ColumnID: "bad"}, wantErr: true},
		{name: "long description", req: dto.TaskPOST{Label: "task", Description: strings.Repeat("a", 301)}, wantErr: true},
		{name: "bad type", req: dto.TaskPOST{Label: "task", TypeID: "bad"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTaskCreate(tt.req, columns, types)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateTaskCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTaskUpdate(t *testing.T) {
	validLabel := "task"
	emptyLabel := ""
	longLabel := strings.Repeat("a", 51)
	longDescription := strings.Repeat("a", 301)
	validType := "type-1"
	badType := "bad"

	tests := []struct {
		name    string
		req     dto.TaskPATCH
		wantErr bool
	}{
		{name: "valid nil"},
		{name: "valid label", req: dto.TaskPATCH{Label: &validLabel}},
		{name: "valid type", req: dto.TaskPATCH{TypeID: &validType}},
		{name: "empty label", req: dto.TaskPATCH{Label: &emptyLabel}, wantErr: true},
		{name: "long label", req: dto.TaskPATCH{Label: &longLabel}, wantErr: true},
		{name: "long description", req: dto.TaskPATCH{Description: &longDescription}, wantErr: true},
		{name: "bad type", req: dto.TaskPATCH{TypeID: &badType}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTaskUpdate(tt.req, []string{"type-1"})
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateTaskUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
