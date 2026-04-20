package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
	"megadoge1337/maxon/pkg/jwt"
	"time"
)

type RefreshUseCaseConfig struct {
	JwtSecret string
}

type RefreshUseCaseDeps struct {
	Repo   repository.AuthRepository
	Config RefreshUseCaseConfig
}

type RefreshUseCase struct {
	repo   repository.AuthRepository
	config RefreshUseCaseConfig
}

func NewRefreshUseCase(deps RefreshUseCaseDeps) *RefreshUseCase {
	return &RefreshUseCase{
		repo:   deps.Repo,
		config: deps.Config,
	}
}

func (r *RefreshUseCase) Execute(ctx context.Context, refreshCommand domain.RefreshCommand) (acess string, refresh string, err error) {
	session, err := r.repo.GetById(ctx, refreshCommand.Refresh)
	if err != nil {
		return "", "", err
	}

	userId := session.UserID

	_, err = r.repo.DeleteById(ctx, *session)
	if err != nil {
		return "", "", err
	}

	access, err := jwt.GenerateAccessToken(userId, r.config.JwtSecret, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	b := make([]byte, 16)
	rand.Read(b)
	refresh = hex.EncodeToString(b)

	refreshSession := domain.Session{
		ID:        refresh,
		UserID:    userId,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	_, err = r.repo.Create(ctx, refreshSession)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
