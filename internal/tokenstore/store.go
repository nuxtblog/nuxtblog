// Package tokenstore provides a pluggable store for refresh-token hashes.
// Configure the backend via auth.tokenStore in config.yaml:
//
//	"memory" (default) — in-process map; fast but lost on restart
//	"db"               — SQL table refresh_tokens; survives restarts
//	"redis"            — Redis keys with TTL; survives restarts + multi-instance
package tokenstore

import (
	"context"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Store is the pluggable interface for refresh-token persistence.
type Store interface {
	// Save persists a refresh-token hash for userID until expiresAt.
	Save(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error
	// Exists returns true when the hash is known and not yet expired.
	Exists(ctx context.Context, tokenHash string) (bool, error)
	// Delete revokes a single refresh token by hash.
	Delete(ctx context.Context, tokenHash string) error
	// DeleteByUser revokes all refresh tokens for a user (e.g. on password change).
	DeleteByUser(ctx context.Context, userID int64) error
}

var (
	once    sync.Once
	current Store
)

// Default returns the singleton store, initialised once from config.
func Default(ctx context.Context) Store {
	once.Do(func() {
		val, _ := g.Cfg().Get(ctx, "auth.tokenStore")
		driver := val.String()
		switch driver {
		case "db":
			current = newDBStore()
		case "redis":
			current = newRedisStore()
		default:
			driver = "memory"
			current = newMemoryStore()
		}
		g.Log().Infof(ctx, "[tokenstore] driver: %s", driver)
	})
	return current
}
