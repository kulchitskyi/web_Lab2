package repository

import (
	"deu/internal/models"
	"context"
)

type PlaceRepository interface {
    GetAll(ctx context.Context) ([]models.Place, error)
    GetByID(ctx context.Context, id string) (*models.Place, error)
    Create(ctx context.Context, p *models.Place) error
    Update(ctx context.Context, id string, p *models.PlaceUpdateRequest) error
    Delete(ctx context.Context, id string) error
    DeleteAll(ctx context.Context) error
}