package database

import (
	"context"
	"log"

	"cdlab.cdlan.net/cdlan/uservices/ms-template/internal/config"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type State struct {
	DBPool *pgxpool.Pool
}

var DBState State

func Init() {
	// Start a new span for DB initialization
	_, span := config.C.Otel.NewSpan(context.Background(), "DB.Init")
	if span != nil {
		defer span.End()
	}

	cfg, err := pgxpool.ParseConfig(config.C.DB.GetDSN())
	if err != nil {
		span.RecordError(err)
		log.Printf("Error initializing database: %v", err)
		return
	}

	if config.C.Otel.Enabled {
		cfg.ConnConfig.Tracer = otelpgx.NewTracer()
	}

	dbPool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		span.RecordError(err)
		log.Printf("Error initializing database: %v", err)
		return
	}
	DBState.DBPool = dbPool
}

func Close() {
	_, span := config.C.Otel.NewSpan(context.Background(), "DB.Close")
	if span != nil {
		defer span.End()
	}

	DBState.DBPool.Close()
}
