package repository

import (
	"context"
	"megadoge1337/maxon/internal/domain"
)

type RoleRepository interface {
	Create(ctx context.Context, roleDomain domain.Role) (*domain.Role, error)
	GetAll(ctx context.Context) ([]domain.Role, error)
	GetById(ctx context.Context, id int) (*domain.Role, error)
	GetByUserId(ctx context.Context, userId int) (*domain.Role, error)
	Update(ctx context.Context, roleDomain domain.Role) (*domain.Role, error)
	Delete(ctx context.Context, roleDomain domain.Role) (*domain.Role, error)
}
