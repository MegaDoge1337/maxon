package role

import (
	"context"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
)

type GetAllRolesUseCase struct {
	repo repository.RoleRepository
}

func NewGetAllRolesUseCase(r repository.RoleRepository) *GetAllRolesUseCase {
	return &GetAllRolesUseCase{
		repo: r,
	}
}

func (uc *GetAllRolesUseCase) Execute(ctx context.Context) ([]domain.Role, error) {
	roles, err := uc.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
