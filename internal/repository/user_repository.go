package repository

import (
	"context"

	"megadoge1337/maxon/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, userDomain domain.User) (*domain.User, error)
	GetAll(ctx context.Context) ([]domain.User, error)
	GetById(ctx context.Context, id int) (*domain.User, error)
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	Update(ctx context.Context, userDomain domain.User) (*domain.User, error)
	Delete(ctx context.Context, userDomain domain.User) (*domain.User, error)
}
