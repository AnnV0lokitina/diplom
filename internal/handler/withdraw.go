package handler

import (
	"context"
	"encoding/json"
	"errors"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (h *Handler) Withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			log.Info("order: no session in")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		request, err := io.ReadAll(r.Body)
		if err != nil || len(request) == 0 {
			log.Info("invalid request format")
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		var parsedRequest JSONWithdrawRequest
		if err := json.Unmarshal(request, &parsedRequest); err != nil {
			http.Error(w, "Invalid request 7", http.StatusBadRequest)
			return
		}
		err = h.service.UserOrderWithdraw(ctx, sessionID, parsedRequest.Order, parsedRequest.Sum)
		if err != nil {
			processWithdrawError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func processWithdrawError(w http.ResponseWriter, err error) {
	var labelErr *labelError.LabelError
	if errors.As(err, &labelErr) {
		if labelErr.Label == labelError.TypeUnauthorized {
			log.Info("user unauthorized")
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}
		if labelErr.Label == labelError.TypeNotEnoughPoints {
			log.Info("not enough points")
			http.Error(w, "Not enough points", http.StatusPaymentRequired)
			return
		}
		if labelErr.Label == labelError.TypeNotFound {
			log.Info("order not found")
			http.Error(w, "Order not found", http.StatusUnprocessableEntity)
			return
		}
	}
	log.Info("server error")
	http.Error(w, "Server error", http.StatusInternalServerError)
}
