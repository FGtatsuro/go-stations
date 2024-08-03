package middleware

import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"
)

// UserAgentRecord は、[context.Context] にリクエストのOS情報をセットするHTTPミドルウェアである。
type UserAgentRecordMiddleware struct{}

type uaContextKey string

const (
	UAContextKeyOS = uaContextKey("os")
)

// ServeNext は、[UAContextKeyOS] をキーとしてリクエストのOS情報を保存する。
func (m *UserAgentRecordMiddleware) ServeNext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), UAContextKeyOS, ua.OS)

		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
