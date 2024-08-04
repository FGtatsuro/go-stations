package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

// Deprecated: [NewHandler] を使用して下さい。
func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	mux.Handle("/healthz", handler.NewHealthzHandler())

	// NOTE: RecoveryMiddleware を最後の引数とすることで、一番始めに評価される。
	// これにより、後述の処理でのpanicは全て RecoveryMiddleware で処理される。
	mux.Handle("/todos",
		middleware.With(
			handler.NewTODOHandler(service.NewTODOService(todoDB)),
			middleware.NewAccessLogMiddleware(),
			middleware.NewUserAgentRecordMiddleware(),
			middleware.NewRecoveryMiddleware(),
		),
	)
	mux.Handle("/do-panic",
		middleware.With(
			handler.NewPanicHandler(),
			middleware.NewAccessLogMiddleware(),
			middleware.NewUserAgentRecordMiddleware(),
			middleware.NewRecoveryMiddleware(),
		),
	)

	return mux
}

// NewHandler は、ルーティングを設定したHTTPハンドラを返す。
func NewHandler(todoDB *sql.DB) http.Handler {
	return newHandler(todoDB,
		middleware.NewAccessLogMiddleware(),
		middleware.NewUserAgentRecordMiddleware(),
		middleware.NewRecoveryMiddleware(),
	)
}

func newHandler(todoDB *sql.DB, ms ...middleware.HTTPMiddleware) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/healthz", handler.NewHealthzHandler())

	// NOTE: 初級編の課題のテストが /todos に依存しているため、下記のパスは残したままとする
	mux.Handle("/todos", handler.NewTODOHandler(service.NewTODOService(todoDB)))

	// NOTE: 認証の範囲を限定する(e.g. ヘルスチェックには認証を設定したくない)ため、/api 以下のパスにのみ認証を設定する。
	//
	// Ref: https://forum.golangbridge.org/t/is-it-possible-to-combine-http-servemux/7495/4
	apiMux := http.NewServeMux()
	apiMux.Handle("/todos", handler.NewTODOHandler(service.NewTODOService(todoDB)))
	apiMux.Handle("/do-panic", handler.NewPanicHandler())
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	// *http.ServeMux は http.Handler interfaceを満たすため、他のハンドラ同様ミドルウェアが適用できる。
	//
	// Ref: https://blog.afoolishmanifesto.com/posts/nesting-middleware-in-golang/
	return middleware.With(
		mux,
		ms...,
	)
}
