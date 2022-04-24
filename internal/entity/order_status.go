package entity

import (
	"errors"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
)

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

func NewOrderStatusFromExternal(externalStatus string) (OrderStatus, error) {
	switch externalStatus {
	case "REGISTERED":
		// заказ зарегистрирован, но не начисление не рассчитано
		return OrderStatusNew, nil
	case "INVALID":
		// заказ не принят к расчёту, и вознаграждение не будет начислено
		return OrderStatusInvalid, nil
	case "PROCESSING":
		// расчёт начисления в процессе
		return OrderStatusProcessed, nil
	case "PROCESSED":
		// расчёт начисления окончен
		return OrderStatusProcessed, nil
	default:
		// неизвестный статус
		return OrderStatusUndefined, labelError.NewLabelError(
			labelError.TypeInvalidExternalStatus,
			errors.New("invalid external status"),
		)
	}
}
