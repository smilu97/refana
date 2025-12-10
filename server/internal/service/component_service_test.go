package service_test

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/smilu97/refana/internal/pkg/domain"
	"github.com/smilu97/refana/internal/repository"
	"github.com/smilu97/refana/internal/service"
	"github.com/smilu97/refana/internal/storage"
)

func TestComponentService_CreateValidate(t *testing.T) {
	svc := newComponentService(t)
	ctx := context.Background()

	_, err := svc.Create(ctx, domain.CreateComponentOptions{
		Name:            "",
		VisualisationID: "",
	})
	if err != service.ErrBadRequest {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestComponentService_CRUD_LastWriteWins(t *testing.T) {
	svc := newComponentService(t)
	ctx := context.Background()

	comp, err := svc.Create(ctx, domain.CreateComponentOptions{
		Name:            "comp",
		VisualisationID: "table",
		Queries: []domain.Query{{
			Name:         "main",
			DataSourceID: domain.NewDataSourceID(1),
		}},
		Properties: map[domain.PropertyKey]domain.PropertyValue{"p": "v"},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := svc.Get(ctx, comp.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != comp.Name {
		t.Fatalf("Get Name = %s, want %s", got.Name, comp.Name)
	}

	// older update ignored
	older := domain.UpdateComponentOptions{
		Name:            "older",
		VisualisationID: "table",
		Queries:         []domain.Query{{Name: "main"}},
		Properties:      map[domain.PropertyKey]domain.PropertyValue{},
	}
	if err := svc.Update(ctx, comp.ID, older, comp.UpdatedAt.Add(-time.Minute)); err != nil {
		t.Fatalf("Update older: %v", err)
	}
	got, _ = svc.Get(ctx, comp.ID)
	if got.Name != "comp" {
		t.Fatalf("after older update Name = %s, want %s", got.Name, "comp")
	}

	// newer update applied
	newer := domain.UpdateComponentOptions{
		Name:            "newer",
		VisualisationID: "table",
		Queries:         []domain.Query{{Name: "main"}},
		Properties:      map[domain.PropertyKey]domain.PropertyValue{},
	}
	if err := svc.Update(ctx, comp.ID, newer, comp.UpdatedAt.Add(time.Minute)); err != nil {
		t.Fatalf("Update newer: %v", err)
	}
	got, _ = svc.Get(ctx, comp.ID)
	if got.Name != "newer" {
		t.Fatalf("after newer update Name = %s, want %s", got.Name, "newer")
	}

	// delete then get => not found
	if err := svc.Delete(ctx, comp.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := svc.Get(ctx, comp.ID); err != service.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

// helpers
func newComponentService(t *testing.T) *service.ComponentService {
	t.Helper()
	db := openServiceDB(t)
	return service.NewComponentService(repository.NewComponentRepository(db))
}

func openServiceDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:svc-comp?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := storage.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}
