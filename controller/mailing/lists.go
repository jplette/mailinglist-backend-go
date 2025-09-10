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

// Lists godoc
// @Summary      Get mailing lists
// @Description  Returns all mailing lists available to the system.
// @Tags         mailing
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   mailgun.APIMailingList
// @Failure      401  {string}  string  "Unauthorized"
// @Failure      500  {string}  string  "Internal Server Error"
// @Router       /lists [get]
// Lists returns an [http.Handler] that returns a list of mailing lists.
func Lists(lg *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Authorization is handled by middleware

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

// Subscribe godoc
// @Summary      Subscribe a member to a list
// @Description  Subscribes the specified member email to the given list address.
// @Tags         mailing
// @Accept       application/x-www-form-urlencoded
// @Produce      json
// @Security     BearerAuth
// @Param        list    formData  string  true  "List address"
// @Param        member  formData  string  true  "Member email"
// @Success      200     {string}  string  "OK"
// @Failure      400     {string}  string  "Bad Request"
// @Failure      401     {string}  string  "Unauthorized"
// @Router       /subscribe [post]
func Subscribe(lg *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Auth handled by middleware; fetch claims from context
		claims, err := requestValidator.ClaimsFromRequest(r)
		if err != nil {
			httpErrorUnauthorized(w, r, lg, err)
			return
		}

		listAddress := r.PostFormValue("list")
		memberAddress := r.PostFormValue("member")

		if listAddress == "" || memberAddress == "" {
			httpErrorBadRequest(w, r, lg, fmt.Errorf("failed to parse form: %w", err))
			return
		}

		user := requestValidator.CurrentUser(claims)
		// If not admin, you can only subscribe yourself
		if (user.Admin == false) && (memberAddress != user.Email) {
			httpErrorBadRequest(w, r, lg, fmt.Errorf("only admins can (un)subscribe other users"))
		}

		// TODO: Check if the list is blocked

		err = mailgun.Subscribe(listAddress, memberAddress)
		if err != nil {
			httpErrorBadRequest(w, r, lg, fmt.Errorf("failed to subsribe: %w", err))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}

// Unsubscribe godoc
// @Summary      Unsubscribe a member from a list
// @Description  Unsubscribes the specified member email from the given list address.
// @Tags         mailing
// @Accept       application/x-www-form-urlencoded
// @Produce      json
// @Security     BearerAuth
// @Param        list    formData  string  true  "List address"
// @Param        member  formData  string  true  "Member email"
// @Success      200     {string}  string  "OK"
// @Failure      400     {string}  string  "Bad Request"
// @Failure      401     {string}  string  "Unauthorized"
// @Router       /unsubscribe [post]
func Unsubscribe(lg *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Auth handled by middleware; fetch claims from context
		claims, err := requestValidator.ClaimsFromRequest(r)
		if err != nil {
			httpErrorUnauthorized(w, r, lg, err)
			return
		}

		listAddress := r.PostFormValue("list")
		memberAddress := r.PostFormValue("member")

		if listAddress == "" || memberAddress == "" {
			httpErrorBadRequest(w, r, lg, fmt.Errorf("failed to parse form: %w", err))
			return
		}

		user := requestValidator.CurrentUser(claims)
		// If not admin, you can only subscribe yourself
		if (user.Admin == false) && (memberAddress != user.Email) {
			httpErrorBadRequest(w, r, lg, fmt.Errorf("only admins can (un)subscribe other users"))
		}

		// TODO: Check if the list is blocked

		err = mailgun.Unsubscribe(listAddress, memberAddress)
		if err != nil {
			httpErrorBadRequest(w, r, lg, fmt.Errorf("failed to unsubsribe: %w", err))
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
	http.Error(w, err.Error(), http.StatusBadRequest)
}
