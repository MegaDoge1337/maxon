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

func (aug *GetAllUsersUseCase) Execute(ctx context.Context) ([]domain.User, error) {
	return aug.repo.GetAll(ctx)
}
