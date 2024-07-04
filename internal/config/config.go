package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
}

type ServerConfig struct {
	Port            int
	Timeouts        Timeouts
	HealthCheckPort int `mapstructure:"health_check_port"`
}

type Timeouts struct {
	Write time.Duration
	Read  time.Duration
	Idle  time.Duration
}

func LoadConfig(configPaths []string) *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	return &c
}
