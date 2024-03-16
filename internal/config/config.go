package config

import (
	"fmt"
	"log"
	"ms-template-go/internal/database"
	"ms-template-go/pkg/otel"
	"ms-template-go/pkg/utils"
	"os"

	"github.com/spf13/viper"
)

type GlobalConfig struct {
	DB   database.Config `mapstructure:"database"`
	Otel otel.Config     `mapstructure:"open_telemetry"`

	// Debug if true shows more logs and info
	Debug bool `mapstructure:"debug_active"`

	// GrpcPort is the port the grpc server will listen to
	GrpcPort int `mapstructure:"grpc_port"`
}

// Default populates GlobalConfig with default values
func (C *GlobalConfig) Default() {

	C.GrpcPort = 4445
	C.Debug = true
	C.DB.Default()
	C.Otel.Default()
}

// Make sure we conform to ConfigInterface
var _ ConfigInterface = (*GlobalConfig)(nil)
var _ ConfigInterface = (*database.Config)(nil)
var _ ConfigInterface = (*otel.Config)(nil)

// loadVarsFromYaml reads vars from yaml file and overrides current values if new value found
func (C *GlobalConfig) loadVarsFromYaml() {
	viper.SetConfigName("config.yml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/configs")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Config file not found!")
	}

	if err := viper.Unmarshal(C); err != nil {
		log.Println("Viper failed to unmarshall config: ", err)
	}
}

// LoadVarsFromEnvVar looks for GlobalConfig values in ENV VAR, if found, overrides previous values
func (C *GlobalConfig) LoadVarsFromEnv() {

	// Web
	WebPortStr, ok := os.LookupEnv("GRPC_PORT")
	if ok {

		var err error
		C.GrpcPort, err = utils.StringToInt(WebPortStr)
		if err != nil {
			log.Println(err)
		}
	}

	DebugStr, ok := os.LookupEnv("DEBUG_ENABLED")
	if ok {

		var err error
		C.Debug, err = utils.StringToBool(DebugStr)
		if err != nil {
			log.Println(err)
		}
	}

	C.DB.LoadVarsFromEnv()
	C.Otel.LoadVarsFromEnv()
}

// GetListenAddress returns the address string the socket will be listening to
func (C *GlobalConfig) GetListenAddress() string {
	return fmt.Sprintf(":%d", C.GrpcPort)
}

var C GlobalConfig

func LoadConfiguration() {

	C.Default()
	C.loadVarsFromYaml()
	C.LoadVarsFromEnv()

	if C.Debug {
		log.Println("LOADED ENV: ", C)
	}
}
