package repository

import (
	"context"

	"deu/internal/models"

	"gorm.io/gorm"
)

type PostgresUserPlaceRepository struct {
	DB *gorm.DB
}

func NewPostgresUserPlaceRepository(db *gorm.DB) *PostgresUserPlaceRepository {
	return &PostgresUserPlaceRepository{DB: db}
}

func (r *PostgresUserPlaceRepository) AddVisitedPlace(ctx context.Context, userID, placeID string) error {
	userPlace := models.UserPlace{
		UserID: userID,
		PlaceID: placeID,
	}
	result := r.DB.WithContext(ctx).FirstOrCreate(&userPlace, userPlace)
	
	return result.Error
}

func (r *PostgresUserPlaceRepository) HasVisitedPlace(ctx context.Context, userID, placeID string) (bool, error) {
	var userPlace models.UserPlace
	
	err := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("place_id = ?", placeID).
		First(&userPlace).Error

	if err == nil {
		return true, nil
	}

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return false, err
}

func (r *PostgresUserPlaceRepository) RemoveVisitedPlace(ctx context.Context, userID, placeID string) error {
	result := r.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		Where("place_id = ?", placeID).
		Delete(&models.UserPlace{})
	
	return result.Error
}