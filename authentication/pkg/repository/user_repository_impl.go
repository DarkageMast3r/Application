package repository

import (
	"authentication/pkg/database"
	"authentication/pkg/models"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GormUserRepository implements domain.UserRepository using GORM
type GormUserRepository struct {
	db database.Database
}

func NewGormUserRepository(db database.Database) UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) FindByUsernameOrEmail(ctx context.Context, username, email string) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("username = ? OR email = ?", username, email).First(&user)
	if result.Error() != nil {
		if errors.Is(result.Error(), gorm.ErrRecordNotFound) {
			return nil, nil // User not found
		}
		return nil, result.Error()
	}
	return &user, nil
}

func (r *GormUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if result.Error() != nil {
		if errors.Is(result.Error(), gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error()
	}
	return &user, nil
}

func (r *GormUserRepository) Save(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Save(user)
	return result.Error()
}

func (r *GormUserRepository) UpdateLastLogin(ctx context.Context, user *models.User) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(user).Update("last_login", now)
	return result.Error()
}

func (r *GormUserRepository) LoadUserRoles(ctx context.Context, user *models.User) error {

	gormDB := r.db.(*database.GormDatabase).DB.WithContext(ctx) // Type assertion naar de onderliggende GORM DB
	err := gormDB.Model(user).Association("Roles").Find(&user.Roles)
	return err
}

func (r *GormUserRepository) LoadRolePermissions(ctx context.Context, role *models.Role) error {
	gormDB := r.db.(*database.GormDatabase).DB.WithContext(ctx) // Type assertion
	err := gormDB.Model(role).Association("Permissions").Find(&role.Permissions)
	return err
}
