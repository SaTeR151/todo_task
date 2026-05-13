package validation

import (
	"testing"

	"github.com/sater-151/todo-list/internal/controller/rest/dto"
	"github.com/sater-151/todo-list/internal/entity"
)

func TestValidateColumnCreate(t *testing.T) {
	columns := entity.Columns{{ID: "todo", Name: "todo", OrderNumber: 1}}
	tests := []struct {
		name    string
		req     dto.ColumnPOST
		wantErr bool
	}{
		{name: "valid", req: dto.ColumnPOST{Name: "done", OrderNumber: 2}},
		{name: "backlog order", req: dto.ColumnPOST{Name: "backlog", OrderNumber: -1}},
		{name: "empty name", req: dto.ColumnPOST{OrderNumber: 2}, wantErr: true},
		{name: "duplicate name", req: dto.ColumnPOST{Name: "todo", OrderNumber: 2}, wantErr: true},
		{name: "bad order", req: dto.ColumnPOST{Name: "done"}, wantErr: true},
		{name: "duplicate order", req: dto.ColumnPOST{Name: "done", OrderNumber: 1}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColumnCreate(columns, tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateColumnCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateColumnUpdate(t *testing.T) {
	columns := entity.Columns{{ID: "todo", Name: "todo", OrderNumber: 1}}
	validName := "done"
	emptyName := ""
	validOrder := 2
	badOrder := 0
	duplicateOrder := 1

	tests := []struct {
		name    string
		req     dto.ColumnPATCH
		wantErr bool
	}{
		{name: "valid nil"},
		{name: "valid name", req: dto.ColumnPATCH{Name: &validName}},
		{name: "valid order", req: dto.ColumnPATCH{OrderNumber: &validOrder}},
		{name: "empty name", req: dto.ColumnPATCH{Name: &emptyName}, wantErr: true},
		{name: "duplicate name", req: dto.ColumnPATCH{Name: ptr("todo")}, wantErr: true},
		{name: "bad order", req: dto.ColumnPATCH{OrderNumber: &badOrder}, wantErr: true},
		{name: "duplicate order", req: dto.ColumnPATCH{OrderNumber: &duplicateOrder}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateColumnUpdate(columns, tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateColumnUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateColumnSwap(t *testing.T) {
	if err := ValidateColumnSwap(entity.Column{ID: "a"}, entity.Column{ID: "b"}); err != nil {
		t.Fatalf("expected valid swap, got %v", err)
	}
	if err := ValidateColumnSwap(entity.Column{ID: "a"}, entity.Column{ID: "a"}); err == nil {
		t.Fatal("expected error for same column")
	}
}

func ptr[T any](value T) *T {
	return &value
}
