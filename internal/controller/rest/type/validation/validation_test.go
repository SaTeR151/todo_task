package validation

import (
	"strings"
	"testing"

	"github.com/sater-151/todo-list/internal/entity"
)

func TestValidateTypeCreate(t *testing.T) {
	types := entity.Types{{ID: "type-1", Name: "bug"}}
	tests := []struct {
		name    string
		req     entity.TypeCreate
		wantErr bool
	}{
		{name: "valid", req: entity.TypeCreate{Name: "feat", Color: "#ff0000"}},
		{name: "empty name", req: entity.TypeCreate{Color: "#ff0000"}, wantErr: true},
		{name: "duplicate", req: entity.TypeCreate{Name: "bug", Color: "#ff0000"}, wantErr: true},
		{name: "empty color", req: entity.TypeCreate{Name: "feat"}, wantErr: true},
		{name: "long name", req: entity.TypeCreate{Name: strings.Repeat("a", 11), Color: "#ff0000"}, wantErr: true},
		{name: "bad name", req: entity.TypeCreate{Name: "Не", Color: "#ff0000"}, wantErr: true},
		{name: "bad color", req: entity.TypeCreate{Name: "feat", Color: "red"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTypeCreate(tt.req, types)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateTypeCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateTypeUpdate(t *testing.T) {
	types := entity.Types{{ID: "type-1", Name: "bug"}}
	validName := "feat"
	emptyName := ""
	longName := strings.Repeat("a", 11)
	badName := "Не"
	duplicateName := "bug"
	validColor := "#00ff00"
	emptyColor := ""
	badColor := "green"

	tests := []struct {
		name    string
		req     entity.TypeUpdate
		wantErr bool
	}{
		{name: "valid nil"},
		{name: "valid name", req: entity.TypeUpdate{Name: &validName}},
		{name: "valid color", req: entity.TypeUpdate{Color: &validColor}},
		{name: "empty name", req: entity.TypeUpdate{Name: &emptyName}, wantErr: true},
		{name: "long name", req: entity.TypeUpdate{Name: &longName}, wantErr: true},
		{name: "bad name", req: entity.TypeUpdate{Name: &badName}, wantErr: true},
		{name: "duplicate name", req: entity.TypeUpdate{Name: &duplicateName}, wantErr: true},
		{name: "empty color", req: entity.TypeUpdate{Color: &emptyColor}, wantErr: true},
		{name: "bad color", req: entity.TypeUpdate{Color: &badColor}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTypeUpdate(tt.req, types)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateTypeUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
