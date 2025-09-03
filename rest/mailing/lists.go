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
		lists, err := mailinglist.Lists(false)

		if err != nil {
			httpError(w, r, lg, fmt.Errorf("failed to get lists: %w", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Need to do this because the json encoder will not encode an empty array. It is nil instead
		// Could change with "encoding/json/v2"
		result := []mailinglist.MGMailingList{}
		if len(lists) > 0 {
			result = lists
		}
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			httpError(w, r, lg, fmt.Errorf("failed to get lists: %w", err))
			return
		}
	})
}

func Subscribe(lg *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := mailinglist.Subscribe("justatest@mailgun.wohnsinn-bessungen.de", "test@jakumba.com")
		if err != nil {
			httpError(w, r, lg, fmt.Errorf("failed to subscribe to lists: %w", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

func Unsubscribe(lg *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := mailinglist.Unsubscribe("justatest@mailgun.wohnsinn-bessungen.de", "test@jakumba.com")
		if err != nil {
			httpError(w, r, lg, fmt.Errorf("failed to unsubscribe from lists: %w", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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
	case errors.Is(err, mailinglist.ErrForbidden):
		code = http.StatusForbidden
	}
	if code == http.StatusInternalServerError {
		lg.ErrorContext(r.Context(), "internal", "error", err)
		err = mailinglist.ErrInternal
	}
	http.Error(w, err.Error(), code)
}
