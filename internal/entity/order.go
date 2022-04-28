package entity

import "time"

type Order struct {
	ID         int
	Number     OrderNumber
	Login      string
	Accrual    PointValue
	Status     OrderStatus
	UploadedAt time.Time
}

type OrderUpdateInfo struct {
	Number  OrderNumber
	Status  OrderStatus
	Accrual PointValue
}
