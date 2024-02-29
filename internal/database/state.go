package database

import "github.com/jackc/pgx/v5/pgxpool"

type State struct {
	DBPool *pgxpool.Pool
}

// DBState is the global var that holds the database pool
var DBState State
