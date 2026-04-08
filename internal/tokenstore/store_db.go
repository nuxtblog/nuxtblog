package tokenstore

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

type dbStore struct{}

func newDBStore() *dbStore {
	s := &dbStore{}
	s.createTable()
	return s
}

// createTable ensures the refresh_tokens table exists.
// Called once at startup so the DB driver is self-contained.
func (s *dbStore) createTable() {
	db := g.DB()
	ctx := context.Background()

	var tableSQL string
	if db.GetConfig().Type == "pgsql" {
		tableSQL = `CREATE TABLE IF NOT EXISTS refresh_tokens (
			id         BIGSERIAL   PRIMARY KEY,
			user_id    BIGINT      NOT NULL,
			token_hash TEXT        NOT NULL UNIQUE,
			expires_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`
	} else {
		tableSQL = `CREATE TABLE IF NOT EXISTS refresh_tokens (
			id         INTEGER  PRIMARY KEY AUTOINCREMENT,
			user_id    INTEGER  NOT NULL,
			token_hash TEXT     NOT NULL UNIQUE,
			expires_at DATETIME NOT NULL,
			created_at DATETIME NOT NULL DEFAULT (strftime('%Y-%m-%d %H:%M:%f', 'now'))
		)`
	}

	sqls := []string{
		tableSQL,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens (user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash ON refresh_tokens (token_hash)`,
	}
	for _, sql := range sqls {
		if _, err := db.Exec(ctx, sql); err != nil {
			g.Log().Warningf(ctx, "[tokenstore/db] init table: %v", err)
		}
	}
}

func (s *dbStore) Save(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	_, err := g.DB().Ctx(ctx).Model("refresh_tokens").Data(g.Map{
		"user_id":    userID,
		"token_hash": tokenHash,
		"expires_at": expiresAt.Format("2006-01-02 15:04:05"),
	}).Insert()
	return err
}

func (s *dbStore) Exists(ctx context.Context, tokenHash string) (bool, error) {
	cnt, err := g.DB().Ctx(ctx).Model("refresh_tokens").
		Where("token_hash", tokenHash).
		WhereGT("expires_at", time.Now().Format("2006-01-02 15:04:05")).
		Count()
	return cnt > 0, err
}

func (s *dbStore) Delete(ctx context.Context, tokenHash string) error {
	_, err := g.DB().Ctx(ctx).Model("refresh_tokens").
		Where("token_hash", tokenHash).Delete()
	return err
}

func (s *dbStore) DeleteByUser(ctx context.Context, userID int64) error {
	_, err := g.DB().Ctx(ctx).Model("refresh_tokens").
		Where("user_id", userID).Delete()
	return err
}
