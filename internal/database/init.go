package database

import (
	"context"
	"log"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Init start up the DB connections
func (C *Config) Init(OpenTelemetryEnabled bool) {

	cfg, err := pgxpool.ParseConfig(C.GetDSN())
	if err != nil {
		log.Printf("Error initializing database: %v", err)
		return
	}

	// if otel enabled add sql tracer
	if OpenTelemetryEnabled {
		cfg.ConnConfig.Tracer = otelpgx.NewTracer()
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Printf("Error initializing database: %v", err)
		return
	}

	DBState.DBPool = dbPool
}

// Close closes DB connection
func Close() {
	DBState.DBPool.Close()
}
