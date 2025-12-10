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
		&dataSourceClassModel{},
	)
}

// Minimal models to establish tables. Columns are intentionally
// narrow; evolve alongside repository layer as needed.
type componentModel struct {
	ID               int64     `gorm:"primaryKey;autoIncrement:false"`
	VisualisationID  string    `gorm:"size:64;index"`
	QueryJSON        string    `gorm:"type:text"`
	Name             string    `gorm:"size:256"`
	CoordinationJSON string    `gorm:"type:text"`
	PropertiesJSON   string    `gorm:"type:text"`
	UpdatedAt        time.Time `gorm:"index"`
	CreatedAt        time.Time
}

func (componentModel) TableName() string { return "components" }

type dataSourceModel struct {
	ID             int64     `gorm:"primaryKey;autoIncrement:false"`
	ClassID        string    `gorm:"size:64;index"`
	Name           string    `gorm:"size:256"`
	Alias          string    `gorm:"size:64;index"`
	PropertiesJSON string    `gorm:"type:text"`
	CreatedAt      time.Time
	UpdatedAt      time.Time `gorm:"index"`
}

func (dataSourceModel) TableName() string { return "data_sources" }

type dataSourceClassModel struct {
	ID                      int64  `gorm:"primaryKey;size:64"`
	Name                    string `gorm:"size:256"`
	PropertyDescriptorsJSON string `gorm:"type:text"`
	CreatedAt               time.Time
	UpdatedAt               time.Time `gorm:"index"`
}

func (dataSourceClassModel) TableName() string { return "data_source_classes" }
