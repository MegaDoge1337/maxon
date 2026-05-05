package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/internal/repository"
	"megadoge1337/maxon/pkg/hasher"
	"megadoge1337/maxon/pkg/jwt"
	"time"
)

type LoginUseCaseConfig struct {
	JwtSecret string
}

type LoginUseCaseDeps struct {
	AuthRepo repository.AuthRepository
	UserRepo repository.UserRepository
	RoleRepo repository.RoleRepository
	Config   LoginUseCaseConfig
}

type LoginUseCase struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	config   LoginUseCaseConfig
}

func NewLoginUseCase(deps LoginUseCaseDeps) *LoginUseCase {
	return &LoginUseCase{
		authRepo: deps.AuthRepo,
		userRepo: deps.UserRepo,
		roleRepo: deps.RoleRepo,
		config:   deps.Config,
	}
}

func (l *LoginUseCase) Execute(ctx context.Context, loginCommand domain.LoginCommand) (acess string, refresh string, err error) {
	user, err := l.userRepo.GetByUsername(ctx, loginCommand.Username)
	if err != nil {
		return "", "", err
	}

	isPasswordValid, err := hasher.CompareHashAndPassword(user.Password, loginCommand.Password)
	if err != nil {
		return "", "", err
	}

	if !isPasswordValid {
		return "", "", fmt.Errorf("user credentials are wrong")
	}

	role, err := l.roleRepo.GetByUserId(ctx, user.ID)
	if err != nil {
		return "", "", err
	}

	access, err := jwt.GenerateAccessToken(user.ID, role.Name, l.config.JwtSecret, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	b := make([]byte, 16)
	rand.Read(b)
	refresh = hex.EncodeToString(b)

	newSession := domain.Session{
		ID:        refresh,
		UserID:    user.ID,
		Role:      role.Name,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	_, err = l.authRepo.Create(ctx, newSession)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
