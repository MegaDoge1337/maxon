package infrastructure

import (
	"context"
	"database/sql"

	models "megadoge1337/maxon/gen"
	"megadoge1337/maxon/internal/domain"

	"github.com/aarondl/sqlboiler/v4/boil"
)

type SqlBoilerAuthRepository struct {
	db *sql.DB
}

func NewSqlBoilerAuthRepository(db *sql.DB) *SqlBoilerAuthRepository {
	return &SqlBoilerAuthRepository{
		db: db,
	}
}

func sessionModelToDomain(m models.Session) domain.Session {
	return domain.Session{
		ID:        m.ID,
		UserID:    m.UserID,
		ExpiresAt: m.ExpiresAt,
		CreatedAt: m.CreatedAt,
	}
}

func sessionDomainToModel(d domain.Session) models.Session {
	return models.Session{
		ID:        d.ID,
		UserID:    d.UserID,
		ExpiresAt: d.ExpiresAt,
		CreatedAt: d.CreatedAt,
	}
}

func (r *SqlBoilerAuthRepository) Create(ctx context.Context, seesionDomain domain.Session) (*domain.Session, error) {
	sessionModel := sessionDomainToModel(seesionDomain)
	err := sessionModel.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	seesionDomain = sessionModelToDomain(sessionModel)
	return &seesionDomain, nil
}

func (r *SqlBoilerAuthRepository) GetById(ctx context.Context, id string) (*domain.Session, error) {
	sessionModel, err := models.FindSession(ctx, r.db, id)
	if err != nil {
		return nil, err
	}
	sessionDomain := sessionModelToDomain(*sessionModel)
	return &sessionDomain, nil
}

func (r *SqlBoilerAuthRepository) DeleteById(ctx context.Context, sessionDomain domain.Session) (*domain.Session, error) {
	sessionModel := sessionDomainToModel(sessionDomain)
	_, err := sessionModel.Delete(ctx, r.db)
	if err != nil {
		return nil, err
	}
	sessionDomain = sessionModelToDomain(sessionModel)
	return &sessionDomain, nil
}
