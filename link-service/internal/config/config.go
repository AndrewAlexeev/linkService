package config

import (
	"os"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Name     string
	Password string
}

type Config struct {
	Port string
}

func InitDbConfig() DbConfig {
	config := DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		Name:     os.Getenv("DB_NAME")}
	return config

}

func InitConfig() Config {
	config := Config{
		Port: os.Getenv("PORT")}
	return config

}
