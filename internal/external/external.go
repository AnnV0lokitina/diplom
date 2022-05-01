package external

import (
	"encoding/json"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

type AccrualSystem struct {
	address string
	client  http.Client
}

type JSONOrderStatusResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func NewAccrualSystem(address string) *AccrualSystem {
	client := http.Client{
		Timeout: time.Second * 1,
	}

	return &AccrualSystem{
		address: address,
		client:  client,
	}
}

func createOrderStatusFromExternal(externalStatus string) (entity.OrderStatus, error) {
	switch externalStatus {
	case "REGISTERED":
		// заказ зарегистрирован, но не начисление не рассчитано
		return entity.OrderStatusProcessing, nil
	case "INVALID":
		// заказ не принят к расчёту, и вознаграждение не будет начислено
		return entity.OrderStatusInvalid, nil
	case "PROCESSING":
		// расчёт начисления в процессе
		return entity.OrderStatusProcessing, nil
	case "PROCESSED":
		// расчёт начисления окончен
		return entity.OrderStatusProcessed, nil
	default:
		// неизвестный статус
		return entity.OrderStatusUndefined, labelError.NewLabelError(
			labelError.TypeInvalidExternalStatus,
			errors.New("invalid external status"),
		)
	}
}

func (ac *AccrualSystem) GetOrderInfo(number entity.OrderNumber) (*entity.OrderUpdateInfo, error) {
	url := ac.address + "/api/orders/" + string(number)
	log.WithField("url", url).Info("send request")
	response, err := ac.client.Get(url)
	if err != nil {
		log.WithError(err).Warning("error while send request to external")
		return nil, err
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithError(err).Warning("error while read info from external")
		return nil, err
	}
	log.WithField("body", string(respBody)).Info("receive body")
	var parsedResponse JSONOrderStatusResponse
	if err := json.Unmarshal(respBody, &parsedResponse); err != nil {
		log.WithError(err).Warning("invalid response from external")
		return nil, err
	}

	log.WithFields(log.Fields{
		"status":  parsedResponse.Status,
		"order":   parsedResponse.Order,
		"accrual": parsedResponse.Accrual,
	}).Info("parsed info")

	orderNumber := entity.OrderNumber(parsedResponse.Order)
	status, err := createOrderStatusFromExternal(parsedResponse.Status)
	if err != nil {
		log.WithError(err).Warning("error creating status from external")
		return nil, err
	}
	accrual := entity.NewPointValue(parsedResponse.Accrual)

	return &entity.OrderUpdateInfo{
		Number:  orderNumber,
		Status:  status,
		Accrual: accrual,
	}, nil
}
