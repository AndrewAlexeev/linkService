package repository

import (
	"database/sql"
	"fmt"
	"link-service/internal/config"
	"log"
)

type LinkRepository struct {
	db *sql.DB
}

func Init(dbConfig config.DbConfig) (*LinkRepository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("failed to connect to database: %w", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to ping database: %w", err)
		return nil, err
	}

	return &LinkRepository{
		db: db}, nil

}

func (l LinkRepository) SaveUrl(link string, shortLink string) error {

	_, err := l.db.Exec("insert into links (original_url, short_code, created_at, visits) values ($1,$2,NOW(), 0)",
		link, shortLink)
	if err != nil {
		log.Fatal("error: ", err)
		return err
	}

	return nil

}
