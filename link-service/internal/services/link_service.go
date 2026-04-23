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

func (ls *LinkService) FindLinkByShortCode(ctx context.Context, shortCode string) (*models.LinkDto, error) {
	//TODO тут скорее всего выполняется как go так и ниже стоящий код (тот, что если не нашли в кеше, надо исправить, чтобы нижний код исполнялся ТОЛЬКО если не нашли в кешеы)
	linkDto, err := ls.linkCache.GetLinkInfo(ctx, shortCode)
	if err == nil && linkDto != nil {
		// Асинхронно увеличиваем счетчик в базе
		go func() {
			if err := ls.linkRepository.IncrementVisit(context.Background(), shortCode); err != nil {
				// Логируем ошибку, но не блокируем ответ
				log.Printf("failed to increment visits: %v\n", err)
			}
		}()
		linkDto.Visits++
		return linkDto, nil
	}

	linkDto, err = ls.linkRepository.FindLinkByShortCode(ctx, shortCode)
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	go ls.linkRepository.IncrementVisit(context.Background(), shortCode)
	linkDto.Visits = linkDto.Visits + 1

	cacheDto := models.CacheDto{
		ShortCode: shortCode,
		Url:       linkDto.Url,
		Visits:    linkDto.Visits}

	if err := ls.linkCache.PutLinkInfo(ctx, cacheDto); err != nil {
		// Логируем ошибку, но не блокируем ответ
		log.Printf("failed to set cache: %v\n", err)
	}

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
