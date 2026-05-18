package services

import (
	"context"
	"link-service/internal/cache"
	"link-service/internal/models"
	"link-service/internal/repository"
	"log"
	"math/rand"
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
	visits, err := ls.linkRepository.IncrementVisit(ctx, shortCode)
	if err != nil {
		log.Printf("failed to increment visits: %v\n", err)
		return
	}

	cacheDto := models.CacheDto{
		ShortCode: shortCode,
		Url:       linkDto.Url,
		Visits:    visits}

	if err := ls.linkCache.PutLinkInfo(ctx, cacheDto); err != nil {
		log.Printf("failed to set cache: %v\n", err)
	}

}

func (ls *LinkService) FindLinkByShortCode(ctx context.Context, shortCode string) (*models.LinkDto, error) {

	linkDto, err := ls.linkCache.GetLinkInfo(ctx, shortCode)

	if err == nil && linkDto != nil {
		linkDto.Visits++
		// Асинхронно увеличиваем счетчик в базе и в кеше
		linkDtoCopy := *linkDto

		go func() {
			timeout := 5 * time.Second
			updateCtx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			ls.updateData(updateCtx, linkDtoCopy, shortCode)

		}()
		return linkDto, nil
	}

	if err != nil {
		log.Printf("failed to get cache: %v\n", err)
	}

	linkDto, err = ls.linkRepository.FindLinkByShortCode(ctx, shortCode)

	if err != nil {
		return nil, err
	}
	linkDto.Visits++
	ls.updateData(ctx, *linkDto, shortCode)

	return linkDto, err
}

func (ls *LinkService) FindLinkStatsByShortCode(ctx context.Context, shortCode string) (models.LinkDto, error) {
	linkDto, err := ls.linkRepository.FindLinkStatsByShortCode(ctx, shortCode)
	return linkDto, err
}

func (ls *LinkService) DeleteByShortCode(ctx context.Context, shortCode string) error {

	if err := ls.linkCache.DeleteLinkInfo(ctx, shortCode); err != nil {
		log.Printf("failed to delete cache: %v\n", err)
	}
	return ls.linkRepository.DeleteByShortCode(ctx, shortCode)
}

func (ls *LinkService) GetByPage(ctx context.Context, limit, offset int) ([]models.LinkDto, error) {
	return ls.linkRepository.GetByPage(ctx, limit, offset)
}

func CreateRandomString(size int) string {

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, size)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)

}
