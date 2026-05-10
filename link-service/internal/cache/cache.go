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

func (lc *LinkCache) PutLinkInfo(ctx context.Context, dto models.CacheDto) error {

	value, err := json.Marshal(dto)
	if err != nil {
		log.Printf("failed to marshal link: %s", err)
		return err
	}

	key := dto.ShortCode
	maxRetries := 5
	// Начальная задержка между попытками
	retryBackoff := 10 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		// CAS с использованием WATCH
		err = lc.client.Watch(ctx, func(tx *redis.Tx) error {
			// 1. Получаем текущее значение
			val, err := tx.Get(ctx, key).Result()
			if err != nil && err != redis.Nil {
				return err
			}
			if val == "" {
				_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
					pipe.Set(ctx, key, value, 0)
					return nil
				})

				if err != nil {
					return err
				}
				return nil
			}
			var linkDto models.LinkDto

			if err = json.Unmarshal([]byte(val), &linkDto); err != nil {
				return err
			}

			// 2. если в кеше уже более актуальное значение счетчика, пропускаем обновление кеша
			if linkDto.Visits > dto.Visits {
				return nil
			}

			// 3. Выполняем изменения в транзакции
			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, key, value, 0)
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
		log.Printf("failed to get link in cache: %s", err)
		return nil, err
	}

	var linkDto models.LinkDto

	if err = json.Unmarshal([]byte(val), &linkDto); err != nil {
		return nil, err
	}

	return &linkDto, nil

}

func (lc *LinkCache) DeleteLinkInfo(ctx context.Context, key string) {
	lc.client.Del(ctx, key)
}
