package main

import (
	"database/sql"
	"embed"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

const dbType = "postgres"

func doMigrates(dsn string) error {
	db, err := sql.Open(dbType, dsn)
	if err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dbType); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	return nil
}
