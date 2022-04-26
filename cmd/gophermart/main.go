package main

import (
	"context"
	handlerPkg "github.com/AnnV0lokitina/diplom/internal/handler"
	"github.com/AnnV0lokitina/diplom/internal/repo"
	"github.com/AnnV0lokitina/diplom/internal/service"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

const nOfWorkers = 3

func main() {
	cfg := initConfig()
	initParams(cfg)
	doMigrates(cfg.DataBaseURI)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()
	repo, err := repo.NewRepo(ctx, cfg.DataBaseURI)
	if err != nil {
		log.WithError(err).Fatal("error repo init")
	}
	defer repo.Close(ctx)

	service := service.NewService(repo)
	handler := handlerPkg.NewHandler(service)
	application := NewApp(handler)

	go func() {
		service.CreateGetOrderInfoProcess(ctx, cfg.AccrualSystemAddress, nOfWorkers)
	}()

	err = application.Run(ctx, cfg.RunAddress)
	if err != nil {
		log.Fatal(err)
	}
}
