package user

import (
	"context"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
)

type GetUserByIdUseCase struct {
	repo repository.UserRepository
}

func NewGetUserByIdUseCase(r repository.UserRepository) *GetUserByIdUseCase {
	return &GetUserByIdUseCase{
		repo: r,
	}
}

func (uc *GetUserByIdUseCase) Execute(ctx context.Context, id int) (*domain.User, error) {
	user, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
