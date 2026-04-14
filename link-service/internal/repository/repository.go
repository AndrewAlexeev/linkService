package repository

import (
	"database/sql"
	"fmt"
	"link-service/internal/config"
	"link-service/internal/models"
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

func (l LinkRepository) SaveUrl(link string, shortCode string) error {

	_, err := l.db.Exec("insert into links (original_url, short_code, created_at, visits) values ($1,$2,NOW(), 0)", link, shortCode)
	if err != nil {
		log.Print("error: ", err)
		return err
	}

	return nil

}

func (l LinkRepository) FindLinkByShortCode(shortCode string) (*models.LinkDto, error) {

	row := l.db.QueryRow("select original_url, visits from links where short_code = $1", shortCode)
	lDto := models.LinkDto{}
	err := row.Scan(&lDto.Url, &lDto.Visits)
	return &lDto, err

}
