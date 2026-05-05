package auth

import (
	"context"
	"fmt"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
	"megadoge1337/maxon/pkg/hasher"
)

type RegisterUseCaseDeps struct {
	AuthRepo repository.AuthRepository
	UserRepo repository.UserRepository
	RoleRepo repository.RoleRepository
}

type RegisterUseCase struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewRegisterUseCase(deps RegisterUseCaseDeps) *RegisterUseCase {
	return &RegisterUseCase{
		authRepo: deps.AuthRepo,
		userRepo: deps.UserRepo,
		roleRepo: deps.RoleRepo,
	}
}

func (r *RegisterUseCase) Execute(ctx context.Context, registerCommand domain.RegisterCommand) (*domain.User, error) {
	if registerCommand.Username == "" {
		return nil, fmt.Errorf("username cannot be empty")
	}

	if registerCommand.Password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	if registerCommand.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	existingUser, _ := r.userRepo.GetByUsername(ctx, registerCommand.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("user with same username already exists")
	}

	hashedPassword, err := hasher.HashPassword(registerCommand.Password)
	if err != nil {
		return nil, err
	}

	userDomain := domain.User{
		Username: registerCommand.Username,
		Email:    registerCommand.Email,
		Password: hashedPassword,
	}

	newUser, err := r.userRepo.Create(ctx, userDomain)
	if err != nil {
		return nil, err
	}

	roleDomain := domain.Role{
		UserId: newUser.ID,
		Name:   domain.RolesRegistry["user"],
	}

	_, err = r.roleRepo.Create(ctx, roleDomain)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
