package repository_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/smilu97/refana/internal/pkg/domain"
	"github.com/smilu97/refana/internal/repository"
	"github.com/smilu97/refana/internal/storage"
)

func TestDataSourceRepositoryCRUD_LastWriteWins(t *testing.T) {
	db := openDSDB(t)
	ctx := context.Background()
	repo := repository.NewDataSourceRepository(db)

	now := time.Now()
	ds := domain.DataSource{
		ID:      domain.NewDataSourceID(1),
		ClassID: "postgres",
		Name:    "Main DS",
		Alias:   "primary",
		Properties: map[domain.PropertyKey]domain.PropertyValue{
			"host": "localhost",
		},
	}

	// Create
	if err := repo.Create(ctx, ds); err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Get
	got, err := repo.Get(ctx, ds.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != ds.Name {
		t.Fatalf("Get Name = %s, want %s", got.Name, ds.Name)
	}

	// Update newer
	newer := ds
	newer.Name = "Updated"
	newerAlias := "updated-alias"
	newer.Alias = domain.Alias(newerAlias)
	newerUpdated := now.Add(time.Minute)
	if err := repo.Update(ctx, newer, newerUpdated); err != nil {
		t.Fatalf("Update newer: %v", err)
	}
	got, _ = repo.Get(ctx, ds.ID)
	if got.Name != "Updated" || got.Alias != domain.Alias(newerAlias) {
		t.Fatalf("after newer update = (%s,%s), want (%s,%s)", got.Name, got.Alias, "Updated", newerAlias)
	}

	// Update older should be ignored
	older := newer
	older.Name = "ShouldNotPersist"
	olderUpdated := now.Add(-time.Minute)
	if err := repo.Update(ctx, older, olderUpdated); err != nil {
		t.Fatalf("Update older: %v", err)
	}
	got, _ = repo.Get(ctx, ds.ID)
	if got.Name != "Updated" {
		t.Fatalf("after older update Name = %s, want %s", got.Name, "Updated")
	}

	// Delete
	if err := repo.Delete(ctx, ds.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := repo.Get(ctx, ds.ID); err == nil {
		t.Fatalf("expected error after delete, got nil")
	}
}

func TestDataSourceRepositoryList(t *testing.T) {
	db := openDSDB(t)
	ctx := context.Background()
	repo := repository.NewDataSourceRepository(db)

	for i := 0; i < 2; i++ {
		if err := repo.Create(ctx, domain.DataSource{
			ID:      domain.NewDataSourceID(int64(i + 1)),
			ClassID: "postgres",
			Name:    domain.Name("ds"),
			Alias:   domain.Alias("a"),
			Properties: map[domain.PropertyKey]domain.PropertyValue{
				"host": "localhost",
			},
		}); err != nil {
			t.Fatalf("Create %d: %v", i, err)
		}
	}

	all, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("List len = %d, want 2", len(all))
	}
}

// Helpers
func openDSDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:ds-repo-%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := storage.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}
