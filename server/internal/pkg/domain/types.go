package domain

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
type GeneratedID string
type ComponentID GeneratedID
type DataSourceID GeneratedID

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
