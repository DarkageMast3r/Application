package repository

import (
	"context"
	"errors"
	"time"

	"authentication/pkg/database"
	"authentication/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormAuthTokenRepository struct {
	db database.Database
}

func NewGormAuthTokenRepository(db database.Database) AuthTokenRepository {
	return &GormAuthTokenRepository{db: db}
}

func (r *GormAuthTokenRepository) Create(ctx context.Context, token *models.RefreshToken) error {
	result := r.db.WithContext(ctx).Create(token)
	return result.Error()
}

func (r *GormAuthTokenRepository) FindByRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	result := r.db.WithContext(ctx).Where("refresh_token = ?", token).First(&refreshToken)
	if result.Error() != nil {
		if errors.Is(result.Error(), gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error()
	}
	return &refreshToken, nil
}

func (r *GormAuthTokenRepository) DeleteByID(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&models.RefreshToken{}, id)
	return result.Error()
}

func (r *GormAuthTokenRepository) DeleteExpiredTokens(ctx context.Context) error {
	result := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{})
	return result.Error()
}

func (r *GormAuthTokenRepository) DeleteUserTokens(ctx context.Context, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.RefreshToken{})
	return result.Error()
}

func (r *GormAuthTokenRepository) Save(ctx context.Context, token *models.RefreshToken) error {
	result := r.db.WithContext(ctx).Save(token)
	return result.Error()
}

func (r *GormAuthTokenRepository) AddBlacklistedAccessToken(ctx context.Context, token string, expiresAt time.Time) error {
	// Niet ideaal voor GORM, zie Redis implementatie
	return errors.New("AddBlacklistedAccessToken not implemented in GormAuthTokenRepository; use CacheRepository instead")
}

func (r *GormAuthTokenRepository) IsAccessTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	// Niet ideaal voor GORM, zie Redis implementatie
	return false, errors.New("IsAccessTokenBlacklisted not implemented in GormAuthTokenRepository; use CacheRepository instead")
}

func (r *GormAuthTokenRepository) CleanExpiredBlacklistedTokens(ctx context.Context) error {
	// Niet ideaal voor GORM, zie Redis implementatie
	return errors.New("CleanExpiredBlacklistedTokens not implemented in GormAuthTokenRepository; use CacheRepository instead")
}
