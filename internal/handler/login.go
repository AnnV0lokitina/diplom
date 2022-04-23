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

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		request, err := io.ReadAll(r.Body)
		if err != nil || len(request) == 0 {
			log.Info("invalid request format")
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		var parsedRequest JSONUserRequest
		if err := json.Unmarshal(request, &parsedRequest); err != nil {
			log.WithError(err).Info("invalid request format")
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		user, err := h.service.LoginUser(ctx, parsedRequest.Login, parsedRequest.Password)
		if err != nil {
			var labelErr *labelError.LabelError
			if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeUnauthorized {
				log.Info("illegal login or password")
				http.Error(w, "Illegal login or password", http.StatusUnauthorized)
				return
			}
			log.WithError(err).Info("error when register")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		addSessionIDToCookie(w, user.ActiveSessionID)
		w.WriteHeader(http.StatusOK)
	}
}
