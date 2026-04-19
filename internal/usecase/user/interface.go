package user

import (
	"context"
	"megadoge1337/maxon/internal/domain"
)

type UserCreator interface {
	Execute(ctx context.Context, createUserCommand domain.CreateUserCommand) (*domain.User, error)
}

type AllUsersGetter interface {
	Execute(ctx context.Context) ([]domain.User, error)
}

type UserByIdGetter interface {
	Execute(ctx context.Context, id int) (*domain.User, error)
}

type UserByUsernameGetter interface {
	Execute(ctx context.Context, username string) (*domain.User, error)
}

type UserUpdater interface {
	Execute(ctx context.Context, updateUserCommand domain.UpdateUserCommand, id int) (*domain.User, error)
}

type UserDeleter interface {
	Execute(ctx context.Context, id int) (*domain.User, error)
}
