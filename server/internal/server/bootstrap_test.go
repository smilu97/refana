package server_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/smilu97/refana/internal/server"
)

// TDD: 서버 부트스트랩에서 제공해야 할 기본 라우터 동작을 먼저 정의한다.
func TestNewRouter_Healthz(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := server.NewRouter(context.Background(), server.Deps{})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("healthz status = %d, want %d", w.Code, http.StatusOK)
	}
	if got := w.Header().Get("Content-Type"); got != "application/json; charset=utf-8" {
		t.Fatalf("content-type = %q, want %q", got, "application/json; charset=utf-8")
	}
	const wantBody = `{"status":"ok"}`
	if body := w.Body.String(); body != wantBody {
		t.Fatalf("healthz body = %q, want %q", body, wantBody)
	}
}

func TestNewRouter_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := server.NewRouter(context.Background(), server.Deps{})

	req := httptest.NewRequest(http.MethodGet, "/not-found", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("not-found status = %d, want %d", w.Code, http.StatusNotFound)
	}
}
