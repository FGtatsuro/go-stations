package middleware

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// accessLog は、日時、処理時間等のアクセスログを表す構造体である。
type accessLog struct {
	Timestamp time.Time
	Latency   int64
	Path      string
	OS        string
	Status    int
}

type accessLogMiddleware struct {
	w io.Writer
}

// NewAccessLogMiddleware は、 書き込み先として標準出力を指定したミドルウェアを返す。
func NewAccessLogMiddleware() *accessLogMiddleware {
	return &accessLogMiddleware{
		w: os.Stdout,
	}
}

// statusResponseWriter は、HTTPステータスをログに記録するために、デフォルトの [net/http.ResponseWriter] を拡張した構造体である
//
// Ref: https://tutuz-tech.hatenablog.com/entry/2020/06/14/191416
type statusResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// ServeNext は、 h の前後で取得した情報を元に、 アクセスログを記録する。
func (m *accessLogMiddleware) ServeNext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		sw := &statusResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		now := time.Now()
		h.ServeHTTP(sw, r)

		// NOTE: OS情報がContextに記録されている事を前提とする。
		os, ok := r.Context().Value(UAContextKeyOS).(string)
		if !ok {
			log.Printf("AccessLogMiddleware: os can not be fetched\n")
		}

		al := accessLog{
			Timestamp: now,
			Latency:   time.Since(now).Milliseconds(),
			Path:      r.URL.Path,
			OS:        os,
			Status:    sw.status,
		}
		if err := json.NewEncoder(m.w).Encode(al); err != nil {
			log.Printf("AccessLogMiddleware: could not write access log, err =%v\n", err)
			return
		}
	}
	return http.HandlerFunc(fn)
}
