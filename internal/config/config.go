package config

import (
	"fmt"
	"log"
	"os"

	"cdlab.cdlan.net/cdlan/uservices/ms-template/pkg/otel"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	DB       DBConfig    `mapstructure:"database"`
	GrpcPort int         `mapstructure:"grpc_port"`
	Debug    bool        `mapstructure:"debug_active"`
	Otel     otel.Config `mapstructure:"open_telemetry"`
}

// Default generates a GlobalConfig with default values
func Default() GlobalConfig {

	return GlobalConfig{
		GrpcPort: 4445,
		Debug:    false,
		DB:       DefaultDB(),
		Otel:     otel.Default(),
	}
}

// loadVarsFromYaml reads vars from yaml file and overrides current values if new value found
func (C *GlobalConfig) loadVarsFromYaml() {
	viper.SetConfigName("config.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Config file not found!")
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Println(err)
	}
}

// loadVarsFromEnv looks for GlobalConfig values in ENV VAR, if found, overrides previous values
func (C *GlobalConfig) loadVarsFromEnv() {

	// Web
	WebPortStr, ok := os.LookupEnv("GRPC_PORT")
	if ok {

		C.GrpcPort = stringToInt(WebPortStr)
	}

	DebugStr, ok := os.LookupEnv("DEBUG_ENABLED")
	if ok {

		C.Debug = stringToBool(DebugStr)
	}

	C.DB.loadVarsFromEnv()
	C.Otel.LoadVarsFromEnv()
}

// GetListenAddress returns the address string the socket will be listening to
func (C *GlobalConfig) GetListenAddress() string {
	return fmt.Sprintf(":%d", C.GrpcPort)
}

var C GlobalConfig

func LoadConfiguration() {

	C = Default()
	C.loadVarsFromYaml()
	C.loadVarsFromEnv()

	if C.Debug {
		log.Println("LOADED ENV: ", C)
	}
}
