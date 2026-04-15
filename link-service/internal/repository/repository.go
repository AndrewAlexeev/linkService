package repository

import (
	"database/sql"
	"fmt"
	"link-service/internal/config"
	"link-service/internal/models"
	"log"

	_ "github.com/lib/pq"
)

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(dbConfig config.DbConfig) (*LinkRepository, error) {
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

func (l *LinkRepository) SaveUrl(link string, shortCode string) error {

	_, err := l.db.Exec("INSERT INTO links (original_url, short_code, created_at, visits) VALUES ($1,$2,NOW(), 0)", link, shortCode)
	if err != nil {
		log.Print("error: ", err)
		return err
	}

	return nil

}

func (l *LinkRepository) FindLinkByShortCode(shortCode string) (models.LinkDto, error) {

	row := l.db.QueryRow("SELECT original_url, visits FROM links WHERE short_code = $1", shortCode)
	lDto := models.LinkDto{}
	err := row.Scan(&lDto.Url, &lDto.Visits)
	return lDto, err

}

func (l *LinkRepository) IncrementVisit(shortCode string) error {

	_, err := l.db.Exec("UPDATE links SET visits = visits+1 WHERE short_code = $1", shortCode)
	if err != nil {
		log.Print("error: ", err)
		return err
	}
	return nil
}

func (l *LinkRepository) FindLinkStatsByShortCode(shortCode string) (models.LinkDto, error) {
	row := l.db.QueryRow("SELECT original_url, short_code, visits, created_at FROM links WHERE short_code = $1", shortCode)
	lDto := models.LinkDto{}
	err := row.Scan(&lDto.Url, &lDto.ShortCode, &lDto.Visits, &lDto.CreatedAt)
	return lDto, err

}

func (l *LinkRepository) DeleteByShortCode(shortCode string) error {
	_, err := l.db.Exec("DELETE FROM links WHERE short_code = $1", shortCode)
	if err != nil {
		log.Print("error: ", err)
		return err
	}
	return nil
}
