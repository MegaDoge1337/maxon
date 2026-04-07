package service

import (
	"context"
	"time"

	"github.com/megadoge1337/maxon/internal/domain"
	"github.com/megadoge1337/maxon/internal/repository"
	"github.com/megadoge1337/maxon/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(createUserCommand domain.CreateUserCommand) (*domain.User, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(createUserCommand.Password), bcrypt.DefaultCost)
	hashedPassword := string(hashBytes)

	var userModel models.User
	userModel.Username = createUserCommand.Username
	userModel.Email = createUserCommand.Email
	userModel.Password = hashedPassword

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser, err := s.repo.Create(ctx, &userModel)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
		Password: newUser.Password,
		Created:  newUser.Created.Time,
	}

	return &user, nil
}

func (s *UserService) GetAll() ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userModels, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var users []domain.User

	for _, user := range userModels {
		users = append(users, domain.User{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Password: user.Password,
			Created:  user.Created.Time,
		})
	}

	return users, nil
}

func (s *UserService) GetById(id int) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userModel, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Password: userModel.Password,
		Created:  userModel.Created.Time,
	}

	return &user, err
}

func (s *UserService) GetByUsername(username string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userModel, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Password: userModel.Password,
		Created:  userModel.Created.Time,
	}

	return &user, err
}

func (s *UserService) UpdateById(updateUserCommand domain.UpdateUserCommand, id int) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userModel, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	if updateUserCommand.Username != "" {
		userModel.Username = updateUserCommand.Username
	}

	if updateUserCommand.Email != "" {
		userModel.Email = updateUserCommand.Email
	}

	if updateUserCommand.OldPassword != "" && updateUserCommand.NewPassword != "" {
		err := bcrypt.CompareHashAndPassword([]byte(userModel.Password), []byte(updateUserCommand.OldPassword))
		if err != nil {
			return nil, err
		}

		hashBytes, err := bcrypt.GenerateFromPassword([]byte(updateUserCommand.NewPassword), bcrypt.DefaultCost)
		hashedPassword := string(hashBytes)

		userModel.Password = hashedPassword
	}

	updateUserModel, err := s.repo.UpdateById(ctx, userModel)
	if err != nil {
		return nil, err
	}

	updateUser := domain.User{
		ID:       updateUserModel.ID,
		Username: updateUserModel.Username,
		Email:    updateUserModel.Email,
		Password: updateUserModel.Password,
		Created:  updateUserModel.Created.Time,
	}

	return &updateUser, err
}

func (s *UserService) DeleteById(id int) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userModel, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	deleteUserModel, err := s.repo.DeleteById(ctx, userModel)
	if err != nil {
		return nil, err
	}

	deleteUser := domain.User{
		ID:       deleteUserModel.ID,
		Username: deleteUserModel.Username,
		Email:    deleteUserModel.Email,
		Password: deleteUserModel.Password,
		Created:  deleteUserModel.Created.Time,
	}

	return &deleteUser, nil
}
