package service

import (
	"context"
	"time"

	"github.com/megadoge1337/maxon/internal/domain"
	"github.com/megadoge1337/maxon/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(ctx context.Context, createUserCommand domain.CreateUserCommand) (*domain.User, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(createUserCommand.Password), bcrypt.DefaultCost)
	hashedPassword := string(hashBytes)

	userDomain := createUserCommand.ToDomain()
	userDomain.Password = hashedPassword

	newUser, err := s.repo.Create(ctx, userDomain)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *UserService) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *UserService) GetById(ctx context.Context, id int) (*domain.User, error) {
	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *UserService) UpdateById(updateUserCommand domain.UpdateUserCommand, id int) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if updateUserCommand.Username != "" {
		user.Username = updateUserCommand.Username
	}

	if updateUserCommand.Email != "" {
		user.Email = updateUserCommand.Email
	}

	if updateUserCommand.OldPassword != "" && updateUserCommand.NewPassword != "" {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updateUserCommand.OldPassword))
		if err != nil {
			return nil, err
		}

		hashBytes, err := bcrypt.GenerateFromPassword([]byte(updateUserCommand.NewPassword), bcrypt.DefaultCost)
		hashedPassword := string(hashBytes)

		user.Password = hashedPassword
	}

	updateUser, err := s.repo.UpdateById(ctx, *user)
	if err != nil {
		return nil, err
	}

	return updateUser, err
}

func (s *UserService) DeleteById(id int) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	deleteUser, err := s.repo.DeleteById(ctx, *user)
	if err != nil {
		return nil, err
	}

	return deleteUser, nil
}
