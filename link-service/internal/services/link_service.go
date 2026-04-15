package services

import (
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

func (ls *LinkService) SaveUrl(url string) (string, error) {
	shortCode := CreateRandomString(10)
	return shortCode, ls.linkRepository.SaveUrl(url, shortCode)
}

func (ls *LinkService) FindLinkByShortCode(shortCode string) (models.LinkDto, error) {
	linkDto, err := ls.linkRepository.FindLinkByShortCode(shortCode)
	go ls.linkRepository.IncrementVisit(shortCode)
	linkDto.Visits = linkDto.Visits + 1
	return linkDto, err
}

func (ls *LinkService) FindLinkStatsByShortCode(shortCode string) (models.LinkDto, error) {
	linkDto, err := ls.linkRepository.FindLinkStatsByShortCode(shortCode)
	return linkDto, err
}

func (ls *LinkService) DeleteByShortCode(shortCode string) error {
	return ls.linkRepository.DeleteByShortCode(shortCode)
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
