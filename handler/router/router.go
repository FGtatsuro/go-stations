package router

import (
	"database/sql"
	"net/http"

	"github.com/TechBowl-japan/go-stations/handler"
	"github.com/TechBowl-japan/go-stations/handler/middleware"
	"github.com/TechBowl-japan/go-stations/service"
)

func NewRouter(todoDB *sql.DB) *http.ServeMux {
	// register routes
	mux := http.NewServeMux()

	mux.Handle("/healthz", handler.NewHealthzHandler())

	// NOTE: RecoveryMiddleware を最後の引数とすることで、一番始めに評価される。
	// これにより、後述の処理でのpanicは全て RecoveryMiddleware で処理される。
	mux.Handle("/todos",
		middleware.With(
			handler.NewTODOHandler(service.NewTODOService(todoDB)),
			&middleware.UserAgentRecordMiddleware{},
			&middleware.RecoveryMiddleware{},
		),
	)
	mux.Handle("/do-panic",
		middleware.With(
			handler.NewPanicHandler(),
			&middleware.UserAgentRecordMiddleware{},
			&middleware.RecoveryMiddleware{},
		),
	)

	return mux
}
