package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"link-service/internal/config"
	"link-service/internal/errs"
	"link-service/internal/models"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type LinkRepository struct {
	db *sql.DB
}

func NewLinkRepository(dbConfig *config.DbConfig) (*LinkRepository, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name)

	db, err := sql.Open("postgres", psqlInfo)

	// Ограничиваем количество открытых соединений
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)
	// Ограничиваем количество простаивающих соединений
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)
	// Устанавливаем время жизни соединения
	db.SetConnMaxLifetime(time.Duration(dbConfig.ConnMaxLifetimeInMin) * time.Minute)

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

func (l *LinkRepository) SaveUrl(ctx context.Context, link string, shortCode string) error {

	_, err := l.db.ExecContext(ctx, "INSERT INTO links (original_url, short_code, created_at, visits) VALUES ($1,$2,NOW(), 0)", link, shortCode)
	if err != nil {
		log.Print("save url error: ", err)
		return err
	}

	return nil

}

func (l *LinkRepository) FindLinkByShortCode(ctx context.Context, shortCode string) (*models.LinkDto, error) {

	row := l.db.QueryRowContext(ctx, "SELECT original_url, visits FROM links WHERE short_code = $1", shortCode)
	lDto := models.LinkDto{}
	err := row.Scan(&lDto.Url, &lDto.Visits)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NewNotFoundLinkError()
		} else {
			return nil, err
		}
	}

	return &lDto, err

}

func (l *LinkRepository) IncrementVisit(ctx context.Context, shortCode string) (int32, error) {

	row := l.db.QueryRowContext(ctx, "UPDATE links SET visits = visits+1 WHERE short_code = $1 RETURNING visits", shortCode)
	var visits int32
	err := row.Scan(&visits)

	if err != nil {
		log.Print("error: ", err)
		return visits, err
	}
	return visits, nil
}

func (l *LinkRepository) FindLinkStatsByShortCode(ctx context.Context, shortCode string) (models.LinkDto, error) {
	row := l.db.QueryRowContext(ctx, "SELECT original_url, short_code, visits, created_at FROM links WHERE short_code = $1", shortCode)
	lDto := models.LinkDto{}
	err := row.Scan(&lDto.Url, &lDto.ShortCode, &lDto.Visits, &lDto.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return lDto, errs.NewNotFoundLinkError()
		} else {
			return lDto, err
		}
	}

	return lDto, err

}

func (l *LinkRepository) DeleteByShortCode(ctx context.Context, shortCode string) error {
	_, err := l.db.ExecContext(ctx, "DELETE FROM links WHERE short_code = $1", shortCode)
	if err != nil {
		log.Print("error: ", err)
		return err
	}
	return nil
}

func (l *LinkRepository) GetByPage(ctx context.Context, limit, offset int) ([]models.LinkDto, error) {

	rows, err := l.db.QueryContext(ctx, "SELECT original_url, short_code, visits, created_at FROM links ORDER BY created_at DESC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		log.Print("error: ", err)
		return make([]models.LinkDto, 0), err
	}

	defer rows.Close()

	var response []models.LinkDto

	for rows.Next() {

		lDto := models.LinkDto{}

		if err := rows.Scan(&lDto.Url, &lDto.ShortCode, &lDto.Visits, &lDto.CreatedAt); err != nil {
			log.Print(err)
			return make([]models.LinkDto, 0), err
		}

		response = append(response, lDto)

	}
	return response, nil

}
