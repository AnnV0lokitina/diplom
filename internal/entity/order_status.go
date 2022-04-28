package entity

type OrderStatus int

const (
	// OrderStatusNew заказ загружен в систему, но не попал в обработку
	OrderStatusNew OrderStatus = iota
	// OrderStatusProcessing вознаграждение за заказ рассчитывается
	OrderStatusProcessing
	// OrderStatusInvalid система расчёта вознаграждений отказала в расчёте
	OrderStatusInvalid
	// OrderStatusProcessed данные по заказу проверены и информация о расчёте успешно получена
	OrderStatusProcessed
	// OrderStatusUndefined неверный статус
	OrderStatusUndefined = 1000
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
