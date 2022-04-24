package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net/http"
	"time"
)

func (s *Service) CreateGetOrderInfoProcess(ctx context.Context, accrualSystemAddress string, nOfWorkers int) {
	s.createCheckOrderWorkerPull(ctx, accrualSystemAddress, nOfWorkers)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			orderNumberList, err := s.repo.GetOrdersListForRequest(ctx)
			if err != nil {
				var labelErr *labelError.LabelError
				if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeNotFound {
					time.Sleep(200 * time.Millisecond)
					continue
				}
				log.WithError(err).Warning("error get orders")
				return
			}
			for _, orderNumber := range orderNumberList {
				job := entity.JobCheckOrder{
					Number: orderNumber,
				}
				if s.jobCheckOrder != nil {
					s.jobCheckOrder <- &job
				}
			}
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (s *Service) updateStatus(ctx context.Context, accrualSystemAddress string, job *entity.JobCheckOrder) error {
	client := http.Client{}
	client.Timeout = time.Second * 1
	response, err := client.Get(accrualSystemAddress + "/api/orders/" + string(job.Number))
	if err != nil {
		log.WithError(err).Warning("error while send request to external")
		return err
	}
	defer response.Body.Close()
	reader := response.Body
	respBody, err := ioutil.ReadAll(reader)
	if err != nil {
		log.WithError(err).Warning("error while read info from external")
		return err
	}
	var parsedResponse entity.JSONOrderStatusResponse
	if err := json.Unmarshal(respBody, &parsedResponse); err != nil {
		log.WithError(err).Warning("invalid response from external")
		return err
	}
	orderNumber := entity.OrderNumber(parsedResponse.Order)
	status, err := entity.NewOrderStatusFromExternal(parsedResponse.Status)
	if err != nil {
		log.WithError(err).Warning("error creating status from external")
		return err
	}
	accrual := entity.NewPointValue(parsedResponse.Accrual)
	err = s.repo.AddOrderInfo(ctx, orderNumber, status, accrual)
	if err != nil {
		log.WithError(err).Warning("error while updating order info")
		return err
	}
	return nil
}

func (s *Service) createCheckOrderWorkerPull(ctx context.Context, accrualSystemAddress string, nOfWorkers int) {
	s.jobCheckOrder = make(chan *entity.JobCheckOrder)
	g, _ := errgroup.WithContext(ctx)

	for i := 1; i <= nOfWorkers; i++ {
		g.Go(func() error {
			for job := range s.jobCheckOrder {
				err := s.updateStatus(ctx, accrualSystemAddress, job)
				if err != nil {
					continue
				}
			}
			return nil
		})
	}

	go func() {
		<-ctx.Done()
		close(s.jobCheckOrder)
	}()

	if err := g.Wait(); err != nil {
		log.WithError(err).Warning("error in worker pool")
	}
}