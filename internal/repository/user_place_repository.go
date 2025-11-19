package repository

import (
	"context"
)

type UserPlaceRepository interface {
    AddVisitedPlace(ctx context.Context, userID, placeID string) error
    HasVisitedPlace(ctx context.Context, userID, placeID string) (bool, error)
}