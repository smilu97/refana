package domain

import "github.com/rushysloth/go-tsid"

// Common value objects and descriptors used across the server.
// Validation tags follow the specification document.

// Name and Alias
type Name string
type Alias string

// Rect and Coordination
type Rect struct {
	Left   uint32
	Top    uint32
	Width  uint32
	Height uint32
}

type Coordination struct {
	Rect
	ZIndex uint32
}

// IDs
type GeneratedID struct{ tsid tsid.Tsid }

func NewGeneratedID(v int64) GeneratedID {
	return GeneratedID{tsid: *tsid.FromNumber(v)}
}

func (id GeneratedID) Int64() int64 {
	return id.tsid.ToNumber()
}

type ComponentID struct{ GeneratedID }

func NewComponentID(v int64) ComponentID {
	return ComponentID{GeneratedID: NewGeneratedID(v)}
}

type DataSourceID struct{ GeneratedID }

func NewDataSourceID(v int64) DataSourceID {
	return DataSourceID{GeneratedID: NewGeneratedID(v)}
}

type DesignatedID string
type VisualisationID DesignatedID
type DataSourceClassID DesignatedID

// Properties
type PropertyType string

const (
	PropertyTypeString PropertyType = "string"
	PropertyTypeNumber PropertyType = "number"
)

type PropertyKey string
type PropertyValue string

type PropertyDescriptor struct {
	Key        PropertyKey
	Name       Name
	Type       PropertyType
	Category   Name
	Order      uint32
	IsRequired bool
	IsSecret   bool
	Candidates []PropertyValue
}
