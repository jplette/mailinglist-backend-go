package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"mailinglist-backend-go/controller/health"
	"mailinglist-backend-go/controller/mailing"
	"mailinglist-backend-go/services/configReader"
	"mailinglist-backend-go/services/requestValidator"
	"net/http"
	"os"
	"strings"
)

type config struct {
	http struct {
		addr string
	}
	lg *slog.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.http.addr, "http.addr", ":8080", "http listen address")
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
	// Unprotected health endpoint
	mux.HandleFunc("/health", health.Ping)
	// Protected endpoints wrapped by authMiddleware
	mux.Handle("GET /lists", authMiddleware(mailing.Lists(cfg.lg)))
	mux.Handle("POST /subscribe", authMiddleware(mailing.Subscribe(cfg.lg)))
	mux.Handle("POST /unsubscribe", authMiddleware(mailing.Unsubscribe(cfg.lg)))

	// Setup CORS middleware with allowed origins from environment
	allowed := configReader.Values("CORS_ALLOWED_ORIGINS")
	handler := corsMiddleware(allowed)(mux)

	err := http.ListenAndServe(cfg.http.addr, handler)
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server closed unexpectedly: %w", err)
	}
	return nil
}

// authMiddleware validates JWT from Authorization header and stores claims in context.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := requestValidator.ValidateRequest(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := requestValidator.WithClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// corsMiddleware returns a middleware that sets CORS headers based on allowed origins.
// allowedOrigins is a list of origins (scheme://host[:port]) or "*" to allow any.
func corsMiddleware(allowedOrigins []string) func(http.Handler) http.Handler {
	// Normalize and trim values
	var normalized []string
	for _, o := range allowedOrigins {
		ot := strings.TrimSpace(o)
		if ot != "" {
			normalized = append(normalized, ot)
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			allowed := isOriginAllowed(origin, normalized)

			if origin != "" && allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
				// If you need credentials in future, enable and ensure specific origins only
				// w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			// Always advertise what methods/headers are accepted for preflight
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

			if r.Method == http.MethodOptions {
				if origin == "" || !allowed {
					// Not a valid CORS preflight; respond with 204 to be safe
					w.WriteHeader(http.StatusNoContent)
					return
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func isOriginAllowed(origin string, allowed []string) bool {
	if origin == "" {
		return false
	}
	for _, a := range allowed {
		if a == "*" || origin == a {
			return true
		}
	}
	return false
}
