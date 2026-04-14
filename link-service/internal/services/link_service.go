package services

import (
	"link-service/internal/models"
	"link-service/internal/repository"
	"math/rand"
	"time"
)

type LinkService struct {
	LinkRepository *repository.LinkRepository
}

func (l LinkService) SaveUrl(url string) (string, error) {
	shortCode := CreateRandomString(10)
	return shortCode, l.LinkRepository.SaveUrl(url, shortCode)
}

func (l LinkService) FindLinkByShortCode(shortCode string) (*models.LinkDto, error) {
	return l.LinkRepository.FindLinkByShortCode(shortCode)
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
