package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/megadoge1337/maxon/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) (*models.User, error) {
	err := u.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetAll(ctx context.Context) (models.UserSlice, error) {
	u, err := models.Users().All(ctx, r.db)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetById(ctx context.Context, id int) (*models.User, error) {
	u, err := models.FindUser(ctx, r.db, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	u, err := models.Users(
		qm.Where("username = ?", username),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) UpdateById(ctx context.Context, u *models.User) (*models.User, error) {
	_, err := u.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) DeleteById(ctx context.Context, u *models.User) (*models.User, error) {
	_, err := u.Delete(ctx, r.db)
	if err != nil {
		return nil, err
	}

	return u, nil
}
