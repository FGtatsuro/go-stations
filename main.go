package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/TechBowl-japan/go-stations/db"
	"github.com/TechBowl-japan/go-stations/handler/router"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatalln("main: failed to exit successfully, err =", err)
	}
}

func realMain() error {
	// config values
	const (
		defaultPort   = ":8080"
		defaultDBPath = ".sqlite3/todo.db"
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = defaultDBPath
	}

	// set time zone
	var err error
	time.Local, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return err
	}

	// set up sqlite3
	todoDB, err := db.NewDB(dbPath)
	if err != nil {
		return err
	}
	defer todoDB.Close()

	mux, err := router.NewHandlerWithBasicAuth(
		todoDB,
		os.Getenv("BASIC_AUTH_USER_ID"),
		os.Getenv("BASIC_AUTH_PASSWORD"),
	)
	if err != nil {
		return err
	}
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()
	var wg sync.WaitGroup

	// NOTE: serverの数だけAddする
	wg.Add(1)
	go run(ctx, &wg, server)
	wg.Wait()

	return nil
}

// run はHTTPサーバに対するGraceful shutdownを提供する。
//
// [context.Context] 及び [sync.WaitGroup]を共有する事で複数サーバのGraceful shutdownを同時に制御できる。
func run(ctx context.Context, wg *sync.WaitGroup, srv *http.Server) {
	go func() {
		defer wg.Done()

		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("main: could not gracefully shutdown the server, err =%v\n", err)
		} else {
			log.Printf("main: server is completely shutdown\n")
		}
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Printf("main: could not listen on %v, err =%v\n", srv.Addr, err)
	} else {
		log.Printf("main: listen port %v is closed\n", srv.Addr)
	}
}
