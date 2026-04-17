package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/megadoge1337/maxon/internal/domain"
	"github.com/megadoge1337/maxon/internal/repository"
	"github.com/megadoge1337/maxon/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceConfig struct {
	JwtSecret string
}

type AuthSerivceDeps struct {
	AuthRepo repository.AuthRepository
	UserRepo repository.UserRepository
	Config   AuthServiceConfig
}

type AuthService struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
	config   AuthServiceConfig
}

func NewAuthService(deps AuthSerivceDeps) *AuthService {
	return &AuthService{
		authRepo: deps.AuthRepo,
		userRepo: deps.UserRepo,
		config:   deps.Config,
	}
}

func (s *AuthService) Login(createSessionCommand domain.CreateSessionCommand) (acess string, refresh string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.userRepo.GetByUsername(ctx, createSessionCommand.Login)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(createSessionCommand.Password))
	if err != nil {
		return "", "", err
	}

	access, err := jwt.GenerateAccessToken(user.ID, s.config.JwtSecret, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	b := make([]byte, 16)
	rand.Read(b)
	refresh = hex.EncodeToString(b)

	newSession := domain.Session{
		ID:        refresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	_, err = s.authRepo.Create(ctx, newSession)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) Refresh(refreshSessionCommand domain.RefreshSessionCommand) (acess string, refresh string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := s.authRepo.GetById(ctx, refreshSessionCommand.Refresh)
	if err != nil {
		return "", "", err
	}

	userId := session.UserID

	_, err = s.authRepo.DeleteById(ctx, *session)
	if err != nil {
		return "", "", err
	}

	access, err := jwt.GenerateAccessToken(userId, s.config.JwtSecret, 15*time.Minute)
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

	_, err = s.authRepo.Create(ctx, refreshSession)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
