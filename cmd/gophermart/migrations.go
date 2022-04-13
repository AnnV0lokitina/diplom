package main

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const pathToMigrations = "file://../../migrations"

func doMigrates(dsn string) {
	m, err := migrate.New(pathToMigrations, dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		fmt.Println(err)
	}
}
