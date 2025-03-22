package redis

import (
	"UserService/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	LLen(ctx context.Context, key string) *redis.IntCmd
	RPop(ctx context.Context, key string) *redis.StringCmd
	LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd
	LPush(ctx context.Context, key string, values ...interface{}) *redis.IntCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Type(ctx context.Context, key string) *redis.StatusCmd
	Ping(ctx context.Context) *redis.StatusCmd
}

type Redis struct {
	client RedisClient
	ctx    context.Context
}

func NewRedis(client RedisClient, ctx context.Context) *Redis {

	return &Redis{
		client: client,
		ctx:    ctx,
	}

}

func NewClient(redisHost, redisPort, redisPassword string) RedisClient {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       0,
	})
}

func (r *Redis) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {

	if key == "" {
		return models.ErrKeyIsEmpty

	}

	jsonValue, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, jsonValue, expiration).Err()

}

func (r *Redis) LLen(ctx context.Context, key string) (int64, error) {
	if key == "" {
		return 0, models.ErrKeyIsEmpty
	}
	return r.client.LLen(ctx, key).Result()
}

func (r *Redis) RPop(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", models.ErrKeyIsEmpty
	}
	return r.client.RPop(ctx, key).Result()
}

func (r *Redis) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	if key == "" {
		return nil, models.ErrKeyIsEmpty
	}
	return r.client.LRange(ctx, key, start, stop).Result()
}

func (r *Redis) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	if key == "" {
		return 0, models.ErrKeyIsEmpty
	}
	return r.client.LPush(ctx, key, values...).Result()
}

func (r *Redis) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	if key == "" {
		return false, models.ErrKeyIsEmpty
	}
	return r.client.Expire(ctx, key, expiration).Result()
}

func (r *Redis) Del(ctx context.Context, keys ...string) (int64, error) {
	if len(keys) == 0 {
		return 0, models.ErrKeyIsEmpty
	}

	return r.client.Del(ctx, keys...).Result()
}

func (r *Redis) Type(ctx context.Context, key string) (string, error) {
	if key == "" {
		return "", models.ErrKeyIsEmpty
	}

	return r.client.Type(ctx, key).Result()
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
