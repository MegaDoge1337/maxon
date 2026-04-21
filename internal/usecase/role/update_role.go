package role

import (
	"context"
	"fmt"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
)

type UpdateRoleUseCaseDeps struct {
	RoleRepo repository.RoleRepository
	UserRepo repository.UserRepository
}

type UpdateRoleUseCase struct {
	roleRepo repository.RoleRepository
	userRepo repository.UserRepository
}

func NewUpdateRoleUseCase(deps UpdateRoleUseCaseDeps) *UpdateRoleUseCase {
	return &UpdateRoleUseCase{
		roleRepo: deps.RoleRepo,
		userRepo: deps.UserRepo,
	}
}

func (uc *UpdateRoleUseCase) Execute(ctx context.Context, roleDomain domain.Role) (*domain.Role, error) {
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

	// update role
	updateRole, err := uc.roleRepo.Update(ctx, roleDomain)
	if err != nil {
		return nil, err
	}

	return updateRole, nil
}
