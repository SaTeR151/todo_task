package type_service

import (
	"context"

	"github.com/sater-151/todo-list/internal/entity"
	"github.com/sater-151/todo-list/internal/repository/postgres"
	"github.com/sater-151/todo-list/pkg/utils"
)

type TypeService struct {
	repo *postgres.Repository
}

func (s *TypeService) Get(ctx context.Context, opts entity.GetTypesOpts) (types entity.Types, err error) {
	types, err = s.repo.Type.Get(ctx, opts)
	if err != nil {
		return
	}

	if len(types) == 0 {
		return nil, entity.ErrNotFound
	}

	return
}

func (s *TypeService) GetByUserID(ctx context.Context, userID string) (res entity.Types, err error) {
	users, err := s.Get(ctx, entity.GetTypesOpts{UserID: userID})
	if err != nil {
		return
	}

	return users, nil
}

func (s *TypeService) GetByID(ctx context.Context, typeID string) (res entity.Type, err error) {
	types, err := s.repo.Type.Get(ctx, entity.GetTypesOpts{
		ID: typeID,
	})

	if err != nil {
		return
	}

	return types[0], err
}

func (s *TypeService) Create(ctx context.Context, typeCreate entity.TypeCreate) (res entity.Type, err error) {
	defer utils.AddFuncLabel("[service-create-type]", err)

	newUserID, err := s.repo.Type.Create(ctx, typeCreate)
	if err != nil {
		return
	}

	return s.GetByID(ctx, newUserID)
}

func (s *TypeService) Update(ctx context.Context, userID string, typeUpdate entity.TypeUpdate) (res entity.Type, err error) {
	defer utils.AddFuncLabel("[service-update-type]", err)

	_, err = s.Get(ctx, entity.GetTypesOpts{
		UserID: userID,
		ID:     typeUpdate.ID,
	})
	if err != nil {
		return
	}

	if err = s.repo.Type.Update(ctx, typeUpdate); err != nil {
		return
	}

	return s.GetByID(ctx, typeUpdate.ID)
}

func (s *TypeService) Delete(ctx context.Context, typeID string) (err error) {
	defer utils.AddFuncLabel("[service-delete-type]", err)

	return s.repo.Type.Delete(ctx, typeID)
}
