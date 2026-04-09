package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/megadoge1337/maxon/internal/domain"
	"github.com/megadoge1337/maxon/internal/repository"
	"github.com/megadoge1337/maxon/models"
	"github.com/megadoge1337/maxon/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo  *repository.AuthRepository
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(authRepo *repository.AuthRepository, userRepo *repository.UserRepository, secret string) *AuthService {
	return &AuthService{
		authRepo:  authRepo,
		userRepo:  userRepo,
		jwtSecret: secret,
	}
}

func (s *AuthService) Login(createSessionCommand domain.CreateSessionCommand) (string, string, error) {
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

	access, err := jwt.GenerateAccessToken(user.ID, s.jwtSecret, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	b := make([]byte, 16)
	rand.Read(b)
	refresh := hex.EncodeToString(b)

	session := &models.Session{
		ID:        refresh,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
	}

	_, err = s.authRepo.Create(ctx, session)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) Refresh(refreshSessionCommand domain.RefreshSessionCommand) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	session, err := s.authRepo.GetById(ctx, refreshSessionCommand.Refresh)
	if err != nil {
		return "", "", err
	}

	userId := session.UserID

	_, err = s.authRepo.DeleteById(ctx, session)
	if err != nil {
		return "", "", err
	}

	access, err := jwt.GenerateAccessToken(userId, s.jwtSecret, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	b := make([]byte, 16)
	rand.Read(b)
	refresh := hex.EncodeToString(b)

	refreshSession := &models.Session{
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
