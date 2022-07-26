package handler

import (
	"context"
	"encoding/json"
	"errors"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			log.Info("order: no session in")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		balance, err := h.service.GetUserBalance(ctx, sessionID)
		if err != nil {
			var labelErr *labelError.LabelError
			if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeUnauthorized {
				log.Info("illegal login or password")
				http.Error(w, "Illegal login or password", http.StatusUnauthorized)
				return
			}
			log.WithError(err).Info("error when register balance handler")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		balanceResponse := JSONBalanceResponse{
			Current:   balance.Current.ToFloat(),
			Withdrawn: balance.Withdrawn.ToFloat(),
		}
		w.Header().Set(headerContentType, jsonContentType)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&balanceResponse); err != nil {
			http.Error(w, "Error while json conversion", http.StatusInternalServerError)
			return
		}
	}
}
