package domain

import "time"

// Query describes how to fetch data for a visualisation.
type Query struct {
	Name         Name                           `json:"name"`
	DataSourceID DataSourceID                   `json:"dataSourceId"`
	Properties   map[PropertyKey]PropertyValue  `json:"properties"`
}

// Component binds a visualisation to its data and layout.
type Component struct {
	ID              ComponentID                  `json:"id"`
	VisualisationID VisualisationID              `json:"visualisationId"`
	Query           Query                        `json:"query"`
	Name            Name                         `json:"name"`
	Coordination    Coordination                 `json:"coordination"`
	Properties      map[PropertyKey]PropertyValue `json:"properties"`
	UpdatedAt       time.Time                    `json:"updatedAt"`
}

type CreateComponentOptions struct {
	VisualisationID VisualisationID              `json:"visualisationId"`
	Queries         []Query                      `json:"queries"`
	Name            Name                         `json:"name"`
	Coordination    Coordination                 `json:"coordination"`
	Properties      map[PropertyKey]PropertyValue `json:"properties"`
}

type UpdateComponentOptions struct {
	VisualisationID VisualisationID              `json:"visualisationId"`
	Queries         []Query                      `json:"queries"`
	Name            Name                         `json:"name"`
	Coordination    Coordination                 `json:"coordination"`
	Properties      map[PropertyKey]PropertyValue `json:"properties"`
}

// Tabular data returned to the frontend.
type ColumnData struct {
	Name   Name           `json:"name"`
	Type   PropertyType   `json:"type"`
	Values []PropertyValue `json:"values"`
}

type TableData struct {
	Columns []ColumnData `json:"columns"`
}

// DataSource describes a configured backend data provider.
type DataSource struct {
	ID         DataSourceID                  `json:"id"`
	ClassID    DataSourceClassID             `json:"classId"`
	Name       Name                          `json:"name"`
	Alias      Alias                         `json:"alias"`
	Properties map[PropertyKey]PropertyValue `json:"properties"`
	UpdatedAt  time.Time                     `json:"updatedAt"`
}

type CreateDataSourceOptions struct {
	ClassID    DataSourceClassID             `json:"classId"`
	Name       Name                          `json:"name"`
	Alias      Alias                         `json:"alias"`
	Properties map[PropertyKey]PropertyValue `json:"properties"`
}

type UpdateDataSourceOptions struct {
	ClassID    DataSourceClassID             `json:"classId"`
	Name       Name                          `json:"name"`
	Alias      Alias                         `json:"alias"`
	Properties map[PropertyKey]PropertyValue `json:"properties"`
}

// DataSourceClass defines how to turn properties and queries into TableData.
type DataSourceClass struct {
	ID                  DataSourceClassID    `json:"id"`
	Name                Name                 `json:"name"`
	PropertyDescriptors []PropertyDescriptor `json:"propertyDescriptors"`
}
