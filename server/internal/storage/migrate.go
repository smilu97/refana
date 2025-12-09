package storage

import (
	"time"

	"gorm.io/gorm"
)

// Migrate applies database schema using GORM's AutoMigrate.
// Idempotent: safe to run multiple times.
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&componentModel{},
		&dataSourceModel{},
	)
}

// Minimal models to establish tables. Columns are intentionally
// narrow; evolve alongside repository layer as needed.
type componentModel struct {
	ID               string `gorm:"primaryKey"`
	VisualisationID  string
	QueryJSON        string
	Name             string
	CoordinationJSON string
	PropertiesJSON   string
	UpdatedAt        time.Time
	CreatedAt        time.Time
}

func (componentModel) TableName() string { return "components" }

type dataSourceModel struct {
	ID             string `gorm:"primaryKey"`
	ClassID        string
	Name           string
	Alias          string
	PropertiesJSON string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (dataSourceModel) TableName() string { return "data_sources" }
