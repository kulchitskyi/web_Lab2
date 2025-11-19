package places

import (
    "time"
    "context"
	"github.com/google/uuid"

    er "deu/internal/errors"
    repo "deu/internal/repository"
	"deu/internal/models"
)

type PlaceService struct {
    repo repo.PlaceRepository
}

func NewPlaceService(repo repo.PlaceRepository) *PlaceService {
    return &PlaceService{repo: repo}
}

func (s *PlaceService) GetAll(ctx context.Context) ([]models.Place, error) {
    return s.repo.GetAll(ctx)
}

func (s *PlaceService) GetById(ctx context.Context, id string) (*models.Place, error) {
    if id == "" {
        return nil, er.ErrInvalidPlaceData
    }
    return s.repo.GetByID(ctx, id)
}

func (s *PlaceService) Create(ctx context.Context, p *models.PlaceCreateRequest) (*models.Place, error) {

    if p.Name == "" {
        return nil, er.ErrInvalidPlaceData
    }

	id := uuid.New()

	place := models.Place{
		Id: 			id.String(),
		Name: 			p.Name,
		Description:	p.Description,
		Location: 		p.Location,
		Address: 		p.Address,
		Rating: 		p.Rating,
		CreatedAt: 		time.Now(),
	}


    return &place, s.repo.Create(ctx, &place)
}

func (s *PlaceService) Update(ctx context.Context, id string, p *models.PlaceUpdateRequest) error {

    if id == "" {
        return er.ErrInvalidPlaceData
    }

    return s.repo.Update(ctx, id, p)
}

func (s *PlaceService) DeleteById(ctx context.Context, id string) error {
    if id == "" {
        return er.ErrInvalidPlaceData
    }
    return s.repo.Delete(ctx, id)
}

func (s *PlaceService) DeleteAll(ctx context.Context) error {
    return s.repo.DeleteAll(ctx)
}