package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/smilu97/refana/internal/pkg/domain"
	"github.com/smilu97/refana/internal/repository"
)

type ComponentService struct {
	repo *repository.ComponentRepository
}

func NewComponentService(repo *repository.ComponentRepository) *ComponentService {
	return &ComponentService{repo: repo}
}

func (s *ComponentService) Create(ctx context.Context, opts domain.CreateComponentOptions) (domain.Component, error) {
	if opts.Name == "" || opts.VisualisationID == "" {
		return domain.Component{}, ErrBadRequest
	}
	var query domain.Query
	if len(opts.Queries) > 0 {
		query = opts.Queries[0]
	}

	comp := domain.Component{
		ID:              domain.NewComponentID(time.Now().UnixNano()),
		VisualisationID: opts.VisualisationID,
		Query:           query,
		Name:            opts.Name,
		Coordination:    opts.Coordination,
		Properties:      opts.Properties,
		UpdatedAt:       time.Now(),
	}

	if err := s.repo.Create(ctx, comp); err != nil {
		return domain.Component{}, err
	}
	return comp, nil
}

func (s *ComponentService) Get(ctx context.Context, id domain.ComponentID) (domain.Component, error) {
	c, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Component{}, ErrNotFound
		}
		return domain.Component{}, err
	}
	return c, nil
}

func (s *ComponentService) List(ctx context.Context) ([]domain.Component, error) {
	return s.repo.List(ctx)
}

func (s *ComponentService) Update(
	ctx context.Context,
	id domain.ComponentID,
	opts domain.UpdateComponentOptions,
	updatedAt time.Time,
) error {
	if opts.Name == "" || opts.VisualisationID == "" {
		return ErrBadRequest
	}

	var query domain.Query
	if len(opts.Queries) > 0 {
		query = opts.Queries[0]
	}

	comp := domain.Component{
		ID:              id,
		VisualisationID: opts.VisualisationID,
		Query:           query,
		Name:            opts.Name,
		Coordination:    opts.Coordination,
		Properties:      opts.Properties,
		UpdatedAt:       updatedAt,
	}
	if err := s.repo.Update(ctx, comp); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *ComponentService) Delete(ctx context.Context, id domain.ComponentID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
