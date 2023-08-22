package config

import (
	"fmt"
	"os"
)

type DBConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"hostname"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Schema   string `mapstructure:"schema"`
}

func DefaultDB() DBConfig {

	return DBConfig{
		Username: "root",
		Password: "secret",
		Host:     "localhost",
		Port:     5432,
		Database: "products",
		Schema:   "products",
	}
}

func (Conf *DBConfig) loadVarsFromEnv() {

	Username, ok := os.LookupEnv("POSTGRES_USER")
	if ok {
		Conf.Username = Username
	}

	Password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if ok {
		Conf.Password = Password
	}

	Host, ok := os.LookupEnv("POSTGRES_HOSTNAME")
	if ok {
		Conf.Host = Host
	}

	Database, ok := os.LookupEnv("POSTGRES_DB")
	if ok {
		Conf.Database = Database
	}

	PortStr, ok := os.LookupEnv("POSTGRES_PORT")
	if ok {
		Conf.Port = stringToInt(PortStr)
	}

	schema, ok := os.LookupEnv("POSTGRES_SCHEMA")
	if ok {
		Conf.Schema = schema
	}
}

// GetDSN DSN: user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10
func (Conf *DBConfig) GetDSN() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable pool_max_conns=10", Conf.Username, Conf.Password, Conf.Host, Conf.Port, Conf.Database)
}
