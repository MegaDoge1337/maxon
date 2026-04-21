package role

import (
	"context"
	"megadoge1337/maxon/internal/domain"
)

type CreateRole interface {
	Execute(ctx context.Context, roleDomain domain.Role) (*domain.Role, error)
}

type GetAllRoles interface {
	Execute(ctx context.Context) ([]domain.Role, error)
}

type GetRoleById interface {
	Execute(ctx context.Context, id int) (*domain.Role, error)
}

type GetRoleByUserId interface {
	Execute(ctx context.Context, userId int) (*domain.Role, error)
}

type UpdateRole interface {
	Execute(ctx context.Context, roleDomain domain.Role) (*domain.Role, error)
}

type DeleteRole interface {
	Execute(ctx context.Context, roleDomain domain.Role) (*domain.Role, error)
}
