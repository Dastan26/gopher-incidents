// Сервис справочника брендов. Запуск: make run, затем
// curl -s localhost:8080/brands
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopher-incidents/internal/brand"
	"gopher-incidents/internal/httpapi"
	"gopher-incidents/internal/store"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	repo := store.NewMemory()
	seed(log, repo)

	srv := &http.Server{
		Addr:              ":" + env("PORT", "8080"),
		Handler:           httpapi.New(repo, log, env("ALLOWED_ORIGIN", "http://localhost:5173")),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server stopped", "err", err)
			stop()
		}
	}()

	<-ctx.Done()
	log.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown", "err", err)
	}
}

func seed(log *slog.Logger, repo *store.Memory) {
	for _, name := range []string{"Рахат", "Кока-Кола", "Простоквашино"} {
		if _, err := repo.Create(context.Background(), brand.Brand{NameKz: name, NameRu: name}); err != nil {
			log.Error("seed", "name", name, "err", err)
			os.Exit(1)
		}
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
