package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"stockchallenge/backend/internal/config"
	"stockchallenge/backend/internal/db"
	httpapi "stockchallenge/backend/internal/http"
	"stockchallenge/backend/internal/stocks"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	stocksService := stocks.NewService(
		pool,
		cfg.StocksAPIURL,
		cfg.StocksAPIToken,
		cfg.SyncTimeout,
		cfg.SyncMaxPages,
	)
	if err := db.Migrate(ctx, pool); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	if cfg.AutoSyncOnStartup {
		result, err := stocksService.Sync(ctx, cfg.SyncMaxPages)
		if err != nil {
			log.Printf("startup sync failed: %v", err)
		} else {
			log.Printf("startup sync completed: pages=%d stocks=%d", result.PagesProcessed, result.StocksSaved)
		}
	}

	handler := httpapi.NewRouter(pool, stocksService, cfg.CORSAllowedOrigins)
	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("server shutdown error: %v", err)
		}
	}()

	log.Printf("backend listening on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
