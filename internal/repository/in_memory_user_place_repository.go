package repository

import (
	"context"
	"sync"
)


type MemoryUserPlaceRepository struct {
	visitedMap map[string]map[string]bool
	mu         sync.RWMutex
}

func NewMemoryUserPlaceRepository() *MemoryUserPlaceRepository {
	return &MemoryUserPlaceRepository{
		visitedMap: make(map[string]map[string]bool),
	}
}

func (r *MemoryUserPlaceRepository) AddVisitedPlace(ctx context.Context, userID, placeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.visitedMap[userID] == nil {
		r.visitedMap[userID] = make(map[string]bool)
	}

	r.visitedMap[userID][placeID] = true
	
	return nil
}

func (r *MemoryUserPlaceRepository) HasVisitedPlace(ctx context.Context, userID, placeID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	places, ok := r.visitedMap[userID]
	if !ok {
		return false, nil
	}

	_, visited := places[placeID]
	
	return visited, nil
}