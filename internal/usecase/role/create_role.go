package role

import (
	"context"
	"fmt"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
)

type CreateRoleUseCaseDeps struct {
	RoleRepo repository.RoleRepository
	UserRepo repository.UserRepository
}

type CreateRoleUseCase struct {
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
}

func NewCreateRoleUseCase(deps CreateRoleUseCaseDeps) *CreateRoleUseCase {
	return &CreateRoleUseCase{
		roleRepo: deps.RoleRepo,
		userRepo: deps.UserRepo,
	}
}

func (uc *CreateRoleUseCase) Execute(ctx context.Context, roleDomain domain.Role) (*domain.Role, error) {
	// validate role name via registry
	roleName, roleExists := domain.RolesRegistry[roleDomain.Name]
	if !roleExists {
		return nil, fmt.Errorf("role does not exists")
	}
	// update role name
	roleDomain.Name = roleName

	// check if user exists
	existingUser, err := uc.userRepo.GetById(ctx, roleDomain.UserId)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, fmt.Errorf("user does not exists")
	}

	// create role
	newRole, err := uc.roleRepo.Create(ctx, roleDomain)
	if err != nil {
		return nil, err
	}
	return newRole, nil
}
