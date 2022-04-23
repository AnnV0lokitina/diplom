package handler

import (
	"context"
	"encoding/json"
	"errors"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) GetOrdersList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		sessionID, err := getSessionIDFromCookie(r)
		if err != nil {
			log.Info("order: no session in")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		orderList, err := h.service.GetOrderList(ctx, sessionID)
		if err != nil {
			processOrderListError(w, err)
			return
		}
		orderResponseList := make([]JSONOrderResponse, len(orderList))
		for _, order := range orderList {
			o := JSONOrderResponse{
				Number:     string(order.Number),
				Status:     order.Status.String(),
				Accrual:    order.Accrual.ToFloat(),
				UploadedAt: order.UploadedAt,
			}
			orderResponseList = append(orderResponseList, o)
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&orderResponseList); err != nil {
			http.Error(w, "Error while json conversion", http.StatusInternalServerError)
			return
		}
	}
}

func processOrderListError(w http.ResponseWriter, err error) {
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
	return
}
