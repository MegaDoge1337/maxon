package user

import (
	"context"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
)

type DeleteUserUseCase struct {
	repo repository.UserRepository
}

func NewDeleteUserUseCase(r repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		repo: r,
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, id int) (*domain.User, error) {
	user, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	deleteUser, err := uc.repo.Delete(ctx, *user)
	if err != nil {
		return nil, err
	}

	return deleteUser, nil
}
