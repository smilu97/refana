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

func TestComponentRepositoryCRUD_LastWriteWins(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	repo := repository.NewComponentRepository(db)

	now := time.Now()
	comp := domain.Component{
		ID:              domain.NewComponentID(1),
		VisualisationID: "table",
		Query: domain.Query{
			Name:         "main",
			DataSourceID: domain.NewDataSourceID(1),
			Properties: map[domain.PropertyKey]domain.PropertyValue{
				"sql": "select 1",
			},
		},
		Name:         "Test Component",
		Coordination: domain.Coordination{Rect: domain.Rect{Left: 0, Top: 0, Width: 10, Height: 10}, ZIndex: 1},
		Properties: map[domain.PropertyKey]domain.PropertyValue{
			"title": "hello",
		},
		UpdatedAt: now,
	}

	// Create
	if err := repo.Create(ctx, comp); err != nil {
		t.Fatalf("Create: %v", err)
	}

	// Get
	got, err := repo.Get(ctx, comp.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != comp.Name {
		t.Fatalf("Get Name = %s, want %s", got.Name, comp.Name)
	}

	// Update with newer timestamp should override
	newer := comp
	newer.Name = "Updated Name"
	newer.UpdatedAt = comp.UpdatedAt.Add(time.Minute)
	if err := repo.Update(ctx, newer); err != nil {
		t.Fatalf("Update newer: %v", err)
	}
	got, _ = repo.Get(ctx, comp.ID)
	if got.Name != "Updated Name" {
		t.Fatalf("after newer update Name = %s, want %s", got.Name, "Updated Name")
	}

	// Update with older timestamp should be ignored (last-write-wins)
	older := newer
	older.Name = "Should Not Persist"
	older.UpdatedAt = comp.UpdatedAt.Add(-time.Minute)
	if err := repo.Update(ctx, older); err != nil {
		t.Fatalf("Update older: %v", err)
	}
	got, _ = repo.Get(ctx, comp.ID)
	if got.Name != "Updated Name" {
		t.Fatalf("after older update Name = %s, want %s", got.Name, "Updated Name")
	}

	// Delete
	if err := repo.Delete(ctx, comp.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := repo.Get(ctx, comp.ID); err == nil {
		t.Fatalf("expected error after delete, got nil")
	}
}

func TestComponentRepositoryList(t *testing.T) {
	db := openTestDB(t)
	ctx := context.Background()
	repo := repository.NewComponentRepository(db)

	base := time.Now()
	for i := 0; i < 3; i++ {
		if err := repo.Create(ctx, domain.Component{
			ID:              domain.NewComponentID(int64(i + 1)),
			VisualisationID: "table",
			Query:           domain.Query{Name: "q"},
			Name:            domain.Name("c"),
			Coordination:    domain.Coordination{},
			Properties:      map[domain.PropertyKey]domain.PropertyValue{},
			UpdatedAt:       base,
		}); err != nil {
			t.Fatalf("Create %d: %v", i, err)
		}
	}

	all, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("List len = %d, want 3", len(all))
	}
}

// Helpers
func openTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := fmt.Sprintf("file:comp-repo-%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := storage.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}
