package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

type uaContextKey string

const (
	UAContextKeyOS = uaContextKey("os")
)

type userAgentRecordMiddleware struct{}

// NewUserAgentRecordMiddleware は、 [context.Context] にリクエストのOS情報をセットするミドルウェアを返す。
func NewUserAgentRecordMiddleware() *userAgentRecordMiddleware {
	return &userAgentRecordMiddleware{}
}

// ServeNext は、[UAContextKeyOS] をキーとしてリクエストのOS情報を保存する。
func (m *userAgentRecordMiddleware) ServeNext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), UAContextKeyOS, ua.OS)

		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
