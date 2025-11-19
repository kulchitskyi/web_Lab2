package repository

import (
	"deu/internal/models"
	"context"
    "sync"

    er "deu/internal/errors"
)

type MemoryUserRepository struct {
    mu      sync.RWMutex
    users  map[string]models.User
}

func NewMemoryUserRepository() *MemoryUserRepository {
    return &MemoryUserRepository{
        users: make(map[string]models.User),
    }
}

func (r *MemoryUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    result := make([]models.User, 0, len(r.users))
    for _, u := range r.users {
        result = append(result, u)
    }
    return result, nil
}

func (r *MemoryUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    result, ok := r.users[id]
    if ok {
        return &result, nil
    }

    return nil, er.ErrUserNotFound
}

func (r *MemoryUserRepository) Create(ctx context.Context, u *models.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.users[u.Id] = *u
    return nil
}

func (r *MemoryUserRepository) Update(ctx context.Context, id string, u *models.UserUpdateRequest) error { //impove
    r.mu.Lock()
    defer r.mu.Unlock()

    value, ok := r.users[id]
    if !ok {
        return er.ErrUserNotFound
    }

    if u.Name != nil {
        value.Name = *u.Name
    }
    if u.Email != nil {
        value.Email = *u.Email
    }

    r.users[id] = value

    return nil
}

func (r *MemoryUserRepository) Delete(ctx context.Context, id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, ok := r.users[id]; !ok {
        return er.ErrUserNotFound
    }

    delete(r.users, id)
    return nil
}

func (r *MemoryUserRepository) DeleteAll(ctx context.Context) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    r.users = make(map[string]models.User)
    return nil
}