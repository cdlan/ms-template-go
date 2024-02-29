package config

// ConfigInterface is the interface that must be implemented by estensions for their config to be added to the main one
type ConfigInterface interface {
	Default()
	LoadVarsFromEnv()
}