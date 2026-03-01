package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSyncLimitFromRequest(t *testing.T) {
	t.Run("query limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/stocks/sync?limit=12", nil)
		limit, err := syncLimitFromRequest(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if limit != 12 {
			t.Fatalf("expected 12, got %d", limit)
		}
	})

	t.Run("body limit overrides query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/stocks/sync?limit=12", strings.NewReader(`{"limit":7}`))
		limit, err := syncLimitFromRequest(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if limit != 7 {
			t.Fatalf("expected 7, got %d", limit)
		}
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/stocks/sync", strings.NewReader(`{`))
		_, err := syncLimitFromRequest(req)
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}

func TestWithCORS(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := withCORS(next, "http://localhost:5173")

	req := httptest.NewRequest(http.MethodGet, "/stocks", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("unexpected allow origin: %s", rec.Header().Get("Access-Control-Allow-Origin"))
	}
}
