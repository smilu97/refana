# Server TODOs
각 항목은 하나의 git commit 에 대응한다.

1. `go.mod`/`go.sum` 초기화 및 기본 의존성 추가 (gorm, sqlite, gin, tsid, validator)
2. 서버 부트스트랩 스켈레톤: `cmd/server/main.go` 에 HTTP 서버, gin 라우터, healthz, 기본 DI 구성
3. 공통 도메인 타입 정의: `internal/pkg/domain`에 ID/Name/Alias/Coordination/Property 타입, 검증 태그 포함
4. 엔티티 스키마 모델링: Component, DataSource, DataSourceClass, Query, TableData 등 struct 정의 및 JSON 태깅
5. SQLite3 GORM 설정 및 마이그레이션 파이프라인 추가 (auto-migrate 또는 명시적 migration)
6. Repository 계층 구현: Component, DataSource, DataSourceClass CRUD + UpdatedAt 기반 last-write-wins 정책 반영
7. Service 계층 구현: 비즈니스 로직/검증/에러 매핑 (BadRequest/NotFound/ErrorResponse)
8. HTTP 핸들러 - Components: `GET /components`, `GET /components/:id`, `GET /components/:id/data`, `POST`, `PATCH`, `DELETE`
9. HTTP 핸들러 - DataSources: `GET /data-sources`, `GET /data-sources/:id`, `POST`, `PATCH`, `DELETE`
10. HTTP 핸들러 - DataSourceClasses: `GET /data-source-classes`, `GET /data-source-classes/:id`
11. `TableData` 빌더: DataSource/Query 프로퍼티를 받아 테이블 형태로 변환하는 인터페이스/스텁 구현
12. PostgreSQL DataSourceClass 구현: 연결 프로퍼티 검증, 쿼리 실행, 결과를 `TableData`로 변환
13. JSON import/export 구현: Component/DataSource bulk import/export 엔드포인트 혹은 파일 기반 CLI 추가
14. 동시 수정 대응: UpdatedAt 비교 기반 last-write-wins 적용 및 낙관적 잠금 테스트 추가
15. 단위/통합 테스트: 서비스, 핸들러, PostgreSQL 스텁 테스트 및 sqlite 인메모리 사용
16. 구성/런 명령: `configs/` 예시, `.env` 로드, `Makefile`에 `make run`/`make test` 추가
17. Dockerfile/DevContainer: Go 빌드, sqlite3 도구, optional frontend 정적파일 서빙 자리 마련
