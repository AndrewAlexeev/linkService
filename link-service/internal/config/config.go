package config

import (
	"fmt"
	"os"
	"strconv"
)

type DbConfig struct {
	Host     string
	Port     string
	User     string
	Name     string
	Password string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	CacheTTL int
}

type Config struct {
	Port string
}

func InitDbConfig() DbConfig {
	config := DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Name:     os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD")}
	return config

}

func InitRedisConfig() (*RedisConfig, error) {

	dbStr := os.Getenv("REDIS_DB")

	db, err := strconv.Atoi(dbStr)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return nil, err
	}

	cacheTTLStr := os.Getenv("CACHE_TTL")

	cacheTTL, err := strconv.Atoi(cacheTTLStr)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return nil, err
	}

	config := RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
		CacheTTL: cacheTTL}
	return &config, nil

}

func InitConfig() Config {
	config := Config{
		Port: os.Getenv("PORT")}
	return config

}
