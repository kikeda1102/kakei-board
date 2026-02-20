package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kikeda1102/kakei-board/backend/internal/database"
	"github.com/kikeda1102/kakei-board/backend/internal/eventstore"
	"github.com/kikeda1102/kakei-board/backend/internal/expense"
	"github.com/kikeda1102/kakei-board/backend/internal/middleware"
	"github.com/kikeda1102/kakei-board/backend/migrations"
)

const shutdownTimeout = 10 * time.Second

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	cfg, err := database.ConfigFromEnv()
	if err != nil {
		log.Fatalf("database config: %v", err)
	}

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("database open: %v", err)
	}
	defer db.Close()

	if err := migrations.Run(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}

	handler := buildHandler(db)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-done
	log.Println("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
	log.Println("server stopped")
}

func buildHandler(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := db.PingContext(ctx); err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			log.Printf("health check failed: %v", err)
			if _, wErr := w.Write([]byte(`{"status":"unhealthy"}`)); wErr != nil {
				log.Printf("failed to write response: %v", wErr)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	})

	store := eventstore.NewMySQLStore(db)
	projector := expense.NewProjector(db)
	repo := expense.NewRepository(db)
	expenseHandler := expense.NewHandler(store, projector, repo)
	expenseHandler.Register(mux)

	return middleware.CORS(mux)
}
