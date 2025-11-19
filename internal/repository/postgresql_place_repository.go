package repository

import (
	"context"
	"time"

	er "deu/internal/errors"
	"deu/internal/models"

	"gorm.io/gorm"
)

type PostgresPlaceRepository struct {
	DB *gorm.DB
}

func NewPostgresPlaceRepository(db *gorm.DB) *PostgresPlaceRepository {
	return &PostgresPlaceRepository{DB: db}
}

func (r *PostgresPlaceRepository) GetAll(ctx context.Context) ([]models.Place, error) {
	var places []models.Place
	if err := r.DB.WithContext(ctx).Find(&places).Error; err != nil {
		return nil, err
	}
	return places, nil
}

func (r *PostgresPlaceRepository) GetByID(ctx context.Context, id string) (*models.Place, error) {
	var place models.Place
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&place).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, er.ErrPlaceNotFound
		}
		return nil, err
	}
	return &place, nil
}

func (r *PostgresPlaceRepository) Create(ctx context.Context, p *models.Place) error {
	return r.DB.WithContext(ctx).Create(p).Error
}

func (r *PostgresPlaceRepository) Update(ctx context.Context, id string, p *models.PlaceUpdateRequest) error {
	updates := map[string]interface{}{}
	if p.Name != nil {
		updates["name"] = *p.Name
	}
	if p.Description != nil {
		updates["description"] = *p.Description
	}
	if p.Location != nil {
		updates["location"] = *p.Location
	}
	if p.Address != nil {
		updates["address"] = *p.Address
	}
	if p.Rating != nil {
		updates["rating"] = *p.Rating
	}

	updates["updated_at"] = time.Now()

	result := r.DB.WithContext(ctx).Model(&models.Place{}).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return er.ErrPlaceNotFound
	}
	return nil
}

func (r *PostgresPlaceRepository) Delete(ctx context.Context, id string) error {
	result := r.DB.WithContext(ctx).Where("id = ?", id).Delete(&models.Place{})
	
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return er.ErrPlaceNotFound
	}
	return nil
}

func (r *PostgresPlaceRepository) DeleteAll(ctx context.Context) error {
	return r.DB.WithContext(ctx).Unscoped().Where("1 = 1").Delete(&models.Place{}).Error
}