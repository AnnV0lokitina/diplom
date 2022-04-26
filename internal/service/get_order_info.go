package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AnnV0lokitina/diplom/internal/entity"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net/http"
	"time"
)

func (s *Service) CreateGetOrderInfoProcess(ctx context.Context, accrualSystemAddress string, nOfWorkers int) {
	s.jobCheckOrder = make(chan *entity.JobCheckOrder)
	s.client = http.Client{}
	s.sendOrdersToCheck(ctx)
	fmt.Println("CreateGetOrderInfoProcess")
	s.createCheckOrderWorkerPull(ctx, accrualSystemAddress, nOfWorkers)
}

func (s *Service) sendOrdersToCheck(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				fmt.Println("start")
				orderNumberList, err := s.repo.GetOrdersListForRequest(ctx)
				fmt.Println(len(orderNumberList))
				if err != nil {
					var labelErr *labelError.LabelError
					if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeNotFound {
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
			}
		}
	}()
}

func (s *Service) updateStatus(ctx context.Context, accrualSystemAddress string, job *entity.JobCheckOrder) error {
	s.client.Timeout = time.Second * 1
	url := accrualSystemAddress + "/api/orders/" + string(job.Number)
	log.WithField("url", url).Info("send request")
	response, err := s.client.Get(url)
	if err != nil {
		log.WithError(err).Warning("error while send request to external")
		return err
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.WithError(err).Warning("error while read info from external")
		return err
	}
	log.WithField("body", string(respBody)).Info("receive body")
	var parsedResponse entity.JSONOrderStatusResponse
	if err := json.Unmarshal(respBody, &parsedResponse); err != nil {
		log.WithError(err).Warning("invalid response from external")
		return err
	}

	log.WithFields(log.Fields{
		"status":  parsedResponse.Status,
		"order":   parsedResponse.Order,
		"accrual": parsedResponse.Accrual,
	}).Info("parsed info")

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
	log.Info("info added")
	return nil
}

func (s *Service) createCheckOrderWorkerPull(ctx context.Context, accrualSystemAddress string, nOfWorkers int) {
	fmt.Println("createCheckOrderWorkerPull")
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
