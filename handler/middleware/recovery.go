package middleware

import (
	"log"
	"net/http"
)

// RecoveryMiddleware は、panicが発生した際のリカバリ処理を追加するHTTPミドルウェアである。
type RecoveryMiddleware struct{}

// ServeNext は、h でpanicが発生した際にリカバリ処理を行い、ユーザにstatus 500を返す。
//
// panic発生前に [net/http.ResponseWriter] の WriteHeader が呼ばれていた場合、statusは上書きされない。
// [net/http.ResponseWriter] の WriteHeader がpanic発生前に呼ばれない事を保証するのは利用側の責務とする。
func (m *RecoveryMiddleware) ServeNext(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("recovery: panic =%v\n", p)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
