package mailing

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mailinglist-backend-go/mailinglist"
	"net/http"
)

// Lists returns an [http.Handler] that returns a list of mailing lists.
func Lists(lg *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lists, err := mailinglist.Lists()

		if err != nil {
			httpError(w, r, lg, fmt.Errorf("failed to get lists: %w", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		err = json.NewEncoder(w).Encode(lists)
		if err != nil {
			httpError(w, r, lg, fmt.Errorf("failed to get lists: %w", err))
			return
		}
	})
}

func httpError(w http.ResponseWriter, r *http.Request, lg *slog.Logger, err error) {
	code := http.StatusInternalServerError
	switch {
	case errors.Is(err, mailinglist.ErrBadRequest):
		code = http.StatusBadRequest
	case errors.Is(err, mailinglist.ErrConflict):
		code = http.StatusConflict
	case errors.Is(err, mailinglist.ErrNotFound):
		code = http.StatusNotFound
	}
	if code == http.StatusInternalServerError {
		lg.ErrorContext(r.Context(), "internal", "error", err)
		err = mailinglist.ErrInternal
	}
	http.Error(w, err.Error(), code)
}
