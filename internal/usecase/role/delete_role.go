package role

import (
	"context"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
)

type DeleteRoleUseCase struct {
	repo repository.RoleRepository
}

func NewDeleteRoleUseCase(r repository.RoleRepository) *DeleteRoleUseCase {
	return &DeleteRoleUseCase{
		repo: r,
	}
}

func (uc *DeleteRoleUseCase) Execute(ctx context.Context, roleDomain domain.Role) (*domain.Role, error) {
	deleteRole, err := uc.repo.Delete(ctx, roleDomain)
	if err != nil {
		return nil, err
	}

	return deleteRole, nil
}
