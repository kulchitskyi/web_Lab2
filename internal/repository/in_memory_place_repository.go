package repository

import (
	"deu/internal/models"
	"context"
    "sync"

    er "deu/internal/errors"
)

type MemoryPlaceRepository struct {
    mu      sync.RWMutex
    places  map[string]models.Place
}

func NewMemoryPlaceRepository() *MemoryPlaceRepository {
    return &MemoryPlaceRepository{
        places: make(map[string]models.Place),
    }
}

func (r *MemoryPlaceRepository) GetAll(ctx context.Context) ([]models.Place, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    result := make([]models.Place, 0, len(r.places))
    for _, p := range r.places {
        result = append(result, p)
    }
    return result, nil
}

func (r *MemoryPlaceRepository) GetByID(ctx context.Context, id string) (*models.Place, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    result, ok := r.places[id]
    if ok {
        return &result, nil
    }

    return nil, er.ErrPlaceNotFound
}

func (r *MemoryPlaceRepository) Create(ctx context.Context, p *models.Place) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.places[p.Id] = *p
    return nil
}

func (r *MemoryPlaceRepository) Update(ctx context.Context, id string, p *models.PlaceUpdateRequest) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    value, ok := r.places[id]
    if !ok {
        return er.ErrPlaceNotFound
    }

    if p.Name != nil {
        value.Name = *p.Name
    }
    if p.Description != nil {
        value.Description = *p.Description
    }
	if p.Location != nil {
        value.Location = *p.Location
    }
	if p.Address != nil {
        value.Address = *p.Address
    }
	if p.Rating != nil {
        value.Rating = *p.Rating
    }

    r.places[id] = value

    return nil
}

func (r *MemoryPlaceRepository) Delete(ctx context.Context, id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, ok := r.places[id]; !ok {
        return er.ErrPlaceNotFound
    }

    delete(r.places, id)
    return nil
}

func (r *MemoryPlaceRepository) DeleteAll(ctx context.Context) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.places = make(map[string]models.Place)
    return nil
}