package handler

import (
	"context"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	headerAcceptEncoding  = "Accept-Encoding"
	headerContentEncoding = "Content-Encoding"
	headerContentType     = "Content-Type"
	encoding              = "gzip"
	jsonContentType       = "application/json; charset=UTF-8"
)

type Service interface {
	RegisterUser(ctx context.Context, login string, password string) (*entity.User, error)
	LoginUser(ctx context.Context, login string, password string) (*entity.User, error)
	AddNewOrder(ctx context.Context, sessionID string, orderNumber string) error
	GetOrderList(ctx context.Context, sessionID string) ([]*entity.Order, error)
	GetUserBalance(ctx context.Context, sessionID string) (*entity.UserBalance, error)
	GetUserWithdrawals(ctx context.Context, sessionID string) ([]*entity.Withdrawal, error)
	UserOrderWithdraw(ctx context.Context, sessionID string, orderNumber string, sum float64) error
}

type Handler struct {
	*chi.Mux
	service Service
}

func NewHandler(service Service) *Handler {
	h := &Handler{
		Mux:     chi.NewMux(),
		service: service,
	}

	h.Use(CompressMiddleware)

	h.Post("/api/user/register", h.Register())
	h.Post("/api/user/login", h.Login())
	h.Post("/api/user/orders", h.Order())
	h.Get("/api/user/orders", h.GetOrdersList())
	h.Get("/api/user/balance", h.GetBalance())
	h.Post("/api/user/balance/withdraw", h.Withdraw())
	h.Get("/api/user/balance/withdrawals", h.GetWithdrawals())
	h.MethodNotAllowed(h.ExecIfNotAllowed())

	return h
}

func (h *Handler) ExecIfNotAllowed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid request 5", http.StatusBadRequest)
	}
}
