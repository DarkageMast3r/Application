package repository

import (
	"context"
	"errors"

	"authentication/pkg/database"
	"authentication/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormEndpointRepository struct {
	db database.Database
}

func NewGormEndpointRepository(db database.Database) EndpointRepository {
	return &GormEndpointRepository{db: db}
}

func (r *GormEndpointRepository) Create(ctx context.Context, endpoint *models.APIEndpoint) error {
	result := r.db.WithContext(ctx).Create(endpoint)
	return result.Error()
}

func (r *GormEndpointRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.APIEndpoint, error) {
	var endpoint models.APIEndpoint
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&endpoint)
	if result.Error() != nil {
		if errors.Is(result.Error(), gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error()
	}
	return &endpoint, nil
}

func (r *GormEndpointRepository) FindByPathAndMethod(ctx context.Context, path, method string) (*models.APIEndpoint, error) {
	var endpoint models.APIEndpoint
	result := r.db.WithContext(ctx).Where("path = ? AND method = ?", path, method).First(&endpoint)
	if result.Error() != nil {
		if errors.Is(result.Error(), gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error()
	}
	return &endpoint, nil
}

func (r *GormEndpointRepository) FindAll(ctx context.Context) ([]models.APIEndpoint, error) {
	var endpoints []models.APIEndpoint
	result := r.db.WithContext(ctx).Find(&endpoints)
	if result.Error() != nil {
		return nil, result.Error()
	}
	return endpoints, nil
}

func (r *GormEndpointRepository) Save(ctx context.Context, endpoint *models.APIEndpoint) error {
	result := r.db.WithContext(ctx).Save(endpoint)
	return result.Error()
}

func (r *GormEndpointRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.APIEndpoint{}, id)
	return result.Error()
}

func (r *GormEndpointRepository) LoadEndpointRoles(ctx context.Context, endpoint *models.APIEndpoint) error {
	gormDB := r.db.(*database.GormDatabase).DB.WithContext(ctx) // Aanpassing nodig als GormDatabase anders is genoemd
	err := gormDB.Model(endpoint).Association("Roles").Find(&endpoint.Roles)
	return err
}
