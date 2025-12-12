package repository

import (
	"context"
	"time"

	er "deu/internal/errors"
	"deu/internal/models"

	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	DB *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{DB: db}
}

func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.DB.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	if err := r.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, er.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) Create(ctx context.Context, u *models.User) error {
	return r.DB.WithContext(ctx).Create(u).Error
}

func (r *PostgresUserRepository) Update(ctx context.Context, id string, u *models.UserUpdateRequest) error {
	updates := map[string]interface{}{}
	if u.Name != nil {
		updates["name"] = *u.Name 
	}
	if u.Email != nil {
		updates["email"] = *u.Email
	}
	
	updates["updated_at"] = time.Now() 

	result := r.DB.WithContext(ctx).Model(&models.User{}).Where("id = ?", id).Updates(updates)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return er.ErrUserNotFound
	}
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	result := r.DB.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(&models.User{})
	
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return er.ErrUserNotFound
	}
	return nil
}

func (r *PostgresUserRepository) DeleteAll(ctx context.Context) error {
	return r.DB.WithContext(ctx).Unscoped().Where("1 = 1").Delete(&models.User{}).Error
}