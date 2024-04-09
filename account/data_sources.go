package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type dataSources struct {
	DB *sqlx.DB
}

// INITDS establishes connections to fields in databases

func initDS() (*dataSources, error) {
	log.Printf("Initializing data sources\n")

	// load env variables -
	// top level(main package)
	// helper function ,
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")
	pgDB := os.Getenv("PG_DB")
	pgSSL := os.Getenv("PG_SSL")

	pgConnString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", pgHost, pgPort, pgUser, pgPassword, pgDB, pgSSL)

	log.Printf("Connecting to postgresql\n")
	db, err := sqlx.Connect("postgres", pgConnString)

	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// verify database connection is working
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	return &dataSources{
		DB: db,
	}, nil
}

func (d *dataSources) close() error {
	if err := d.DB.Close(); err != nil {
		return fmt.Errorf(" error closing Postgres Postgresql: %w", err)
	}

	return nil
}