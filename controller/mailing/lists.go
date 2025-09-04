package mailing

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mailinglist-backend-go/services/common"
	"mailinglist-backend-go/services/mailgun"
	"mailinglist-backend-go/services/requestValidator"
	"net/http"
)

// Lists returns an [http.Handler] that returns a list of mailing lists.
func Lists(lg *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Check if the user is authorized
		_, err := requestValidator.ValidateRequest(r)
		if err != nil {
			httpErrorUnauthorized(w, r, lg, err)
			return
		}

		// Get the list of mailing lists
		lists, err := mailgun.Lists(false)

		if err != nil {
			httpError(w, r, lg, fmt.Errorf("failed to get lists: %w", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// Need to do this because the json encoder will not encode an empty array. It is nil instead
		// Could change with "encoding/json/v2"
		result := []mailgun.MGMailingList{}
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

		// Check if the user is authorized
		claims, err := requestValidator.ValidateRequest(r)
		if err != nil {
			httpErrorUnauthorized(w, r, lg, err)
			return
		}

		err = r.ParseForm()
		if err != nil {
			httpErrorBadRequest(w, r, lg, fmt.Errorf("failed to parse form: %w", err))
			return
		}

		listAddress := r.Form.Get("list")
		memberAddress := r.Form.Get("member")

		user := requestValidator.CurrentUser(claims)
		// If not admin, you can only subscribe yourself
		if (user.Admin == false) && (memberAddress != user.Email) {
			httpErrorBadRequest(w, r, lg, fmt.Errorf("only admins can (un)subscribe other users"))
		}

		// TODO: Check if the list is blocked

		err = mailgun.Subscribe(listAddress, memberAddress)
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

		// Check if the user is authorized
		claims, err := requestValidator.ValidateRequest(r)
		if err != nil {
			httpErrorUnauthorized(w, r, lg, err)
			return
		}

		err = r.ParseForm()
		if err != nil {
			httpErrorBadRequest(w, r, lg, fmt.Errorf("failed to parse form: %w", err))
			return
		}

		listAddress := r.Form.Get("list")
		memberAddress := r.Form.Get("member")

		user := requestValidator.CurrentUser(claims)
		// If not admin, you can only subscribe yourself
		if (user.Admin == false) && (memberAddress != user.Email) {
			httpErrorBadRequest(w, r, lg, fmt.Errorf("only admins can (un)subscribe other users"))
		}

		// TODO: Check if the list is blocked

		err = mailgun.Unsubscribe(listAddress, memberAddress)
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
	case errors.Is(err, common.ErrBadRequest):
		code = http.StatusBadRequest
	case errors.Is(err, common.ErrConflict):
		code = http.StatusConflict
	case errors.Is(err, common.ErrNotFound):
		code = http.StatusNotFound
	case errors.Is(err, common.ErrForbidden):
		code = http.StatusForbidden
	}
	if code == http.StatusInternalServerError {
		lg.ErrorContext(r.Context(), "internal", "error", err)
		err = common.ErrInternal
	}
	http.Error(w, err.Error(), code)
}

func httpErrorUnauthorized(w http.ResponseWriter, r *http.Request, lg *slog.Logger, err error) {
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	lg.ErrorContext(r.Context(), "authorization failed", "error", err.Error())
}

func httpErrorBadRequest(w http.ResponseWriter, r *http.Request, lg *slog.Logger, err error) {
	http.Error(w, "Bad request", http.StatusBadRequest)
}
