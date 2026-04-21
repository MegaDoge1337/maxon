package role

import (
	"context"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
)

type GetRoleByIdUseCase struct {
	repo repository.RoleRepository
}

func NewGetRoleByIdUseCase(r repository.RoleRepository) *GetRoleByIdUseCase {
	return &GetRoleByIdUseCase{
		repo: r,
	}
}

func (uc *GetRoleByIdUseCase) Execute(ctx context.Context, id int) (*domain.Role, error) {
	role, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return role, nil
}
