package infrastructure

import (
	"context"
	"database/sql"

	"megadoge1337/maxon/internal/domain"
	"megadoge1337/maxon/models"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type SqlBoilerUserRepository struct {
	db *sql.DB
}

func NewSqlBoilerUserRepository(db *sql.DB) *SqlBoilerUserRepository {
	return &SqlBoilerUserRepository{
		db: db,
	}
}

func userModelToDomain(m models.User) domain.User {
	return domain.User{
		ID:       m.ID,
		Username: m.Username,
		Email:    m.Email,
		Password: m.Password,
		Created:  m.Created.Time,
	}
}

func userModelsToDomains(ms models.UserSlice) []domain.User {
	var ds []domain.User
	for _, m := range ms {
		ds = append(ds, domain.User{
			ID:       m.ID,
			Username: m.Username,
			Email:    m.Email,
			Password: m.Password,
			Created:  m.Created.Time,
		})
	}
	return ds
}

func userDomainToModel(d domain.User) models.User {
	return models.User{
		ID:       d.ID,
		Username: d.Username,
		Email:    d.Email,
		Password: d.Password,
		Created:  null.TimeFrom(d.Created),
	}
}

func (r *SqlBoilerUserRepository) Create(ctx context.Context, userDomain domain.User) (*domain.User, error) {
	userModel := userDomainToModel(userDomain)
	err := userModel.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	userDomain = userModelToDomain(userModel)
	return &userDomain, nil
}

func (r *SqlBoilerUserRepository) GetAll(ctx context.Context) ([]domain.User, error) {
	userModels, err := models.Users().All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return userModelsToDomains(userModels), nil
}

func (r *SqlBoilerUserRepository) GetById(ctx context.Context, id int) (*domain.User, error) {
	userModel, err := models.FindUser(ctx, r.db, id)
	if err != nil {
		return nil, err
	}
	userDomain := userModelToDomain(*userModel)
	return &userDomain, nil
}

func (r *SqlBoilerUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	userModel, err := models.Users(
		qm.Where("username = ?", username),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	userDomain := userModelToDomain(*userModel)
	return &userDomain, nil
}

func (r *SqlBoilerUserRepository) UpdateById(ctx context.Context, userDomain domain.User) (*domain.User, error) {
	userModel := userDomainToModel(userDomain)
	_, err := userModel.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	userDomain = userModelToDomain(userModel)
	return &userDomain, nil
}

func (r *SqlBoilerUserRepository) DeleteById(ctx context.Context, userDomain domain.User) (*domain.User, error) {
	userModel := userDomainToModel(userDomain)
	_, err := userModel.Delete(ctx, r.db)
	if err != nil {
		return nil, err
	}
	userDomain = userModelToDomain(userModel)
	return &userDomain, nil
}
