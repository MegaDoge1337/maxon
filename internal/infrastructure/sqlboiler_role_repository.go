package infrastructure

import (
	"context"
	"database/sql"
	models "megadoge1337/maxon/gen"
	"megadoge1337/maxon/internal/domain"

	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
)

type SqlBoilerRoleRepository struct {
	db *sql.DB
}

func NewSqlBoilerRoleRepository(db *sql.DB) *SqlBoilerRoleRepository {
	return &SqlBoilerRoleRepository{
		db: db,
	}
}

func roleModelToDomain(m models.Role) domain.Role {
	return domain.Role{
		ID:     m.ID,
		UserId: m.UserID,
		Name:   m.Name,
	}
}

func roleModelsToDomains(ms models.RoleSlice) []domain.Role {
	var ds []domain.Role
	for _, m := range ms {
		ds = append(ds, domain.Role{
			ID:     m.ID,
			UserId: m.UserID,
			Name:   m.Name,
		})
	}
	return ds
}

func roleDomainToModel(d domain.Role) models.Role {
	return models.Role{
		ID:     d.ID,
		UserID: d.UserId,
		Name:   d.Name,
	}
}

func (r *SqlBoilerRoleRepository) Create(ctx context.Context, roleDomain domain.Role) (*domain.Role, error) {
	roleModel := roleDomainToModel(roleDomain)
	err := roleModel.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	roleDomain = roleModelToDomain(roleModel)
	return &roleDomain, nil
}

func (r *SqlBoilerRoleRepository) GetAll(ctx context.Context) ([]domain.Role, error) {
	roleModels, err := models.Roles().All(ctx, r.db)
	if err != nil {
		return nil, err
	}
	roleDomains := roleModelsToDomains(roleModels)
	return roleDomains, nil
}

func (r *SqlBoilerRoleRepository) GetById(ctx context.Context, id int) (*domain.Role, error) {
	roleModel, err := models.FindRole(ctx, r.db, id)
	if err != nil {
		return nil, err
	}
	roleDomain := roleModelToDomain(*roleModel)
	return &roleDomain, nil
}

func (r *SqlBoilerRoleRepository) GetByUserId(ctx context.Context, userId int) (*domain.Role, error) {
	roleModel, err := models.Roles(
		qm.Where("user_id = ?", userId),
	).One(ctx, r.db)
	if err != nil {
		return nil, err
	}
	roleDomain := roleModelToDomain(*roleModel)
	return &roleDomain, nil
}

func (r *SqlBoilerRoleRepository) Update(ctx context.Context, roleDomain domain.Role) (*domain.Role, error) {
	roleModel := roleDomainToModel(roleDomain)
	_, err := roleModel.Update(ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}
	roleDomain = roleModelToDomain(roleModel)
	return &roleDomain, nil
}

func (r *SqlBoilerRoleRepository) Delete(ctx context.Context, roleDomain domain.Role) (*domain.Role, error) {
	roleModel := roleDomainToModel(roleDomain)
	_, err := roleModel.Delete(ctx, r.db)
	if err != nil {
		return nil, err
	}
	roleDomain = roleModelToDomain(roleModel)
	return &roleDomain, nil
}
