package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"mailinglist-backend-go/rest"
	"mailinglist-backend-go/rest/mailing"
	"net/http"
	"os"
)

type config struct {
	http struct {
		addr string
	}
	lg *slog.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.http.addr, "http.addr", "localhost:8080", "http listen address")
	flag.Parse()

	cfg.lg = slog.New(slog.NewJSONHandler(os.Stderr, nil)).With("app", "mailinglist-backend-go")
	cfg.lg.Info("starting", "addr", cfg.http.addr)

	if err := run(context.Background(), cfg); err != nil {
		cfg.lg.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

func run(_ context.Context, cfg config) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", rest.Health)
	mux.Handle("GET /lists", mailing.Lists(cfg.lg))
	mux.Handle("GET /add", mailing.Subscribe(cfg.lg))
	mux.Handle("GET /remove", mailing.Unsubscribe(cfg.lg))

	err := http.ListenAndServe(cfg.http.addr, mux)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server closed unexpectedly: %w", err)
	}
	return nil
}
