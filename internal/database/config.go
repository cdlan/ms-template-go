package database

import (
	"ms-template-go/pkg/utils"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
	Host          string `mapstructure:"hostname"`
	Port          int    `mapstructure:"port"`
	Database      string `mapstructure:"database"`
	Schema        string `mapstructure:"schema"`
	MigrationPath string `mapstructure:"migration_path"`
}

func (C *Config) Default() {

	C.Username = "root"
	C.Password = "secret"
	C.Host = "localhost"
	C.Port = 5432
	C.Database = "mistra"
	C.Schema = "users"
	C.MigrationPath = "internal/migrations/sql"
}

func (C *Config) LoadVarsFromEnv() {

	Username, ok := os.LookupEnv("POSTGRES_USER")
	if ok {
		C.Username = Username
	}

	Password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if ok {
		C.Password = Password
	}

	Host, ok := os.LookupEnv("POSTGRES_HOSTNAME")
	if ok {
		C.Host = Host
	}

	Database, ok := os.LookupEnv("POSTGRES_DB")
	if ok {
		C.Database = Database
	}

	PortStr, ok := os.LookupEnv("POSTGRES_PORT")
	if ok {

		var err error
		C.Port, err = utils.StringToInt(PortStr)
		if err != nil {
			log.Println(err)
		}
	}

	schema, ok := os.LookupEnv("POSTGRES_SCHEMA")
	if ok {
		C.Schema = schema
	}

	migrationPath, ok := os.LookupEnv("MIGRATION_PATH")
	if ok {
		C.MigrationPath = migrationPath
	}
}

// GetDSN DSN: user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10 options=--search_path=
func (C *Config) GetDSN() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable pool_max_conns=10 options=--search_path=%s", C.Username, C.Password, C.Host, C.Port, C.Database, C.Schema)
}

// GetURL returns a connection string like postgres://user:secret@localhost:5432/database?sslmode=disable&search_path=your_schema_name
func (C *Config) GetURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&search_path=%s", C.Username, C.Password, C.Host, C.Port, C.Database, C.Schema)
	//return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", C.Username, C.Password, C.Host, C.Port, C.Database)
}
