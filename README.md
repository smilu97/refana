# Refana

Retool 의 낮은 제작 비용과 Grafana 의 개발자 친화적 시각화 경험을 섞어, 비즈니스/운영 모두가 활용할 수 있는 경량 대시보드/툴 빌더를 목표로 합니다.

## 1차 목표
- 내비게이션: 메인(모든 컴포넌트), 컴포넌트 상세(단일 렌더링·프로퍼티 수정), 데이터소스 목록, 데이터소스 상세(수정)
- 컴포넌트·데이터소스 JSON import/export 지원
- 동시 수정 시 가장 늦은 커밋이 승리 (last write wins)
- 내용은 SQLite3 에 영속
- Visualisations: Table, Form, Text 지원
- DataSourceClasses: PostgreSQL 지원

## 아키텍처 개요
- Web SPA: SolidStart 기반 SPA 빌드 후 Go 서버에서 API와 정적 파일을 함께 서빙
- 프록시 예시 (dev): `/api` → `http://localhost:8080`
- 백엔드: Go, gORM + SQLite3 (향후 PostgreSQL 확장 고려)
- 모든 쿼리는 백엔드에서 실행하며 `TableData` 형태로 프론트에 전달

## 디렉터리 구조
- `frontend/`: Solid.js 애플리케이션
- `server/`
  - `cmd/`, `internal/` (`app/`, `pkg/`), `configs/`
  - `go.mod`
- `Makefile`, `Dockerfile`

## 핵심 엔터티 및 API
### Component
- 시각화 대상 데이터 쿼리와 위치/프로퍼티를 정의
- 주요 API: `GET /components`, `GET /components/:id`, `GET /components/:id/data`, `POST /components`, `PATCH /components/:id`, `DELETE /components/:id`
- 스키마 요약:
```go
type Component struct {
  ID ComponentID; VisualisationID VisualisationID
  Query Query; Name Name; Coordination Coordination
  Properties map[PropertyKey]PropertyValue; UpdatedAt time.Time
}
```

### DataSource
- 어떤 DataSourceClass 로 쿼리할지와 프로퍼티 정의
- 주요 API: `GET /data-sources`, `GET /data-sources/:id`, `POST /data-sources`, `PATCH /data-sources/:id`, `DELETE /data-sources/:id`
```go
type DataSource struct {
  ID DataSourceID; ClassID DataSourceClassID
  Name Name; Alias Alias; Properties map[PropertyKey]PropertyValue
}
```

### DataSourceClass
- 빌드타임 정의. DataSource/Query 프로퍼티를 받아 `TableData` 를 만드는 책임
- 주요 API: `GET /data-source-classes`, `GET /data-source-classes/:id`
```go
type DataSourceClass struct {
  ID DataSourceClassID; Name Name; PropertyDescriptors []PropertyDescriptor
}
```

### 공통 타입
```go
type Name string; type Alias string
type VisualisationID string // 예: "table"
type DataSourceClassID string // 예: "mysql"
type Coordination struct { Rect struct{Left,Top,Width,Height uint32}; ZIndex uint32 }
type PropertyDescriptor struct {
  Key PropertyKey; Name Name; Type PropertyType
  Category Name; Order uint32; IsRequired bool; IsSecret bool
  Candidates []PropertyValue
}
```

## TODO
- Page, Project 구현
- PropertyTypeSQL, Form 시각화 구현 및 제출 후 관련 컴포넌트 자동 리프레시
- PostgreSQL 백엔드, MySQL/Prometheus 데이터소스
