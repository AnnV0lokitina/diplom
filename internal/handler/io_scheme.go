package handler

import "time"

type JSONUserRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type JSONWithdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type JSONOrderResponse struct {
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    float64   `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type JSONBalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type JSONWithdrawResponse struct {
	Order       string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

type JSONOrderStatusResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
