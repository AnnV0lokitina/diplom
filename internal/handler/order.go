package handler

import (
	"context"
	"errors"
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
		num, err := io.ReadAll(r.Body)
		if err != nil || len(num) == 0 {
			log.Info("invalid request format")
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}
		err = h.service.AddNewOrder(ctx, sessionID, string(num))
		if err != nil {
			processOrderError(w, err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func processOrderError(w http.ResponseWriter, err error) {
	var labelErr *labelError.LabelError
	if errors.As(err, &labelErr) {
		if labelErr.Label == labelError.TypeUnauthorized {
			log.Info("user unauthorized")
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}
		if labelErr.Label == labelError.TypeInvalidData {
			log.Info("invalid order number")
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
			return
		}
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
	log.WithError(err).Info("server error")
	http.Error(w, "Server error", http.StatusInternalServerError)
}
