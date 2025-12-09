package domain_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/smilu97/refana/internal/pkg/domain"
)

// 사양서 기반 엔터티 스키마 테스트 (TDD).

func TestQueryShape(t *testing.T) {
	qt := reflect.TypeOf(domain.Query{})
	expectType(t, qt, "Name", typeOf(domain.Name("")))
	expectType(t, qt, "DataSourceID", typeOf(domain.DataSourceID("")))
	expectType(t, qt, "Properties", typeOf(map[domain.PropertyKey]domain.PropertyValue{}))
}

func TestComponentShape(t *testing.T) {
	ct := reflect.TypeOf(domain.Component{})
	expectType(t, ct, "ID", typeOf(domain.ComponentID("")))
	expectType(t, ct, "VisualisationID", typeOf(domain.VisualisationID("")))
	expectType(t, ct, "Query", typeOf(domain.Query{}))
	expectType(t, ct, "Name", typeOf(domain.Name("")))
	expectType(t, ct, "Coordination", typeOf(domain.Coordination{}))
	expectType(t, ct, "Properties", typeOf(map[domain.PropertyKey]domain.PropertyValue{}))
	expectType(t, ct, "UpdatedAt", typeOf(time.Time{}))
}

func TestCreateComponentOptionsShape(t *testing.T) {
	ot := reflect.TypeOf(domain.CreateComponentOptions{})
	expectType(t, ot, "VisualisationID", typeOf(domain.VisualisationID("")))
	expectType(t, ot, "Queries", typeOf([]domain.Query{}))
	expectType(t, ot, "Name", typeOf(domain.Name("")))
	expectType(t, ot, "Coordination", typeOf(domain.Coordination{}))
	expectType(t, ot, "Properties", typeOf(map[domain.PropertyKey]domain.PropertyValue{}))
}

func TestUpdateComponentOptionsShape(t *testing.T) {
	ot := reflect.TypeOf(domain.UpdateComponentOptions{})
	expectType(t, ot, "VisualisationID", typeOf(domain.VisualisationID("")))
	expectType(t, ot, "Queries", typeOf([]domain.Query{}))
	expectType(t, ot, "Name", typeOf(domain.Name("")))
	expectType(t, ot, "Coordination", typeOf(domain.Coordination{}))
	expectType(t, ot, "Properties", typeOf(map[domain.PropertyKey]domain.PropertyValue{}))
}

func TestColumnDataShape(t *testing.T) {
	ct := reflect.TypeOf(domain.ColumnData{})
	expectType(t, ct, "Name", typeOf(domain.Name("")))
	expectType(t, ct, "Type", typeOf(domain.PropertyType("")))
	expectType(t, ct, "Values", typeOf([]domain.PropertyValue{}))
}

func TestTableDataShape(t *testing.T) {
	tt := reflect.TypeOf(domain.TableData{})
	expectType(t, tt, "Columns", typeOf([]domain.ColumnData{}))
}

func TestDataSourceShape(t *testing.T) {
	dt := reflect.TypeOf(domain.DataSource{})
	expectType(t, dt, "ID", typeOf(domain.DataSourceID("")))
	expectType(t, dt, "ClassID", typeOf(domain.DataSourceClassID("")))
	expectType(t, dt, "Name", typeOf(domain.Name("")))
	expectType(t, dt, "Alias", typeOf(domain.Alias("")))
	expectType(t, dt, "Properties", typeOf(map[domain.PropertyKey]domain.PropertyValue{}))
}

func TestCreateDataSourceOptionsShape(t *testing.T) {
	ot := reflect.TypeOf(domain.CreateDataSourceOptions{})
	expectType(t, ot, "ClassID", typeOf(domain.DataSourceClassID("")))
	expectType(t, ot, "Name", typeOf(domain.Name("")))
	expectType(t, ot, "Alias", typeOf(domain.Alias("")))
	expectType(t, ot, "Properties", typeOf(map[domain.PropertyKey]domain.PropertyValue{}))
}

func TestUpdateDataSourceOptionsShape(t *testing.T) {
	ot := reflect.TypeOf(domain.UpdateDataSourceOptions{})
	expectType(t, ot, "ClassID", typeOf(domain.DataSourceClassID("")))
	expectType(t, ot, "Name", typeOf(domain.Name("")))
	expectType(t, ot, "Alias", typeOf(domain.Alias("")))
	expectType(t, ot, "Properties", typeOf(map[domain.PropertyKey]domain.PropertyValue{}))
}

func TestDataSourceClassShape(t *testing.T) {
	ct := reflect.TypeOf(domain.DataSourceClass{})
	expectType(t, ct, "ID", typeOf(domain.DataSourceClassID("")))
	expectType(t, ct, "Name", typeOf(domain.Name("")))
	expectType(t, ct, "PropertyDescriptors", typeOf([]domain.PropertyDescriptor{}))
}

// Helpers
func expectType(t *testing.T, typ reflect.Type, field string, want reflect.Type) {
	t.Helper()
	f, ok := typ.FieldByName(field)
	if !ok {
		t.Fatalf("%s missing field %s", typ.Name(), field)
	}
	if f.Type != want {
		t.Fatalf("%s.%s type = %v, want %v", typ.Name(), field, f.Type, want)
	}
}

func typeOf[T any](v T) reflect.Type { return reflect.TypeOf(v) }
