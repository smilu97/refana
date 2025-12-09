package domain_test

import (
	"reflect"
	"testing"

	"github.com/smilu97/refana/internal/pkg/domain"
)

// 사양서 기반 도메인 타입 유효성 검증.
// NOTE: Go는 타입 선언에 태그를 붙일 수 없으므로 값 제약 검증은 런타임 validator에서 다룰 예정.

func TestRectFields(t *testing.T) {
	rt := reflect.TypeOf(domain.Rect{})
	checkField := func(name string) {
		f, ok := rt.FieldByName(name)
		if !ok {
			t.Fatalf("Rect missing field %s", name)
		}
		if f.Type.Kind() != reflect.Uint32 {
			t.Fatalf("Rect.%s type = %s, want uint32", name, f.Type.Kind().String())
		}
	}
	checkField("Left")
	checkField("Top")
	checkField("Width")
	checkField("Height")
}

func TestCoordinationFields(t *testing.T) {
	ct := reflect.TypeOf(domain.Coordination{})

	_, ok := ct.FieldByName("Rect")
	if !ok {
		t.Fatalf("Coordination missing embedded Rect")
	}

	z, ok := ct.FieldByName("ZIndex")
	if !ok {
		t.Fatalf("Coordination missing ZIndex")
	}
	if z.Type.Kind() != reflect.Uint32 {
		t.Fatalf("Coordination.ZIndex type = %s, want uint32", z.Type.Kind().String())
	}
}

func TestPropertyTypeConstants(t *testing.T) {
	if domain.PropertyTypeString != "string" {
		t.Fatalf("PropertyTypeString = %q, want %q", domain.PropertyTypeString, "string")
	}
	if domain.PropertyTypeNumber != "number" {
		t.Fatalf("PropertyTypeNumber = %q, want %q", domain.PropertyTypeNumber, "number")
	}
}

func TestPropertyDescriptorShape(t *testing.T) {
	pt := reflect.TypeOf(domain.PropertyDescriptor{})

	expectField := func(name string, kind reflect.Kind) {
		f, ok := pt.FieldByName(name)
		if !ok {
			t.Fatalf("PropertyDescriptor missing field %s", name)
		}
		if f.Type.Kind() != kind {
			t.Fatalf("PropertyDescriptor.%s kind = %s, want %s", name, f.Type.Kind(), kind)
		}
	}

	expectField("Key", reflect.String)
	expectField("Name", reflect.String)
	expectField("Type", reflect.String)
	expectField("Category", reflect.String)
	expectField("Order", reflect.Uint32)
	expectField("IsRequired", reflect.Bool)
	expectField("IsSecret", reflect.Bool)

	f, ok := pt.FieldByName("Candidates")
	if !ok {
		t.Fatalf("PropertyDescriptor missing Candidates")
	}
	if f.Type.Kind() != reflect.Slice || f.Type.Elem().Kind() != reflect.String {
		t.Fatalf("PropertyDescriptor.Candidates type = %s, want []string", f.Type.Kind().String())
	}
}
