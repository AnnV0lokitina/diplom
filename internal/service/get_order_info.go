package service

import (
	"context"
	"errors"
	labelError "github.com/AnnV0lokitina/diplom/pkg/error"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"time"
)

func (s *Service) CreateGetOrderInfoProcess(ctx context.Context, nOfWorkers int) {
	s.jobCheckOrder = make(chan *JobCheckOrder)
	s.sendOrdersToCheck(ctx)
	s.createCheckOrderWorkerPull(ctx, nOfWorkers)
}

func (s *Service) sendOrdersToCheck(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				orderNumberList, err := s.repo.GetOrdersListForRequest(ctx)
				if err != nil {
					var labelErr *labelError.LabelError
					if errors.As(err, &labelErr) && labelErr.Label == labelError.TypeNotFound {
						continue
					}
					log.WithError(err).Warning("error get orders")
					return
				}
				for _, orderNumber := range orderNumberList {
					job := JobCheckOrder{
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

func (s *Service) updateStatus(ctx context.Context, job *JobCheckOrder) error {
	orderInfo, err := s.accrualSystem.GetOrderInfo(job.Number)
	if err != nil {
		return err
	}
	err = s.repo.AddOrderInfo(ctx, orderInfo)
	if err != nil {
		log.WithError(err).Warning("error while updating order info")
		return err
	}
	log.Info("info added")
	return nil
}

func (s *Service) createCheckOrderWorkerPull(ctx context.Context, nOfWorkers int) {
	g, _ := errgroup.WithContext(ctx)

	for i := 1; i <= nOfWorkers; i++ {
		g.Go(func() error {
			for job := range s.jobCheckOrder {
				err := s.updateStatus(ctx, job)
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
