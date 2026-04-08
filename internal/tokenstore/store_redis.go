package tokenstore

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2" // register GoFrame Redis adapter
)

type redisStore struct{}

func newRedisStore() *redisStore { return &redisStore{} }

// Key layout:
//
//	rt:{hash}    — string key; value = userID; TTL = token lifetime
//	rtu:{userID} — Redis SET; members = token hashes; TTL = max token lifetime
func (s *redisStore) Save(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	rds := g.Redis()
	ttl := int(time.Until(expiresAt).Seconds())
	if ttl <= 0 {
		return nil
	}
	tokenKey := fmt.Sprintf("rt:%s", tokenHash)
	userKey := fmt.Sprintf("rtu:%d", userID)
	if _, err := rds.Do(ctx, "SET", tokenKey, userID, "EX", ttl); err != nil {
		return err
	}
	_, _ = rds.Do(ctx, "SADD", userKey, tokenHash)
	_, _ = rds.Do(ctx, "EXPIRE", userKey, ttl)
	return nil
}

func (s *redisStore) Exists(ctx context.Context, tokenHash string) (bool, error) {
	rds := g.Redis()
	val, err := rds.Do(ctx, "EXISTS", fmt.Sprintf("rt:%s", tokenHash))
	if err != nil {
		return false, err
	}
	return val.Int() == 1, nil
}

func (s *redisStore) Delete(ctx context.Context, tokenHash string) error {
	rds := g.Redis()
	_, err := rds.Do(ctx, "DEL", fmt.Sprintf("rt:%s", tokenHash))
	return err
}

func (s *redisStore) DeleteByUser(ctx context.Context, userID int64) error {
	rds := g.Redis()
	userKey := fmt.Sprintf("rtu:%d", userID)
	val, err := rds.Do(ctx, "SMEMBERS", userKey)
	if err != nil {
		return err
	}
	hashes := val.Strings()
	if len(hashes) == 0 {
		return nil
	}
	// Build DEL args: userKey + all rt:{hash} keys
	keys := make([]interface{}, 0, len(hashes)+1)
	keys = append(keys, userKey)
	for _, h := range hashes {
		keys = append(keys, fmt.Sprintf("rt:%s", h))
	}
	_, err = rds.Do(ctx, "DEL", keys...)
	return err
}
