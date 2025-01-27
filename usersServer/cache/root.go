package cache

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
	"time"

	"github.com/redis/go-redis/v9"
)

type ICache[T any] interface {
	Handle(ctx context.Context, key string, fn func() (*T, error)) (*T, error)
}

type Cache struct {
	c *redis.Client
}

func NewCache(c *redis.Client) *Cache {
	return &Cache{
		c: c,
	}
}
func (c *Cache) Handle(ctx context.Context, key string, fn func() (*types.User, error)) (*types.User, error) {
	val, err := c.c.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		v, err := fn()
		if err != nil {
			return nil, errors.New("Handling went wrong: " + err.Error())
		}
		marshal, err := json.Marshal(v)
		if err != nil {
			return nil, errors.New("Handling went wrong: " + err.Error())
		}
		c.c.Set(ctx, key, marshal, time.Hour)
		return v, nil
	} else if err != nil {
		return nil, errors.New("Handling went wrong: " + err.Error())
	}
	var user types.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return nil, errors.New("Handling went wrong: " + err.Error())
	}
	return &user, nil
}
