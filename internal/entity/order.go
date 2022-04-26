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
