package repository

import (
	"context"
	"deu/internal/models"
	"log/slog"
	"time"
)

type LoggingUserRepository struct {
	Repo   UserRepository
	Logger *slog.Logger
}

func NewLoggingUserRepository(repo UserRepository, logger *slog.Logger) *LoggingUserRepository {
	return &LoggingUserRepository{
		Repo:   repo,
		Logger: logger,
	}
}

func (r *LoggingUserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	r.Logger.Info("Calling GetAll Users")
	start := time.Now()
	users, err := r.Repo.GetAll(ctx)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("GetAll Users failed", "error", err, "duration", duration)
		return nil, err
	}
	r.Logger.Info("GetAll Users success", "count", len(users), "duration", duration)
	return users, nil
}

func (r *LoggingUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	r.Logger.Info("Calling GetByID User", "id", id)
	start := time.Now()
	user, err := r.Repo.GetByID(ctx, id)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("GetByID User failed", "id", id, "error", err, "duration", duration)
		return nil, err
	}
	r.Logger.Info("GetByID User success", "id", id, "duration", duration)
	return user, nil
}

func (r *LoggingUserRepository) Create(ctx context.Context, u *models.User) error {
	r.Logger.Info("Calling Create User", "email", u.Email)
	start := time.Now()
	err := r.Repo.Create(ctx, u)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("Create User failed", "email", u.Email, "error", err, "duration", duration)
		return err
	}
	r.Logger.Info("Create User success", "id", u.Id, "duration", duration)
	return nil
}

func (r *LoggingUserRepository) Update(ctx context.Context, id string, u *models.UserUpdateRequest) error {
	r.Logger.Info("Calling Update User", "id", id)
	start := time.Now()
	err := r.Repo.Update(ctx, id, u)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("Update User failed", "id", id, "error", err, "duration", duration)
		return err
	}
	r.Logger.Info("Update User success", "id", id, "duration", duration)
	return nil
}

func (r *LoggingUserRepository) Delete(ctx context.Context, id string) error {
	r.Logger.Info("Calling Delete User", "id", id)
	start := time.Now()
	err := r.Repo.Delete(ctx, id)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("Delete User failed", "id", id, "error", err, "duration", duration)
		return err
	}
	r.Logger.Info("Delete User success", "id", id, "duration", duration)
	return nil
}

func (r *LoggingUserRepository) DeleteAll(ctx context.Context) error {
	r.Logger.Info("Calling DeleteAll Users")
	start := time.Now()
	err := r.Repo.DeleteAll(ctx)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("DeleteAll Users failed", "error", err, "duration", duration)
		return err
	}
	r.Logger.Info("DeleteAll Users success", "duration", duration)
	return nil
}

type LoggingPlaceRepository struct {
	Repo   PlaceRepository
	Logger *slog.Logger
}

func NewLoggingPlaceRepository(repo PlaceRepository, logger *slog.Logger) *LoggingPlaceRepository {
	return &LoggingPlaceRepository{
		Repo:   repo,
		Logger: logger,
	}
}

func (r *LoggingPlaceRepository) GetAll(ctx context.Context) ([]models.Place, error) {
	r.Logger.Info("Calling GetAll Places")
	start := time.Now()
	places, err := r.Repo.GetAll(ctx)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("GetAll Places failed", "error", err, "duration", duration)
		return nil, err
	}
	r.Logger.Info("GetAll Places success", "count", len(places), "duration", duration)
	return places, nil
}

func (r *LoggingPlaceRepository) GetByID(ctx context.Context, id string) (*models.Place, error) {
	r.Logger.Info("Calling GetByID Place", "id", id)
	start := time.Now()
	place, err := r.Repo.GetByID(ctx, id)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("GetByID Place failed", "id", id, "error", err, "duration", duration)
		return nil, err
	}
	r.Logger.Info("GetByID Place success", "id", id, "duration", duration)
	return place, nil
}

func (r *LoggingPlaceRepository) Create(ctx context.Context, p *models.Place) error {
	r.Logger.Info("Calling Create Place", "name", p.Name)
	start := time.Now()
	err := r.Repo.Create(ctx, p)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("Create Place failed", "name", p.Name, "error", err, "duration", duration)
		return err
	}
	r.Logger.Info("Create Place success", "id", p.Id, "duration", duration)
	return nil
}

func (r *LoggingPlaceRepository) Update(ctx context.Context, id string, p *models.PlaceUpdateRequest) error {
	r.Logger.Info("Calling Update Place", "id", id)
	start := time.Now()
	err := r.Repo.Update(ctx, id, p)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("Update Place failed", "id", id, "error", err, "duration", duration)
		return err
	}
	r.Logger.Info("Update Place success", "id", id, "duration", duration)
	return nil
}

func (r *LoggingPlaceRepository) Delete(ctx context.Context, id string) error {
	r.Logger.Info("Calling Delete Place", "id", id)
	start := time.Now()
	err := r.Repo.Delete(ctx, id)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("Delete Place failed", "id", id, "error", err, "duration", duration)
		return err
	}
	r.Logger.Info("Delete Place success", "id", id, "duration", duration)
	return nil
}

func (r *LoggingPlaceRepository) DeleteAll(ctx context.Context) error {
	r.Logger.Info("Calling DeleteAll Places")
	start := time.Now()
	err := r.Repo.DeleteAll(ctx)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("DeleteAll Places failed", "error", err, "duration", duration)
		return err
	}
	r.Logger.Info("DeleteAll Places success", "duration", duration)
	return nil
}

type LoggingUserPlaceRepository struct {
	Repo   UserPlaceRepository
	Logger *slog.Logger
}

func NewLoggingUserPlaceRepository(repo UserPlaceRepository, logger *slog.Logger) *LoggingUserPlaceRepository {
	return &LoggingUserPlaceRepository{
		Repo:   repo,
		Logger: logger,
	}
}

func (r *LoggingUserPlaceRepository) AddVisitedPlace(ctx context.Context, userID, placeID string) error {
	r.Logger.Info("Calling AddVisitedPlace", "userID", userID, "placeID", placeID)
	start := time.Now()
	err := r.Repo.AddVisitedPlace(ctx, userID, placeID)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("AddVisitedPlace failed", "userID", userID, "placeID", placeID, "error", err, "duration", duration)
		return err
	}
	r.Logger.Info("AddVisitedPlace success", "userID", userID, "placeID", placeID, "duration", duration)
	return nil
}

func (r *LoggingUserPlaceRepository) RemoveVisitedPlace(ctx context.Context, userID, placeID string) error {
	r.Logger.Info("Calling RemoveVisitedPlace", "userID", userID, "placeID", placeID)
	start := time.Now()
	err := r.Repo.RemoveVisitedPlace(ctx, userID, placeID)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("RemoveVisitedPlace failed", "userID", userID, "placeID", placeID, "error", err, "duration", duration)
		return err
	}
	r.Logger.Info("RemoveVisitedPlace success", "userID", userID, "placeID", placeID, "duration", duration)
	return nil
}

func (r *LoggingUserPlaceRepository) HasVisitedPlace(ctx context.Context, userID, placeID string) (bool, error) {
	r.Logger.Info("Calling HasVisitedPlace", "userID", userID, "placeID", placeID)
	start := time.Now()
	visited, err := r.Repo.HasVisitedPlace(ctx, userID, placeID)
	duration := time.Since(start)
	if err != nil {
		r.Logger.Error("HasVisitedPlace failed", "userID", userID, "placeID", placeID, "error", err, "duration", duration)
		return false, err
	}
	r.Logger.Info("HasVisitedPlace success", "userID", userID, "placeID", placeID, "visited", visited, "duration", duration)
	return visited, nil
}
