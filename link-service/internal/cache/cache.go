package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"link-service/internal/config"
	"link-service/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type LinkCache struct {
	client *redis.Client
	ttl    time.Duration
}

func InitLinkCache(redisConfig config.RedisConfig) (*LinkCache, error) {
	rdc := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdc.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	linkCache := LinkCache{
		client: rdc, ttl: time.Duration(redisConfig.CacheTTL) * time.Second}

	return &linkCache, nil

}

func (lc *LinkCache) PutLinkInfo(ctx context.Context, dto models.CacheDto) error {

	value, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("failed to marshal link: %w", err)
	}

	key := dto.ShortCode
	maxRetries := 5
	// Начальная задержка между попытками
	retryBackoff := 10 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		// CAS с использованием WATCH, чтобы решать рейскондишин при обновлении кеша
		err = lc.client.Watch(ctx, func(tx *redis.Tx) error {

			// 1. Получаем текущее значение
			val, err := tx.Get(ctx, key).Result()
			if err != nil && err != redis.Nil {
				return err
			}

			if err != redis.Nil {
				var cacheDto models.CacheDto

				if err = json.Unmarshal([]byte(val), &cacheDto); err != nil {
					return fmt.Errorf("failed to unmarshal existing value: %w", err)

				}

				// 2. если в кеше уже более актуальное значение счетчика, пропускаем обновление кеша
				if cacheDto.Visits > dto.Visits {
					return nil
				}
			}

			// 3. Выполняем изменения в транзакции
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, key, value, lc.ttl)
				return nil
			})
			return err
		}, key)

		if err != nil {
			// Ждем перед следующей попыткой
			time.Sleep(retryBackoff)
			retryBackoff *= 2 // Увеличиваем задержку
			continue
		}
		return nil

	}

	return err

}

func (lc *LinkCache) GetLinkInfo(ctx context.Context, key string) (*models.LinkDto, error) {

	val, err := lc.client.Get(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get link by key %s from cache: %w", key, err)
	}

	var cacheDto models.CacheDto

	if err = json.Unmarshal([]byte(val), &cacheDto); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return &models.LinkDto{
		Url:       cacheDto.Url,
		ShortCode: cacheDto.ShortCode,
		Visits:    cacheDto.Visits,
	}, nil

}

func (lc *LinkCache) DeleteLinkInfo(ctx context.Context, key string) error {
	err := lc.client.Del(ctx, key).Err()

	if err != nil {
		return fmt.Errorf("failed to delete key %s from cache: %w", key, err)
	}
	return nil
}
