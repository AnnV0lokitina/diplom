package main

import (
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

type config struct {
	RunAddress           string `env:"RUN_ADDRESS"  envDefault:"localhost:8080"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:""`
	DataBaseURI          string `env:"DATABASE_URI" envDefault:""`
}

func initConfig() *config {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		log.WithError(err).Fatal("init environment error")
	}
	return &cfg
}
