package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TechBowl-japan/go-stations/handler/middleware"
)

func TestRecovery(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	m := middleware.NewRecoveryMiddleware()
	h := m.ServeNext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test")
	}))
	h.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("期待していない HTTP status code です, got = %d, want = %d", w.Code, http.StatusInternalServerError)
	}
}
