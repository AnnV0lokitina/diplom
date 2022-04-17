package handler

import (
	"context"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func (h *Handler) Order() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			log.Info("order: no session in")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		data, err := io.ReadAll(r.Body)
		if err != nil || len(data) == 0 {
			log.Info("invalid request format")
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		orderNumber := entity.OrderNumber(data)
		if !orderNumber.Valid() {
			log.Info("invalid order number")
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
			return
		}
		err = h.service.AddNewOrder(ctx, sessionID, orderNumber)
		if err == nil {
			w.WriteHeader(http.StatusAccepted)
			return
		}
		var labelErr *labelError.LabelError
		if errors.As(err, &labelErr) {
			if labelErr.Label == labelError.TypeCreated {
				log.Info("order existed")
				w.WriteHeader(http.StatusOK)
				return
			}
			if labelErr.Label == labelError.TypeConflict {
				log.Info("order created with another user")
				w.WriteHeader(http.StatusConflict)
				return
			}
		}

		log.Info("server error")
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func (h *Handler) GetOrdersList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			log.WithFields(log.Fields{
				"session ID": sessionID,
			}).Info("authorization failed")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		log.WithFields(log.Fields{
			"session ID": sessionID,
		}).Info("authorization success")
		w.WriteHeader(http.StatusOK)
	}
}
