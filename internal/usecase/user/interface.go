package user

import (
	"context"
	"megadoge1337/maxon/internal/domain"
)

type CreateUser interface {
	Execute(ctx context.Context, createUserCommand domain.CreateUserCommand) (*domain.User, error)
}

type GetAllUsers interface {
	Execute(ctx context.Context) ([]domain.User, error)
}

type GetUserById interface {
	Execute(ctx context.Context, id int) (*domain.User, error)
}

type GetUserByUsername interface {
	Execute(ctx context.Context, username string) (*domain.User, error)
}

type UpdateUser interface {
	Execute(ctx context.Context, updateUserCommand domain.UpdateUserCommand, id int) (*domain.User, error)
}

type DeleteUser interface {
	Execute(ctx context.Context, id int) (*domain.User, error)
}
