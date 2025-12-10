package service

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/smilu97/refana/internal/pkg/domain"
	"github.com/smilu97/refana/internal/repository"
)

type DataSourceService struct {
	repo *repository.DataSourceRepository
}

func NewDataSourceService(repo *repository.DataSourceRepository) *DataSourceService {
	return &DataSourceService{repo: repo}
}

func (s *DataSourceService) Create(ctx context.Context, opts domain.CreateDataSourceOptions) (domain.DataSource, error) {
	if opts.Name == "" || opts.ClassID == "" {
		return domain.DataSource{}, ErrBadRequest
	}
	ds := domain.DataSource{
		ID:         domain.NewDataSourceID(time.Now().UnixNano()),
		ClassID:    opts.ClassID,
		Name:       opts.Name,
		Alias:      opts.Alias,
		Properties: opts.Properties,
		UpdatedAt:  time.Now(),
	}
	if err := s.repo.Create(ctx, ds); err != nil {
		return domain.DataSource{}, err
	}
	return ds, nil
}

func (s *DataSourceService) Get(ctx context.Context, id domain.DataSourceID) (domain.DataSource, error) {
	ds, err := s.repo.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.DataSource{}, ErrNotFound
		}
		return domain.DataSource{}, err
	}
	return ds, nil
}

func (s *DataSourceService) List(ctx context.Context) ([]domain.DataSource, error) {
	return s.repo.List(ctx)
}

func (s *DataSourceService) Update(
	ctx context.Context,
	id domain.DataSourceID,
	opts domain.UpdateDataSourceOptions,
	updatedAt time.Time,
) error {
	if opts.Name == "" || opts.ClassID == "" {
		return ErrBadRequest
	}
	ds := domain.DataSource{
		ID:         id,
		ClassID:    opts.ClassID,
		Name:       opts.Name,
		Alias:      opts.Alias,
		Properties: opts.Properties,
		UpdatedAt:  updatedAt,
	}
	if err := s.repo.Update(ctx, ds, updatedAt); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *DataSourceService) Delete(ctx context.Context, id domain.DataSourceID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
