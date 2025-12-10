package repository

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"

	"github.com/smilu97/refana/internal/pkg/domain"
)

type DataSourceRepository struct {
	db *gorm.DB
}

func NewDataSourceRepository(db *gorm.DB) *DataSourceRepository {
	return &DataSourceRepository{db: db}
}

func (r *DataSourceRepository) Create(ctx context.Context, ds domain.DataSource) error {
	if ds.UpdatedAt.IsZero() {
		ds.UpdatedAt = time.Now()
	}
	var model dataSourceModel
	if err := toDataSourceModel(ds, &model); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *DataSourceRepository) Get(ctx context.Context, id domain.DataSourceID) (domain.DataSource, error) {
	var model dataSourceModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id.Int64()).Error; err != nil {
		return domain.DataSource{}, err
	}
	return toDataSourceDomain(model)
}

// Update applies last-write-wins using provided updatedAt timestamp.
func (r *DataSourceRepository) Update(ctx context.Context, ds domain.DataSource, updatedAt time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing dataSourceModel
		if err := tx.First(&existing, "id = ?", ds.ID.Int64()).Error; err != nil {
			return err
		}
		if !updatedAt.After(existing.UpdatedAt) {
			return nil
		}

		ds.UpdatedAt = updatedAt
		var model dataSourceModel
		if err := toDataSourceModel(ds, &model); err != nil {
			return err
		}
		return tx.Model(&existing).Updates(model).Error
	})
}

func (r *DataSourceRepository) Delete(ctx context.Context, id domain.DataSourceID) error {
	return r.db.WithContext(ctx).Delete(&dataSourceModel{}, "id = ?", id.Int64()).Error
}

func (r *DataSourceRepository) List(ctx context.Context) ([]domain.DataSource, error) {
	var models []dataSourceModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.DataSource, 0, len(models))
	for _, m := range models {
		ds, err := toDataSourceDomain(m)
		if err != nil {
			return nil, err
		}
		out = append(out, ds)
	}
	return out, nil
}

// Storage model for data_sources.
type dataSourceModel struct {
	ID             int64 `gorm:"primaryKey;autoIncrement:false"`
	ClassID        string
	Name           string
	Alias          string
	PropertiesJSON string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (dataSourceModel) TableName() string { return "data_sources" }

func toDataSourceModel(src domain.DataSource, dst *dataSourceModel) error {
	props, err := json.Marshal(src.Properties)
	if err != nil {
		return err
	}
	dst.ID = src.ID.Int64()
	dst.ClassID = string(src.ClassID)
	dst.Name = string(src.Name)
	dst.Alias = string(src.Alias)
	dst.PropertiesJSON = string(props)
	dst.UpdatedAt = src.UpdatedAt
	return nil
}

func toDataSourceDomain(m dataSourceModel) (domain.DataSource, error) {
	var props map[domain.PropertyKey]domain.PropertyValue
	if err := json.Unmarshal([]byte(m.PropertiesJSON), &props); err != nil {
		return domain.DataSource{}, err
	}
	return domain.DataSource{
		ID:         domain.NewDataSourceID(m.ID),
		ClassID:    domain.DataSourceClassID(m.ClassID),
		Name:       domain.Name(m.Name),
		Alias:      domain.Alias(m.Alias),
		Properties: props,
		UpdatedAt:  m.UpdatedAt,
	}, nil
}
