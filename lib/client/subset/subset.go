package subset

import (
	"context"
	"github.com/rueian/rueidis"
	"time"
)

const redisKeyPrefix = "rmq:traffic:"

//go:bean
type Subset struct {
	client rueidis.Client
}

func New(client rueidis.Client) *Subset {
	return &Subset{
		client: client,
	}
}

func (s *Subset) GetMembers(ctx context.Context, key string) ([]string, bool, error) {
	key = redisKeyPrefix + key
	result := s.client.DoCache(ctx, s.client.B().Smembers().Key(key).Cache(), 1*time.Hour)
	if err := result.Error(); err != nil {
		return nil, false, err
	}
	ss, err := result.AsStrSlice()
	if err != nil {
		return nil, false, err
	}
	return ss, result.IsCacheHit(), nil
}

func (s *Subset) Add(ctx context.Context, key string, members ...string) error {
	for i, m := range members {
		members[i] = redisKeyPrefix + m
	}
	return s.client.Do(ctx, s.client.B().Sadd().Key(key).Member(members...).Build()).Error()
}

func (s *Subset) Remove(ctx context.Context, key string, members ...string) error {
	for i, m := range members {
		members[i] = redisKeyPrefix + m
	}
	return s.client.Do(ctx, s.client.B().Srem().Key(key).Member(members...).Build()).Error()
}
