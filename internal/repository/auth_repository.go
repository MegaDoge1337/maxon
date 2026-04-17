package repository

import (
	"context"

	"github.com/megadoge1337/maxon/internal/domain"
)

type AuthRepository interface {
	Create(ctx context.Context, sessionDomain domain.Session) (*domain.Session, error)
	GetById(ctx context.Context, id string) (*domain.Session, error)
	DeleteById(ctx context.Context, sessionDomain domain.Session) (*domain.Session, error)
}
