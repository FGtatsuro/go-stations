package middleware

import (
	"log"
	"net/http"
)

// Recovery は、[http.Handler] でpanicが発生した際のリカバリ処理を追加するHTTPミドルウェアである。
//
// リカバリ処理が実行された場合、ユーザにはstatus 500を返す。
// panic発生前に [http.ResponseWriter.WriteHeader] が呼ばれていた場合、statusは上書きされない。
// [http.ResponseWriter.WriteHeader] がpanic発生前に呼ばれない事を保証するのは利用側の責務とする。
func Recovery(h http.Handler) http.Handler {
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
