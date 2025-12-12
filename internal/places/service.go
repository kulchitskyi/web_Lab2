package places

import (
    "time"
    "context"
    "sync"
	"github.com/google/uuid"

    er "deu/internal/errors"
    repo "deu/internal/repository"
	"deu/internal/models"
)

type PlaceService struct {
    repo        repo.PlaceRepository
    enableCache bool
    cache       map[string]*models.Place
    mu          sync.RWMutex
}

func NewPlaceService(repo repo.PlaceRepository, enableCache bool) *PlaceService {
    return &PlaceService{
        repo:        repo,
        enableCache: enableCache,
        cache:       make(map[string]*models.Place),
    }
}

func (s *PlaceService) GetAll(ctx context.Context) ([]models.Place, error) {
    return s.repo.GetAll(ctx)
}

func (s *PlaceService) GetById(ctx context.Context, id string) (*models.Place, error) {
    if id == "" {
        return nil, er.ErrInvalidPlaceData
    }

    if s.enableCache {
        s.mu.RLock()
        if place, found := s.cache[id]; found {
            s.mu.RUnlock()
            return place, nil
        }
        s.mu.RUnlock()
    }

    place, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, err
    }

    if s.enableCache {
        s.mu.Lock()
        s.cache[id] = place
        s.mu.Unlock()
    }

    return place, nil
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

	err := s.repo.Create(ctx, &place)
	if err != nil {
		return nil, err
	}

	if s.enableCache {
		s.mu.Lock()
		s.cache[place.Id] = &place
		s.mu.Unlock()
	}

    return &place, nil
}

func (s *PlaceService) Update(ctx context.Context, id string, p *models.PlaceUpdateRequest) error {

    if id == "" {
        return er.ErrInvalidPlaceData
    }

    err := s.repo.Update(ctx, id, p)
    if err != nil {
        return err
    }

    if s.enableCache {
        s.mu.Lock()
        delete(s.cache, id)
        s.mu.Unlock()
    }

    return nil
}

func (s *PlaceService) DeleteById(ctx context.Context, id string) error {
    if id == "" {
        return er.ErrInvalidPlaceData
    }
    
    err := s.repo.Delete(ctx, id)
    if err != nil {
        return err
    }

    if s.enableCache {
        s.mu.Lock()
        delete(s.cache, id)
        s.mu.Unlock()
    }

    return nil
}

func (s *PlaceService) DeleteAll(ctx context.Context) error {
    err := s.repo.DeleteAll(ctx)
    if err != nil {
        return err
    }
    
    if s.enableCache {
        s.mu.Lock()
        s.cache = make(map[string]*models.Place)
        s.mu.Unlock()
    }
    
    return nil
}