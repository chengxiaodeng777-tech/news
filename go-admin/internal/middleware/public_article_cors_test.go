package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPublicArticleCORS(t *testing.T) {
	r := gin.New()
	r.Use(PublicArticleCORS())
	r.NoRoute(func(c *gin.Context) { c.Status(http.StatusNotFound) })
	for _, tc := range []struct{ method, path, origin string }{
		{http.MethodGet, "/api/articles", "*"},
		{http.MethodGet, "/api/articles/1", "*"},
		{http.MethodOptions, "/api/articles/1", "*"},
		{http.MethodGet, "/api/system/articles", ""},
		{http.MethodGet, "/api/articles-other", ""},
		{http.MethodPost, "/api/articles", ""},
	} {
		t.Run(tc.method+tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("Origin", "http://localhost:8123")
			r.ServeHTTP(w, req)
			if got := w.Header().Get("Access-Control-Allow-Origin"); got != tc.origin {
				t.Fatalf("origin = %q, want %q", got, tc.origin)
			}
			if w.Header().Get("Access-Control-Allow-Credentials") != "" {
				t.Fatal("anonymous reads must not enable credentials")
			}
		})
	}
}
