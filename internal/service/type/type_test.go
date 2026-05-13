package type_service

import (
	"context"
	"errors"
	"testing"

	"github.com/sater-151/todo-list/internal/entity"
)

type fakeTypeRepo struct {
	types   entity.Types
	created []entity.TypeCreate
	updated []entity.TypeUpdate
	deleted []string
	err     error
}

func (r *fakeTypeRepo) Get(_ context.Context, opts entity.GetTypesOpts) (entity.Types, error) {
	if r.err != nil {
		return nil, r.err
	}
	var res entity.Types
	for _, typ := range r.types {
		if opts.ID != "" && typ.ID != opts.ID {
			continue
		}
		if opts.UserID != "" && typ.UserID != opts.UserID {
			continue
		}
		if opts.Name != "" && typ.Name != opts.Name {
			continue
		}
		res = append(res, typ)
	}
	return res, nil
}

func (r *fakeTypeRepo) Create(_ context.Context, typeCreate entity.TypeCreate) (string, error) {
	if r.err != nil {
		return "", r.err
	}
	id := "type-new"
	r.created = append(r.created, typeCreate)
	r.types = append(r.types, entity.Type{ID: id, UserID: typeCreate.UserID, Name: typeCreate.Name, Color: typeCreate.Color})
	return id, nil
}

func (r *fakeTypeRepo) Update(_ context.Context, typeUpdate entity.TypeUpdate) error {
	if r.err != nil {
		return r.err
	}
	r.updated = append(r.updated, typeUpdate)
	for i := range r.types {
		if r.types[i].ID != typeUpdate.ID {
			continue
		}
		if typeUpdate.Name != nil {
			r.types[i].Name = *typeUpdate.Name
		}
		if typeUpdate.Color != nil {
			r.types[i].Color = *typeUpdate.Color
		}
	}
	return nil
}

func (r *fakeTypeRepo) Delete(_ context.Context, typeID string) error {
	if r.err != nil {
		return r.err
	}
	r.deleted = append(r.deleted, typeID)
	return nil
}

func TestTypeServiceGetReturnsNotFoundWhenEmpty(t *testing.T) {
	service := New(&fakeTypeRepo{})

	_, err := service.Get(context.Background(), entity.GetTypesOpts{UserID: "missing"})
	if !errors.Is(err, entity.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTypeServiceCreateUpdateDelete(t *testing.T) {
	repo := &fakeTypeRepo{}
	service := New(repo)

	created, err := service.Create(context.Background(), entity.TypeCreate{UserID: "user-1", Name: "bug", Color: "#ff0000"})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.ID != "type-new" || created.Name != "bug" {
		t.Fatalf("unexpected created type: %+v", created)
	}

	name := "feat"
	color := "#00ff00"
	updated, err := service.Update(context.Background(), "user-1", entity.TypeUpdate{ID: "type-new", Name: &name, Color: &color})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Name != name || updated.Color != color {
		t.Fatalf("unexpected updated type: %+v", updated)
	}

	if err := service.Delete(context.Background(), "type-new"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if len(repo.deleted) != 1 || repo.deleted[0] != "type-new" {
		t.Fatalf("unexpected deletes: %+v", repo.deleted)
	}
}
