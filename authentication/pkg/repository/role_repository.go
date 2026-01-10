package repository

import (
	"authentication/pkg/database"
	"authentication/pkg/models"
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRoleRepository struct {
	db database.Database
}

func NewGormRoleRepository(db database.Database) RoleRepository {
	return &GormRoleRepository{db: db}
}

func (r *GormRoleRepository) Create(ctx context.Context, role *models.Role) error {
	result := r.db.WithContext(ctx).Create(role)
	return result.Error()
}

func (r *GormRoleRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&role)
	if result.Error() != nil {
		if errors.Is(result.Error(), gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error()
	}
	return &role, nil
}

func (r *GormRoleRepository) FindByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	result := r.db.WithContext(ctx).Where("name = ?", name).First(&role)
	if result.Error() != nil {
		if errors.Is(result.Error(), gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error()
	}
	return &role, nil
}

func (r *GormRoleRepository) FindAll(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	result := r.db.WithContext(ctx).Find(&roles)
	if result.Error() != nil {
		return nil, result.Error()
	}
	return roles, nil
}

func (r *GormRoleRepository) Save(ctx context.Context, role *models.Role) error {
	result := r.db.WithContext(ctx).Save(role)
	return result.Error()
}

func (r *GormRoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.Role{}, id)
	return result.Error()
}

func (r *GormRoleRepository) LoadRolePermissions(ctx context.Context, role *models.Role) error {

	gormDB := r.db.(*database.GormDatabase).DB.WithContext(ctx) // Aanpassing nodig als GormDatabase anders is genoemd
	err := gormDB.Model(role).Association("Permissions").Find(&role.Permissions)
	return err
}

// Optioneel:
// func (r *GormRoleRepository) FindByNames(ctx context.Context, names []string) ([]models.Role, error) {
// 	var roles []models.Role
// 	result := r.db.WithContext(ctx).Where("name IN ?", names).Find(&roles)
// 	if result.Error() != nil {
// 		return nil, result.Error()
// 	}
// 	return roles, nil
// }
