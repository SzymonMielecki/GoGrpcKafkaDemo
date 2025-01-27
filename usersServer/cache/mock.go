package cache

import (
	"context"
	"errors"
	"github.com/SzymonMielecki/GoGrpcKafkaDemo/types"
)

type MockCache struct {
	cache map[string]*types.User
}

func (m MockCache) Handle(ctx context.Context, key string, fn func() (*types.User, error)) (*types.User, error) {
	if v, ok := m.cache[key]; ok {
		return v, nil
	}
	v, err := fn()
	if err != nil {
		return nil, errors.New("Handling went wrong: " + err.Error())
	}
	m.cache[key] = v
	return v, nil
}

func NewMockCache() *MockCache {
	return &MockCache{make(map[string]*types.User)}
}

func NewMockCacheWithCache(cache map[string]*types.User) *MockCache {
	return &MockCache{cache}
}
