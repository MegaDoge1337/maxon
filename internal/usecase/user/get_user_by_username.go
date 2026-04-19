package user

import (
	"context"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
)

type GetUserByUsernameUseCase struct {
	repo repository.UserRepository
}

func NewGetUserByUsernameUseCase(r repository.UserRepository) *GetUserByUsernameUseCase {
	return &GetUserByUsernameUseCase{
		repo: r,
	}
}

func (ubug *GetUserByUsernameUseCase) Execute(ctx context.Context, username string) (*domain.User, error) {
	user, err := ubug.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
