package auth

import (
	"context"
	"megadoge1337/maxon/internal/domain"
)

type Register interface {
	Execute(ctx context.Context, registerCommand domain.RegisterCommand) (*domain.User, error)
}

type Login interface {
	Execute(ctx context.Context, loginCommand domain.LoginCommand) (acess string, refresh string, err error)
}

type Refresh interface {
	Execute(ctx context.Context, refreshCommand domain.RefreshCommand) (acess string, refresh string, err error)
}
