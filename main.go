package main

import (
	"avocato-db/src/api/handlers"
	"avocato-db/src/config"
	"avocato-db/src/core"
	"avocato-db/src/core/ledger"
	"avocato-db/src/storage/postgres"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.LoadConfig()
	core.LogInfo("Starting avocato-db on port %s...", cfg.APIPort)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. Connect to Postgres
	db, err := postgres.Connect(ctx, cfg.DBURL)
	if err != nil {
		core.FatalError("Failed to connect to database: %v", err)
	}
	defer db.Close(context.Background())

	// 2. Initialize Ledger Service
	ledgerService, err := ledger.NewService(cfg.WALPath, db)
	if err != nil {
		core.FatalError("Failed to initialize ledger service: %v", err)
	}

	// 3. Setup Routes
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/append", handlers.NewAppendHandler(ledgerService))

	server := &http.Server{
		Addr:    ":" + cfg.APIPort,
		Handler: mux,
	}

	// 4. Start Server
	go func() {
		core.LogInfo("API ready on :%s", cfg.APIPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			core.LogError("Server failed: %v", err)
		}
	}()

	<-ctx.Done()
	core.LogInfo("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		core.LogError("Server shutdown failed: %v", err)
	}
	
	core.LogInfo("Exiting.")
}
