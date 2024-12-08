package config

import (
	"fmt"
	"os"
)

func GetEnvOrDie(key string) string {
	value := os.Getenv(key)

	if value == "" {
		err := fmt.Errorf("missing environment variable %s", key)
		panic(err)
	}

	return value
}

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}

	Server struct {
		Port string
	}
}

func NewConfig() *Config {
	return &Config{
		Database: struct {
			Host     string
			Port     string
			User     string
			Password string
			Name     string
		}{
			Host:     GetEnvOrDie("DB_HOST"),
			Port:     GetEnvOrDie("DB_PORT"),
			User:     GetEnvOrDie("DB_USER"),
			Password: GetEnvOrDie("DB_PASSWORD"),
			Name:     GetEnvOrDie("DB_NAME"),
		},

		Server: struct {
			Port string
		}{
			Port: GetEnvOrDie("SERVER_PORT"),
		},
	}
}
