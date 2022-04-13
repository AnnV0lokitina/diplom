package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	"io"
	"net/http"
)

func (h *Handler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Register")
		ctx := context.Background()

		request, err := io.ReadAll(r.Body)
		if err != nil || len(request) == 0 {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		var parsedRequest JSONUserRequest
		if err := json.Unmarshal(request, &parsedRequest); err != nil {
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		user, err := h.service.RegisterUser(ctx, parsedRequest.Login, parsedRequest.Password)
		if err == nil {
			addSessionIDToCookie(w, user.ActiveSessionID)
			w.WriteHeader(http.StatusOK)
		}
		var labelErr *labelError.LabelError
		if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeConflict {
			http.Error(w, "Login existed", http.StatusConflict)
			return
		}
		fmt.Println(err)
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}
