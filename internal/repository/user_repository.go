package repository

import (
	"deu/internal/models"
	"context"
)

type UserRepository interface {
    GetAll(ctx context.Context) ([]models.User, error)
    GetByID(ctx context.Context, id string) (*models.User, error)
    Create(ctx context.Context, u *models.User) error
    Update(ctx context.Context, id string, u *models.UserUpdateRequest) error
    Delete(ctx context.Context, id string) error
    DeleteAll(ctx context.Context) error
}