package middleware_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TechBowl-japan/go-stations/handler/middleware"
)

func TestAccessLog(t *testing.T) {
	wantPath := "/path"
	r := httptest.NewRequest(http.MethodGet, wantPath, nil)
	w := httptest.NewRecorder()
	var buf bytes.Buffer
	m := middleware.NewAccessLogMiddlewareWithWriter(&buf)

	wantOS := "macOS"
	ctx := context.WithValue(r.Context(), middleware.UAContextKeyOS, wantOS)
	r = r.WithContext(ctx)
	h := m.ServeNext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	h.ServeHTTP(w, r)

	var al middleware.AccessLog
	json.NewDecoder(&buf).Decode(&al)

	gotPath := al.Path
	if gotPath != wantPath {
		t.Errorf("正しいPath情報が記録されていません, got = %v, want = %v", gotPath, wantPath)
	}
	gotOS := al.OS
	if gotOS != wantOS {
		t.Errorf("正しいOS情報が記録されていません, got = %v, want = %v", gotOS, wantOS)

	}
}

func TestAccessLogWithoutOSInfo(t *testing.T) {
	wantPath := "/path"
	r := httptest.NewRequest(http.MethodGet, wantPath, nil)
	w := httptest.NewRecorder()
	var buf bytes.Buffer
	m := middleware.NewAccessLogMiddlewareWithWriter(&buf)

	wantOS := ""
	h := m.ServeNext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	h.ServeHTTP(w, r)

	var al middleware.AccessLog
	json.NewDecoder(&buf).Decode(&al)

	gotPath := al.Path
	if gotPath != wantPath {
		t.Errorf("正しいPath情報が記録されていません, got = %v, want = %v", gotPath, wantPath)
	}

	// NOTE: contextにOS情報が設定されていない場合は、OS情報は空文字で出力される
	gotOS := al.OS
	if gotOS != wantOS {
		t.Errorf("正しいOS情報が記録されていません, got = %v, want = %v", gotOS, wantOS)

	}
}
