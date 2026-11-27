package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"base/ent"
	"base/internal/config"
	"base/internal/db"

	"go.uber.org/zap"
)

func TestArticleRouteAccess(t *testing.T) {
	r := New(context.Background(), config.Config{}, zap.NewNop(), Dependencies{
		Database: &db.Resource{Client: ent.NewClient()},
	})
	for _, tc := range []struct {
		method, path string
		status       int
	}{
		{http.MethodGet, "/api/articles?page=invalid", http.StatusBadRequest},
		{http.MethodGet, "/api/articles/not-an-id", http.StatusBadRequest},
		{http.MethodGet, "/api/articles/0", http.StatusBadRequest},
		{http.MethodGet, "/api/articles/-1", http.StatusBadRequest},
		{http.MethodGet, "/api/system/articles", http.StatusUnauthorized},
		{http.MethodPost, "/api/system/articles", http.StatusUnauthorized},
		{http.MethodPut, "/api/system/articles/1", http.StatusUnauthorized},
		{http.MethodDelete, "/api/system/articles/1", http.StatusUnauthorized},
		{http.MethodPost, "/api/articles", http.StatusNotFound},
	} {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(tc.method, tc.path, nil))
			if w.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", w.Code, tc.status, w.Body.String())
			}
		})
	}
}
