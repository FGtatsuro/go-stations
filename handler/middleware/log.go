package middleware

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// accessLog は、日時、処理時間等のアクセス情報を表す構造体である。
type accessLog struct {
	Timestamp time.Time
	Latency   int64
	Path      string
	OS        string
}

// accessLogMiddleware は、 [accessLog] を w に書き込むHTTPミドルウェアである。
type accessLogMiddleware struct {
	w io.Writer
}

// NewAccessLogMiddleware は、 書き込み先として標準出力を指定した [AccessLogMiddleware] を返す。
func NewAccessLogMiddleware() *accessLogMiddleware {
	return &accessLogMiddleware{
		w: os.Stdout,
	}
}

// ServeNext は、 h の前後で取得した情報を元に、 [accessLog] を記録する。
func (m *accessLogMiddleware) ServeNext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		before := time.Now()
		h.ServeHTTP(w, r)
		after := time.Now()

		// NOTE: UserAgentRecordMiddleware によりOS情報がContextに記録されている事を前提とする。
		os, ok := r.Context().Value(UAContextKeyOS).(string)
		if !ok {
			log.Printf("AccessLogMiddleware: os can not be fetched\n")

		}

		al := accessLog{
			Timestamp: before,
			Latency:   after.Sub(before).Milliseconds(),
			Path:      r.URL.Path,
			OS:        os,
		}
		if err := json.NewEncoder(m.w).Encode(al); err != nil {
			log.Printf("AccessLogMiddleware: could not write access log, err =%v\n", err)
			return
		}
	}
	return http.HandlerFunc(fn)
}
