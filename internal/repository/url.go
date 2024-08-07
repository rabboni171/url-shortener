package repository

import (
	"context"
	"fmt"

	"github.com/rabboni171/url-shortener/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) IURLRepository {
	return &RedisRepo{
		client: client,
	}
}

func (r *RedisRepo) Save(shortURL string, originalURL string) error {
	err := r.client.Set(context.Background(), shortURL, originalURL, 0).Err()
	if err != nil {
		return fmt.Errorf("%w: %v", errors.ErrRedisError, err)
	}
	return nil
}

func (r *RedisRepo) Get(shortURL string) (string, error) {
	url, err := r.client.Get(context.Background(), shortURL).Result()
	if err == redis.Nil {
		return "", errors.ErrURLNotFound
	}
	return url, err
}
