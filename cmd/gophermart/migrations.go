package main

import (
	"database/sql"
	"embed"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const dbType = "postgres"

func doMigrates(dsn string) {
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		log.WithError(err).Fatal("migrations error: db open")
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dbType); err != nil {
		log.WithError(err).Fatal("migrations error: set dialect")
	}

	if err := goose.Up(db, "migrations"); err != nil {
		log.WithError(err).Fatal("migrations error: migrations up")
	}
}
