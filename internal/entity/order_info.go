package entity

type JobCheckOrder struct {
	Number OrderNumber
}

type JSONOrderStatusResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}
