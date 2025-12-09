Retool + Grafana 특징을 섞어보고자 함.

Grafana는 개발자 친화적이며 빠른 읽기 및 다양한 시각화 패널들을 제공하는 점을 강점으로 하고 있음. 추가로, Alertmanager 를 통해 다양한 채널로 비정상 상황을 전파 가능. 인프라 관리나 고수준 Observability 를 달성하는 것을 목적으로 함.

Retool은 비개발자가 비즈니스 흐름에 관여할 수 있도록 도와주는 도구 생성 및 관리를 저비용으로 달성할 수 있게 함. 관계자용 인터페이스를 만드는 것을 목적으로 함.

서로 목적은 다르지만 설계 수준에서는 Datasource, Page, Component 등 기반을 지탱하는 요소들이 비슷해 보임. (디테일한 기능들과 관련해서는 겹치는 부분이 크지 않아 보임)

# 1차 목표

아래 사항들을 모두 만족하면 달성인 것으로 함.

1. Navigation
	1. 메인 페이지: 모든 컴포넌트들에 대한 렌더링
	2. 컴포넌트 상세 페이지: 단일 컴포넌트 렌더링 및 프로퍼티 수정
	3. 데이터소스 페이지: 모든 데이터소스들에 대해 리스트 업
	4. 데이터소스 상세 페이지: 데이터소스 수정
2. Component, Datasource는 JSON 으로 import, export 가능 해야함.
3. 여러 명이 Component, Datasource 수정 시 가장 늦게 Commit 된 내용이 살아남음.
4. Component, Datasource 내용은 Sqlite3 에 저장됨.
5. Visualisations: Table, Form, Text 지원
6. DatasourceClasses: PostgreSQL 지원

# 기초 설계

- Web SPA 앱을 빌드해서 Golang 서버에서 API, SPA static files 를 한꺼번에 서빙함
	- Solid.js 기반 SPA
	- 아래 설정으로 개발 환경에서 API 프록시 할 수 있다고 함
```typescript
// app.config.ts
import { defineConfig } from "@solidjs/start/config";

export default defineConfig({
  // SolidStart 쪽 설정
  start: {
    ssr: false, // SPA
    server: {
      preset: "static", // 빌드 결과를 정적 파일로만 쓰겠다는 의미 (Go가 서빙)
      devProxy: {
        "/api": {
          target: "http://localhost:8080", // Go 백엔드 (dev용)
          changeOrigin: true,
        },
      },
    },
  },
});
```
- 영속성은 sqlite3 기반으로 구현
	- gORM 사용: PostgreSQL 기반으로 확장하는 것을 고려할 것
- DataSource가 서술하는 모든 쿼리는 Golang 에서 수행

## 디렉토리 구조

- frontend/: solid.js 어플리케이션
	- package.json
	- ...
- server/
	- cmd/
	- internal/
		- app/
		- pkg/
	- configs/
	- go.mod
- Makefile
- Dockerfile

## Entities

### Component

Visualisation에 필요한 데이터들을 어떻게 쿼리할 것인지에 대해 서술함
어떤 Visualisation으로 시각화하여 화면 어디에 위치할 것인지 서술함

#### API

prefix: /api

- GET /components
	- ResponseBody
		- 200 OK: `[]Component`
		- 500 Internal Server Error: `ErrorResponse`
- GET /components/:componentId
	- ResponseBody
		- 200 OK: `Component`
		- 404 Not Found: `NotFoundResponse`
		- 500 Internal Server Error: `ErrorResponse`
- GET /components/:componentId/data
	- ResponseBody
		- 200 OK: `TableData`
		- 404 Not Found: `NotFoundResponse`
		- 500 Internal Server Error: `ErrorResponse`
- POST /components
	- RequestBody: `CreateComponentOptions`
	- ResponseBody:
		- 201 Created: `Component`
		- 400 Bad Request: `BadRequestResponse`
		- 500 Internal Server Error: `ErrorResponse`
- PATCH /components/:componentId
	- RequestBody: `UpdateComponentOptions`
	- ResponseBody:
		- 200 Created: `Component`
		- 400 Bad Request: `BadRequestResponse`
		- 404 Not Found: `NotFoundResponse`
		- 500 Internal Server Error: `ErrorResponse`
- DELETE /components/:componentId
	- ResponseBody:
		- 200 OK
		- 404 Not Found: `NotFoundResponse`
		- 500 Internal Server Error: `ErrorResponse`

#### Schema

```go
type Query struct {
	Name         Name
	DataSourceID DataSourceID
	Properties   map[PropertyKey]PropertyValue
}

type Component struct {
	ID              ComponentID
	VisualisationID VisualisationID
	Query           Query
	Name            Name
	Coordination    Coordination
	Properties      map[PropertyKey]PropertyValue
	UpdatedAt       time.Time
}

type CreateComponentOptions struct {
	VisualisationID VisualisationID
	Queries         []Queries
	Name            Name
	Coordination    Coordination
	Properties      map[PropertyKey]PropertyValue
}

type UpdateComponentOptions struct {
	VisualisationID VisualisationID
	Queries         []Queries
	Name            Name
	Coordination    Coordination
	Properties      map[PropertyKey]PropertyValue
}

type ColumnData struct {
	Name   Name
	Type   PropertyType
	Values []PropertyValues
}

type TableData struct {
	Columns []ColumnData
}
```

### Visualisation

`TableData` 가 있을 때, 어떤 식의 UI 를 렌더링 할 것인지 서술함.
`Visualisation` 은 오로지 프론트의 책임
백엔드에서 말하는 어떤 VisualisationID 가 실제로는 프론트 앱에서 지원하지 않을 수도 있음...

### DataSource

어떤 DataSourceClass를 이용해서 쿼리할 것인지 서술
이 때 어떤 공통 Properties 를 가져야 하는지 서술. 예시) host, port, credential, ...

#### API

- GET /data-sources
	- ResponseBody
		- 200 OK: `[]DataSource`
		- 500 Internal Server Error: `ErrorResponse`
- GET /data-sources/:datasourceId
	- ResponseBody
		- 200 OK: `DataSource`
		- 404 Not Found: `NotFoundResponse`
		- 500 Internal Server Error: `ErrorResponse`
- POST /data-sources
	- RequestBody: `CreateDataSourceOptions`
	- ResponseBody
		- 201 Created: `DataSource`
		- 400 Bad Request: `BadRequestResponse`
		- 500 Internal Server Error: `ErrorResponse`
- PATCH /data-sources/:datasourceId
	- RequestBody: `UpdateDataSourceOptions`
	- ResponseBody
		- 200 OK: `DataSource`
		- 400 Bad Request: `BadRequestResponse`
		- 500 Internal Server Error: `ErrorResponse`
- DELETE /data-sources/:datasourceId
	- ResponseBody
		- 200 OK:
		- 404 Not Found: `NotFoundResponse`
		- 500 Internal Server Error: `ErrorResponse`
#### Schema

```go
type DataSource struct {
	ID         DataSourceID
	ClassID    DataSourceClassID
	Name       Name
	Alias      Alias
	Properties map[PropertyKey]PropertyValue
}

type CreateDataSourceOptions struct {
	ClassID    DataSourceClassID
	Name       Name
	Alias      Alias
	Properties map[PropertyKey]PropertyValue
}

type UpdateDataSourceOptions struct {
	ClassID    DataSourceClassID
	Name       Name
	Alias      Alias
	Properties map[PropertyKey]PropertyValue
}
```


### DataSourceClass

Backend 기준 빌드타임에 모든 것이 결정됨. Frontend 에서는 API 로 받아서 활용.

`DataSource` 의 `Properties`, `Query` 의 `Properties` 를 받아서 어떻게 `TableData` 를 가져올 것인지를 책임

미리 관련된 백엔드 코드가 준비되어 있어야 함.

추후에 플러그인 형태로 확장할 수 있는 방법을 어떻게 제공할 것인지 고려해야 함.

#### API

- GET /data-source-classes
	- ResponseBody
		- 200 OK: `[]DataSourceClass`
		- 500 Internal Server Error: `ErrorResponse`
- GET /data-source-classes/:dataSourceClassId
	- ResponseBody
		- 200 OK: `DataSourceClass`
		- 404 Not Found: `NotFoundResponse`
		- 500 Internal Server Error: `ErrorResponse`
#### Schema

```go
type DataSourceClass struct {
	ID                  DataSourceClassID
	Name                Name
	PropertyDescriptors []PropertyDescriptor
}
```

### Etc

#### Schemas

```go
// datatypes.go

type Name string `validate:"max=256"`
type Alias string `validate:"max=64"`

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
```

```go
// id.go

// Machine generated ID
type GeneratedID tsid.Tsid
type ComponentID GeneratedID
type DataSourceID GeneratedID

// Human designated ID
type DesignatedID string `validate:"max=64"`
// etc) table, form, time-series, select, text, button
type VisualisationID DesignatedID `validate:"example=table"`
// etc) mysql, postgresql, prometheus
type DataSourceClassID DesignatedID `validate:"example=mysql"`
```

```go
// property.go

type PropertyType string
const (
	PropertyTypeString PropertyType = "string"
	PropertyTypeNumber PropertyType = "number"
)

type PropertyKey string
type PropertyValue string

type PropertyDescriptor struct {
	Key             PropertyKey
	Name            Name
	Type            PropertyType
	// Category 가 같은 것 끼리 뭉쳐야 함
	Category        Name
	// UI 상 배치에 관여
	Order           uint32
	IsRequired      bool
	IsSecret        bool
	// select type 등을 위해? static 해도 상관없을까
	Candidates      []PropertyValue
}
```

# TODO

- Page 구현
- Project 구현
- PropertyTypeSQL 구현
- Form Visualisation 구현
- Form 제출 후 관련 Component 자동 리프레시
- PostgreSQL 백엔드
- MySQL, Prometheus 데이터소스
