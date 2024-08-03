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

// UserAgentRecord は、[context.Context] にリクエストのOS情報をセットするHTTPミドルウェアである。
//
// 利用側は、[UAContextKeyOS] をキーとしてリクエストのOS情報を取り出す。
func UserAgentRecord(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), UAContextKeyOS, ua.OS)

		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
