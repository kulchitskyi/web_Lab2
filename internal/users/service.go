package users

import (
	"time"
    "context"
	"github.com/google/uuid"
    
	er "deu/internal/errors"
    repo "deu/internal/repository"
	"deu/internal/models"
)

type UserService struct {
    repo repo.UserRepository
    userPlaceRepo repo.UserPlaceRepository 
    placeRepo repo.PlaceRepository
}

func NewUserService(userRepo repo.UserRepository, userPlaceRepo repo.UserPlaceRepository, placeRepo repo.PlaceRepository) *UserService {
    return &UserService{
        repo: userRepo, 
        userPlaceRepo: userPlaceRepo, 
        placeRepo: placeRepo,
    }
}

func (s *UserService) GetAll(ctx context.Context) ([]models.User, error) {
    return s.repo.GetAll(ctx)
}

func (s *UserService) GetById(ctx context.Context, id string) (*models.User, error) {
    if id == "" {
        return nil, er.ErrInvalidUserData
    }
    return s.repo.GetByID(ctx, id)
}

func (s *UserService) Create(ctx context.Context, u *models.UserCreateRequest) (*models.User, error) {

    if u.Name == "" || u.Email == "" {
        return nil, er.ErrInvalidUserData
    }

	id := uuid.New()

	user := models.User{
		Id: 		id.String(),
		Name: 		u.Name,
		Email: 		u.Email,
		CreatedAt:	time.Now(),
	}


    return &user, s.repo.Create(ctx, &user)
}

func (s *UserService) Update(ctx context.Context, id string, u *models.UserUpdateRequest) error {

    if id == "" {
        return er.ErrInvalidUserData
    }

    return s.repo.Update(ctx, id, u)
}

func (s *UserService) DeleteById(ctx context.Context, id string) error {
    if id == "" {
        return er.ErrInvalidUserData
    }
    return s.repo.Delete(ctx, id)
}

func (s *UserService) AddVisitedPlace(ctx context.Context, userID, placeID string) error {
    if userID == "" || placeID == "" {
        return er.ErrInvalidUserData
    }

    _, userErr := s.repo.GetByID(ctx, userID)
    if userErr != nil {
        return userErr
    }
    
    _, placeErr := s.placeRepo.GetByID(ctx, placeID)
    if placeErr != nil {
        return placeErr
    }

    return s.userPlaceRepo.AddVisitedPlace(ctx, userID, placeID)
}

func (s *UserService) HasVisitedPlace(ctx context.Context, userID, placeID string) (bool, error) {
    if userID == "" || placeID == "" {
        return false, er.ErrInvalidUserData
    }

    _, userErr := s.repo.GetByID(ctx, userID)
    if userErr != nil {
        return false, userErr
    }
    
    _, placeErr := s.placeRepo.GetByID(ctx, placeID)
    if placeErr != nil {
        return false, placeErr
    }

    return s.userPlaceRepo.HasVisitedPlace(ctx, userID, placeID)
}

func (s *UserService) RemoveVisitedPlace(ctx context.Context, userID, placeID string) error {
    if userID == "" || placeID == "" {
        return er.ErrInvalidUserData
    }

    _, userErr := s.repo.GetByID(ctx, userID)
    if userErr != nil {
        return userErr
    }
    
    _, placeErr := s.placeRepo.GetByID(ctx, placeID)
    if placeErr != nil {
        return placeErr
    }

    return s.userPlaceRepo.RemoveVisitedPlace(ctx, userID, placeID)
}

func (s *UserService) DeleteAll(ctx context.Context) error {
    return s.repo.DeleteAll(ctx)
}