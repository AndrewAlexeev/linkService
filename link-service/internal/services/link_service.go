package services

import (
	"context"
	"link-service/internal/cache"
	"link-service/internal/models"
	"link-service/internal/repository"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type LinkService struct {
	linkRepository *repository.LinkRepository
	linkCache      *cache.LinkCache
}

func NewLinkService(repository *repository.LinkRepository, cache *cache.LinkCache) *LinkService {
	return &LinkService{linkRepository: repository, linkCache: cache}
}

func (ls *LinkService) SaveLink(ctx context.Context, url string) (string, error) {
	shortCode := CreateRandomString(10)
	return shortCode, ls.linkRepository.SaveUrl(ctx, url, shortCode)
}

func (ls *LinkService) updateData(ctx context.Context, linkDto models.LinkDto, shortCode string) {
	err, visits := ls.linkRepository.IncrementVisit(ctx, shortCode)
	if err != nil {
		// Логируем ошибку, но не блокируем ответ
		log.Printf("failed to increment visits: %v\n", err)
		return
	}

	cacheDto := models.CacheDto{
		ShortCode: shortCode,
		Url:       linkDto.Url,
		Visits:    visits}

	if err := ls.linkCache.PutLinkInfo(ctx, cacheDto); err != nil {
		// Логируем ошибку, но не блокируем ответ
		log.Printf("failed to set cache: %v\n", err)
	}

}

func (ls *LinkService) FindLinkByShortCode(ctx context.Context, shortCode string) (*models.LinkDto, error) {

	linkDto, err := ls.linkCache.GetLinkInfo(ctx, shortCode)

	if err == nil && linkDto != nil {
		linkDto.Visits++
		// Асинхронно увеличиваем счетчик в базе и в кеше
		timeout := 5 * time.Second
		uodateCtx, _ := context.WithTimeout(context.Background(), timeout)
		go ls.updateData(uodateCtx, *linkDto, shortCode)
		return linkDto, nil
	}

	linkDto, err = ls.linkRepository.FindLinkByShortCode(ctx, shortCode)
	linkDto.Visits++
	ls.updateData(ctx, *linkDto, shortCode)

	return linkDto, err
}

func (ls *LinkService) FindLinkStatsByShortCode(ctx context.Context, shortCode string) (models.LinkDto, error) {
	linkDto, err := ls.linkRepository.FindLinkStatsByShortCode(ctx, shortCode)
	return linkDto, err
}

func (ls *LinkService) DeleteByShortCode(ctx context.Context, shortCode string) error {
	ls.linkCache.DeleteLinkInfo(ctx, shortCode)
	return ls.linkRepository.DeleteByShortCode(ctx, shortCode)
}

func (ls *LinkService) GetByPage(ctx context.Context, limit, offset string) ([]models.LinkDto, error) {
	limitInt, err1 := strconv.Atoi(limit)

	if err1 != nil {
		return make([]models.LinkDto, 0), err1
	}

	offsetInt, err2 := strconv.Atoi(offset)

	if err2 != nil {
		return make([]models.LinkDto, 0), err1
	}

	return ls.linkRepository.GetByPage(ctx, limitInt, offsetInt)
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
