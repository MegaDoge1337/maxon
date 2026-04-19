package user

import (
	"context"
	"fmt"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
	"megadoge1337/maxon/pkg/hasher"
)

type UpdateUserUseCase struct {
	repo repository.UserRepository
}

func NewUpdateUserUseCase(r repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		repo: r,
	}
}

func (uu UpdateUserUseCase) Execute(ctx context.Context, updateUserCommand domain.UpdateUserCommand, id int) (*domain.User, error) {
	// get updated user
	updateUser, err := uu.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	// search via username
	existingUser, _ := uu.repo.GetByUsername(ctx, updateUserCommand.Username)
	// if finded enties are same
	if existingUser != nil && existingUser.ID != id {
		return nil, fmt.Errorf("user with same username already exists")
	}

	// validate credentials
	if updateUserCommand.OldPassword == "" {
		return nil, fmt.Errorf("user credentials are wrong")
	}

	isPasswordValid, err := hasher.CompareHashAndPassword(updateUser.Password, updateUserCommand.OldPassword)
	if err != nil {
		return nil, err
	}

	if !isPasswordValid {
		return nil, fmt.Errorf("user credentials are wrong")
	}

	// update username
	if updateUserCommand.Username != "" {
		updateUser.Username = updateUserCommand.Username
	}

	// update email
	if updateUserCommand.Email != "" {
		updateUser.Email = updateUserCommand.Email
	}

	// update password via hashing
	if updateUserCommand.NewPassword != "" {
		updateUser.Password, err = hasher.HashPassword(updateUserCommand.NewPassword)
		if err != nil {
			return nil, err
		}
	}

	// make update
	updateUser, err = uu.repo.Update(ctx, *updateUser)
	if err != nil {
		return nil, err
	}

	return updateUser, nil
}
