package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/token/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type sToken struct{}

func New() service.IToken { return &sToken{} }

func init() {
	service.RegisterToken(New())
}

func currentUserID(ctx context.Context) (int64, error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return 0, errors.New(g.I18n().T(ctx, "error.unauthorized"))
	}
	id := r.GetCtxVar("user_id").Int64()
	if id == 0 {
		return 0, errors.New(g.I18n().T(ctx, "error.unauthorized"))
	}
	return id, nil
}

type tokenRow struct {
	Id         int64  `orm:"id"`
	Name       string `orm:"name"`
	Prefix     string `orm:"prefix"`
	ExpiresAt  string `orm:"expires_at"`
	LastUsedAt string `orm:"last_used_at"`
	CreatedAt  string `orm:"created_at"`
}

func nullableStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// List returns all personal API tokens for the current user.
func (s *sToken) List(ctx context.Context, _ *v1.UserTokenListReq) (*v1.UserTokenListRes, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	var rows []tokenRow
	_ = dao.UserTokens.Ctx(ctx).
		Fields("id,name,prefix,expires_at,last_used_at,created_at").
		Where("user_id", userID).
		OrderDesc("id").
		Scan(&rows)

	list := make([]v1.UserTokenItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, v1.UserTokenItem{
			Id:         row.Id,
			Name:       row.Name,
			Prefix:     row.Prefix,
			ExpiresAt:  nullableStr(row.ExpiresAt),
			LastUsedAt: nullableStr(row.LastUsedAt),
			CreatedAt:  row.CreatedAt,
		})
	}
	return &v1.UserTokenListRes{List: list}, nil
}

// Create generates a new personal API token for the current user.
func (s *sToken) Create(ctx context.Context, req *v1.UserTokenCreateReq) (*v1.UserTokenCreateRes, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	raw := make([]byte, 20)
	if _, err := rand.Read(raw); err != nil {
		return nil, errors.New(g.I18n().T(ctx, "token.generate_failed"))
	}
	token := "yblog_" + hex.EncodeToString(raw)
	h := sha256.Sum256([]byte(token))
	hashHex := hex.EncodeToString(h[:])
	prefix := token[:14] // "yblog_" + first 8 hex chars

	now := time.Now()
	data := g.Map{
		"id":         idgen.New(),
		"user_id":    userID,
		"name":       req.Name,
		"prefix":     prefix,
		"token_hash": hashHex,
		"created_at": now,
	}
	var expiresAt *string
	if req.ExpiresInDays > 0 {
		t := now.Add(time.Duration(req.ExpiresInDays) * 24 * time.Hour)
		data["expires_at"] = t
		s := t.Format(time.RFC3339)
		expiresAt = &s
	}

	result, err := dao.UserTokens.Ctx(ctx).Data(data).Insert()
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	return &v1.UserTokenCreateRes{
		Id:        id,
		Name:      req.Name,
		Prefix:    prefix,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: now.Format(time.RFC3339),
	}, nil
}

// Delete removes a token owned by the current user.
func (s *sToken) Delete(ctx context.Context, req *v1.UserTokenDeleteReq) (*v1.UserTokenDeleteRes, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := dao.UserTokens.Ctx(ctx).
		Where("id", req.Id).
		Where("user_id", userID).
		Delete(); err != nil {
		return nil, err
	}
	return &v1.UserTokenDeleteRes{}, nil
}
