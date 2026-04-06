package service

import (
	"context"
	"time"

	"github.com/aarondl/null/v8"
	"github.com/megadoge1337/maxon/internal/domain"
	"github.com/megadoge1337/maxon/internal/repository"
	"github.com/megadoge1337/maxon/models"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserRepository(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Create(u domain.User) (*domain.User, error) {
	var userModel models.User
	userModel.ID = u.ID
	userModel.Username = u.Username
	userModel.Email = u.Email
	userModel.Password = u.Password
	userModel.Created = null.TimeFrom(u.Created)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser, err := s.repo.Create(ctx, &userModel)
	if err != nil {
		return nil, err
	}

	u.ID = newUser.ID
	u.Username = newUser.Username
	u.Email = newUser.Email
	u.Password = newUser.Password
	u.Created = newUser.Created.Time

	return &u, nil
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

func (s *UserService) Update(user domain.User, id int) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userModel := models.User{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Created:  null.TimeFrom(user.Created),
	}

	updateUserModel, err := s.repo.UpdateById(ctx, &userModel)
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
