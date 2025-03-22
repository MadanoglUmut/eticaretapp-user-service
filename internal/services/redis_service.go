package services

import (
	"context"
	"fmt"
	"time"
)

type redisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	LLen(ctx context.Context, key string) (int64, error)
	RPop(ctx context.Context, key string) (string, error)
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	Type(ctx context.Context, key string) (string, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
}

type RedisService struct {
	redisClient redisClient
}

func NewRedisService(redisClient redisClient) *RedisService {

	return &RedisService{
		redisClient: redisClient,
	}

}

func (s *RedisService) GetTokens(userId int) ([]string, error) {
	redisKey := fmt.Sprintf("user:%d:tokens", userId)

	tokens, err := s.redisClient.LRange(context.Background(), redisKey, 0, -1)
	if err != nil {
		return nil, fmt.Errorf("token'lar alınamadı: %v", err)
	}

	return tokens, nil
}

func (s *RedisService) Set(userId int, token string) error {
	redisKey := fmt.Sprintf("user:%d:tokens", userId)

	keyType, err := s.redisClient.Type(context.Background(), redisKey)
	if err != nil {
		return fmt.Errorf("anahtar türü alınamadı: %v", err)
	}

	if keyType != "list" {
		_, err := s.redisClient.Del(context.Background(), redisKey)
		if err != nil {
			return fmt.Errorf("anahtar silinemedi: %v", err)
		}
	}

	_, err = s.redisClient.LPush(context.Background(), redisKey, token)
	if err != nil {
		return fmt.Errorf("token Redis'e eklenemedi: %v", err)
	}

	listLength, err := s.redisClient.LLen(context.Background(), redisKey)
	if err != nil {
		return fmt.Errorf("liste uzunluğu alınamadı: %v", err)
	}

	if listLength > 2 {
		_, err = s.redisClient.RPop(context.Background(), redisKey)
		if err != nil {
			return fmt.Errorf("en eski token silinemedi: %v", err)
		}
	}

	_, err = s.redisClient.Expire(context.Background(), redisKey, 24*time.Hour)
	if err != nil {
		return fmt.Errorf("token TTL ayarlanamadı: %v", err)
	}

	return nil
}
