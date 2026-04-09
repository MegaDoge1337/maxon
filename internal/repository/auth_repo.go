package repository

import (
	"context"
	"database/sql"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/megadoge1337/maxon/models"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Create(ctx context.Context, s *models.Session) (*models.Session, error) {
	err := s.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *AuthRepository) GetById(ctx context.Context, id string) (*models.Session, error) {
	s, err := models.FindSession(ctx, r.db, id)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (r *AuthRepository) DeleteById(ctx context.Context, s *models.Session) (*models.Session, error) {
	_, err := s.Delete(ctx, r.db)
	if err != nil {
		return nil, err
	}

	return s, nil
}
