package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/smilu97/refana/internal/pkg/domain"
)

type ComponentRepository struct {
	db *gorm.DB
}

func NewComponentRepository(db *gorm.DB) *ComponentRepository {
	return &ComponentRepository{db: db}
}

func (r *ComponentRepository) Create(ctx context.Context, comp domain.Component) error {
	var model componentModel
	if err := toComponentModel(comp, &model); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(&model).Error
}

func (r *ComponentRepository) Get(ctx context.Context, id domain.ComponentID) (domain.Component, error) {
	var model componentModel
	if err := r.db.WithContext(ctx).First(&model, "id = ?", id.Int64()).Error; err != nil {
		return domain.Component{}, err
	}
	return toComponentDomain(model)
}

func (r *ComponentRepository) Update(ctx context.Context, comp domain.Component) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing componentModel
		if err := tx.First(&existing, "id = ?", comp.ID.Int64()).Error; err != nil {
			return err
		}

		// last-write-wins: only update if incoming UpdatedAt is newer
		if !comp.UpdatedAt.After(existing.UpdatedAt) {
			return nil
		}

		var updated componentModel
		if err := toComponentModel(comp, &updated); err != nil {
			return err
		}
		return tx.Model(&existing).Updates(updated).Error
	})
}

func (r *ComponentRepository) Delete(ctx context.Context, id domain.ComponentID) error {
	return r.db.WithContext(ctx).Delete(&componentModel{}, "id = ?", id.Int64()).Error
}

func (r *ComponentRepository) List(ctx context.Context) ([]domain.Component, error) {
	var models []componentModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Component, 0, len(models))
	for _, m := range models {
		c, err := toComponentDomain(m)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

// Storage model for components table.
type componentModel struct {
	ID               int64 `gorm:"primaryKey;autoIncrement:false"`
	VisualisationID  string
	QueryJSON        string
	Name             string
	CoordinationJSON string
	PropertiesJSON   string
	UpdatedAt        time.Time
	CreatedAt        time.Time
}

func (componentModel) TableName() string { return "components" }

func toComponentModel(src domain.Component, dst *componentModel) error {
	queryBytes, err := json.Marshal(src.Query)
	if err != nil {
		return err
	}
	coordBytes, err := json.Marshal(src.Coordination)
	if err != nil {
		return err
	}
	propsBytes, err := json.Marshal(src.Properties)
	if err != nil {
		return err
	}

	dst.ID = src.ID.Int64()
	dst.VisualisationID = string(src.VisualisationID)
	dst.QueryJSON = string(queryBytes)
	dst.Name = string(src.Name)
	dst.CoordinationJSON = string(coordBytes)
	dst.PropertiesJSON = string(propsBytes)
	dst.UpdatedAt = src.UpdatedAt
	// CreatedAt is managed by GORM; leave zero to auto-set.
	return nil
}

func toComponentDomain(m componentModel) (domain.Component, error) {
	var query domain.Query
	if err := json.Unmarshal([]byte(m.QueryJSON), &query); err != nil {
		return domain.Component{}, err
	}
	var coord domain.Coordination
	if err := json.Unmarshal([]byte(m.CoordinationJSON), &coord); err != nil {
		return domain.Component{}, err
	}
	var props map[domain.PropertyKey]domain.PropertyValue
	if err := json.Unmarshal([]byte(m.PropertiesJSON), &props); err != nil {
		return domain.Component{}, err
	}

	return domain.Component{
		ID:              domain.NewComponentID(m.ID),
		VisualisationID: domain.VisualisationID(m.VisualisationID),
		Query:           query,
		Name:            domain.Name(m.Name),
		Coordination:    coord,
		Properties:      props,
		UpdatedAt:       m.UpdatedAt,
	}, nil
}

// Ensure interface compliance with errors.Is on not found cases.
var ErrNotFound = errors.New("component not found")
