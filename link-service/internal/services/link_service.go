package services

import (
	"context"
	"link-service/internal/models"
	"link-service/internal/repository"
	"math/rand"
	"time"
)

// type LinkRepository interface {
// 	SaveUrl(url, shortCode string) error
// 	FindLinkByShortCode(shortCode string) (models.LinkDto, error)
// 	FindLinkStatsByShortCode(shortCode string) (models.LinkDto, error)
// 	IncrementVisit(shortCode string) error
// 	DeleteByShortCode(shortCode string) error
// }

type LinkService struct {
	linkRepository *repository.LinkRepository
}

func NewLinkService(repository *repository.LinkRepository) *LinkService {
	return &LinkService{linkRepository: repository}
}

func (ls *LinkService) SaveUrl(ctx context.Context, url string) (string, error) {
	shortCode := CreateRandomString(10)
	return shortCode, ls.linkRepository.SaveUrl(ctx, url, shortCode)
}

func (ls *LinkService) FindLinkByShortCode(ctx context.Context, shortCode string) (models.LinkDto, error) {
	linkDto, err := ls.linkRepository.FindLinkByShortCode(ctx, shortCode)
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	go ls.linkRepository.IncrementVisit(context.Background(), shortCode)
	linkDto.Visits = linkDto.Visits + 1
	return linkDto, err
}

func (ls *LinkService) FindLinkStatsByShortCode(ctx context.Context, shortCode string) (models.LinkDto, error) {
	linkDto, err := ls.linkRepository.FindLinkStatsByShortCode(ctx, shortCode)
	return linkDto, err
}

func (ls *LinkService) DeleteByShortCode(ctx context.Context, shortCode string) error {
	return ls.linkRepository.DeleteByShortCode(ctx, shortCode)
}

func CreateRandomString(size int) string {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano()) // Инициализация сида

	s := make([]rune, size)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)

}
