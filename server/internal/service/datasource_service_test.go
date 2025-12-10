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

func TestDataSourceService_CreateValidate(t *testing.T) {
	svc := newDataSourceService(t)
	ctx := context.Background()

	_, err := svc.Create(ctx, domain.CreateDataSourceOptions{
		Name:   "",
		ClassID: "",
	})
	if err != service.ErrBadRequest {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestDataSourceService_CRUD_LastWriteWins(t *testing.T) {
	svc := newDataSourceService(t)
	ctx := context.Background()

	ds, err := svc.Create(ctx, domain.CreateDataSourceOptions{
		Name:    "ds",
		ClassID: "postgres",
		Alias:   "a",
		Properties: map[domain.PropertyKey]domain.PropertyValue{
			"host": "localhost",
		},
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := svc.Get(ctx, ds.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != ds.Name {
		t.Fatalf("Get Name = %s, want %s", got.Name, ds.Name)
	}

	// older update ignored
	older := domain.UpdateDataSourceOptions{
		Name:    "older",
		ClassID: "postgres",
		Alias:   "a",
	}
	if err := svc.Update(ctx, ds.ID, older, ds.UpdatedAt.Add(-time.Minute)); err != nil {
		t.Fatalf("Update older: %v", err)
	}
	got, _ = svc.Get(ctx, ds.ID)
	if got.Name != "ds" {
		t.Fatalf("after older update Name = %s, want %s", got.Name, "ds")
	}

	// newer update applied
	newer := domain.UpdateDataSourceOptions{
		Name:    "newer",
		ClassID: "postgres",
		Alias:   "b",
	}
	if err := svc.Update(ctx, ds.ID, newer, ds.UpdatedAt.Add(time.Minute)); err != nil {
		t.Fatalf("Update newer: %v", err)
	}
	got, _ = svc.Get(ctx, ds.ID)
	if got.Name != "newer" || got.Alias != "b" {
		t.Fatalf("after newer update = (%s,%s), want (newer,b)", got.Name, got.Alias)
	}

	// delete then get => not found
	if err := svc.Delete(ctx, ds.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := svc.Get(ctx, ds.ID); err != service.ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

// helpers
func newDataSourceService(t *testing.T) *service.DataSourceService {
	t.Helper()
	db := openDSServiceDB(t)
	return service.NewDataSourceService(repository.NewDataSourceRepository(db))
}

func openDSServiceDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:svc-ds?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := storage.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}
