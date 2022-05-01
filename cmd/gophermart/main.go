package main

import (
	"context"
	"github.com/AnnV0lokitina/diplom/internal/external"
	handlerPkg "github.com/AnnV0lokitina/diplom/internal/handler"
	repoPkg "github.com/AnnV0lokitina/diplom/internal/repo"
	servicePkg "github.com/AnnV0lokitina/diplom/internal/service"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := initConfig()
	initParams(cfg)
	err := doMigrates(cfg.DataBaseURI)
	if err != nil {
		log.WithError(err).Fatal("migrations error")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()
	repo, err := repoPkg.NewRepo(ctx, cfg.DataBaseURI)
	if err != nil {
		log.WithError(err).Fatal("error repo init")
	}
	defer repo.Close(ctx)

	accrualSystem := external.NewAccrualSystem(cfg.AccrualSystemAddress)

	service := servicePkg.NewService(repo, accrualSystem)
	handler := handlerPkg.NewHandler(service)
	application := NewApp(handler)

	go func() {
		service.CreateGetOrderInfoProcess(ctx, cfg.NumOfWorkers)
	}()

	err = application.Run(ctx, cfg.RunAddress)
	if err != nil {
		log.Fatal(err)
	}
}
