package handler

import (
	"context"
	"encoding/json"
	"errors"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) GetWithdrawals() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			log.Info("order: no session in")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		withdrawList, err := h.service.GetUserWithdrawals(ctx, sessionID)
		if err != nil {
			processWithdrawListError(w, err)
			return
		}
		withdrawResponseList := make([]JSONWithdrawResponse, 0, len(withdrawList))
		for _, withdraw := range withdrawList {
			w := JSONWithdrawResponse{
				Order:       withdraw.OrderNumber,
				Sum:         withdraw.Sum.ToFloat(),
				ProcessedAt: withdraw.ProcessedAt,
			}
			withdrawResponseList = append(withdrawResponseList, w)
		}
		w.Header().Set(headerContentType, jsonContentType)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&withdrawResponseList); err != nil {
			http.Error(w, "Error while json conversion", http.StatusInternalServerError)
			return
		}
	}
}

func processWithdrawListError(w http.ResponseWriter, err error) {
	var labelErr *labelError.LabelError
	if errors.As(err, &labelErr) {
		if labelErr.Label == labelError.TypeUnauthorized {
			log.Info("user unauthorized")
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}
		if labelErr.Label == labelError.TypeNotFound {
			log.Info("no data in response")
			http.Error(w, "No data", http.StatusNoContent)
			return
		}
	}
	log.Info("server error")
	http.Error(w, "Server error", http.StatusInternalServerError)
}
