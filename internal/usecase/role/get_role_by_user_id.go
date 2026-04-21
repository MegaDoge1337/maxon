package role

import (
	"context"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
)

type GetRoleByUserIdUseCase struct {
	repo repository.RoleRepository
}

func NewGetRoleByUserIdUseCase(r repository.RoleRepository) *GetRoleByUserIdUseCase {
	return &GetRoleByUserIdUseCase{
		repo: r,
	}
}

func (uc *GetRoleByUserIdUseCase) Execute(ctx context.Context, userId int) (*domain.Role, error) {
	role, err := uc.repo.GetByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	return role, nil
}
