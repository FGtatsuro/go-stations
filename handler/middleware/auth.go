package middleware

import (
	"net/http"

	"github.com/TechBowl-japan/go-stations/pkg/basicauth"
)

type basicAuthMiddleware struct {
	bai basicauth.BasicAuthInfo
}

// NewBasicAuthMiddleware は、Basic認証によるアクセス制限を行うミドルウェアを返す。
func NewBasicAuthMiddleware(bai basicauth.BasicAuthInfo) *basicAuthMiddleware {
	return &basicAuthMiddleware{
		bai: bai,
	}
}

// ServeNext は、Basic認証によるアクセス制限を行う。
func (m *basicAuthMiddleware) ServeNext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if err := m.bai.Authenticate(r); err != nil {
			m.bai.Challenge(w)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
