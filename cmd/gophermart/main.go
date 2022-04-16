package main

import (
	"context"
	handlerPkg "github.com/AnnV0lokitina/diplom/internal/handler"
	"github.com/AnnV0lokitina/diplom/internal/repo"
	"github.com/AnnV0lokitina/diplom/internal/service"
	"log"
)

func main() {
	cfg := initConfig()
	initParams(cfg)
	doMigrates(cfg.DataBaseURI)

	ctx := context.Background()
	repo, err := repo.NewRepo(ctx, cfg.DataBaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close(ctx)

	service := service.NewService(repo)
	handler := handlerPkg.NewHandler(service)
	application := NewApp(handler)

	application.Run(cfg.RunAddress)
}
