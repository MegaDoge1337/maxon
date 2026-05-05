package user

import (
	"context"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
)

type GetAllUsersUseCase struct {
	repo repository.UserRepository
}

func NewGetAllUsersUseCase(r repository.UserRepository) *GetAllUsersUseCase {
	return &GetAllUsersUseCase{
		repo: r,
	}
}

func (uc *GetAllUsersUseCase) Execute(ctx context.Context) ([]domain.User, error) {
	return uc.repo.GetAll(ctx)
}
