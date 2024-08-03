package middleware

import "net/http"

// HTTPMiddleware は、HTTP層での共通処理を実装するための構造体である。
//
// ServeNext で h の ServeHTTP を呼ぶ事は利用側の責務とする。
//
// Ref: https://docs.google.com/presentation/d/1BEay5y7U80Ha1tetZvTPZDZZ9QW0-Iv459hyyRv07ME/edit#slide=id.g4fa6665ad1_0_517v
type HTTPMiddleware interface {
	ServeNext(h http.Handler) http.Handler
}

// With は、[HTTPMiddleware] を合成する関数である。
//
// h 及び hms は、引数で与えられた順とは逆順で評価される。
func With(h http.Handler, hms ...HTTPMiddleware) http.Handler {
	for _, hm := range hms {
		h = hm.ServeNext(h)
	}
	return h
}
