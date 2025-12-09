package storage_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/smilu97/refana/internal/storage"
)

// TDD for TODO #5: SQLite + GORM AutoMigrate pipeline.

func TestMigrateCreatesTables(t *testing.T) {
	db := openInMemoryDB(t)

	if err := storage.Migrate(db); err != nil {
		t.Fatalf("Migrate error: %v", err)
	}

	expectTables(t, db, "components", "data_sources")
}

func TestMigrateIsIdempotent(t *testing.T) {
	db := openInMemoryDB(t)

	if err := storage.Migrate(db); err != nil {
		t.Fatalf("first migrate error: %v", err)
	}
	if err := storage.Migrate(db); err != nil {
		t.Fatalf("second migrate error: %v", err)
	}
}

// Helpers
func openInMemoryDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	return db
}

func expectTables(t *testing.T, db *gorm.DB, tables ...string) {
	t.Helper()
	for _, tbl := range tables {
		if !db.Migrator().HasTable(tbl) {
			t.Fatalf("expected table %q to exist after migration", tbl)
		}
	}
}
