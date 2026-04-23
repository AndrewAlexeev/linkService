package cache

import (
	"context"
	"encoding/json"
	"link-service/internal/config"
	"link-service/internal/models"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type LinkCache struct {
	client *redis.Client
	ttl    time.Duration
}

func InitLinkCache(redisConfig config.RedisConfig) *LinkCache {
	rdc := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB})
	time.Duration.Seconds(132)

	linkCache := LinkCache{
		client: rdc, ttl: time.Duration(redisConfig.CacheTTL) * time.Second}

	return &linkCache

}

func (lc LinkCache) PutLinkInfo(ctx context.Context, dto models.CacheDto) error {

	data, err := json.Marshal(dto)
	if err != nil {
		log.Printf("failed to marshal link: %s", err)
		return err
	}

	err = lc.client.Set(ctx, dto.ShortCode, data, lc.ttl).Err()
	if err != nil {
		log.Printf("failed to put link to cache: %s", err)
		return err
	}
	return nil

}

func (rc LinkCache) GetLinkInfo(ctx context.Context, key string) (*models.LinkDto, error) {

	val, err := rc.client.Get(ctx, key).Bytes()
	if err != nil {
		log.Printf("failed to get link in cache: %s", err)
		return nil, err
	}

	var linkDto *models.LinkDto
	err = json.Unmarshal(val, linkDto)

	return linkDto, nil

}

func (rc LinkCache) DeleteLinkInfo(ctx context.Context, key string) {
	rc.client.Del(ctx, key)
}
