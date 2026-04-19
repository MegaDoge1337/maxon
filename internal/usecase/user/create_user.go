package user

import (
	"context"
	"fmt"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
	"megadoge1337/maxon/pkg/hasher"
)

type CreateUserUseCase struct {
	repo repository.UserRepository
}

func NewCreateUserUseCase(r repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		repo: r,
	}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, createUserCommand domain.CreateUserCommand) (*domain.User, error) {
	// convert command to domain
	userDomain := createUserCommand.ToDomain()

	// verify unique username
	existingUser, _ := uc.repo.GetByUsername(ctx, userDomain.Username)
	if existingUser != nil {
		return existingUser, fmt.Errorf("user with same username already exists")
	}

	// hashing password
	hashedPassword, err := hasher.HashPassword(userDomain.Password)
	if err != nil {
		return nil, err
	}
	userDomain.Password = hashedPassword

	// create new user
	newUser, err := uc.repo.Create(ctx, userDomain)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
