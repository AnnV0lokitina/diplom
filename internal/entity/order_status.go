package entity

type OrderStatus int

const (
	OrderStatusNew OrderStatus = iota
	OrderStatusProcessing
	OrderStatusInvalid
	OrderStatusProcessed
)

func (os OrderStatus) String() string {
	switch os {
	case OrderStatusNew:
		return "NEW"
	case OrderStatusProcessing:
		return "PROCESSING"
	case OrderStatusInvalid:
		return "INVALID"
	case OrderStatusProcessed:
		return "PROCESSED"
	default:
		return "UNKNOWN"
	}
}
