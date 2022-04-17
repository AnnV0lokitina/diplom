package entity

import "time"

type Order struct {
	Number     OrderNumber
	Login      string
	Status     OrderStatus
	UploadedAt time.Time
}
