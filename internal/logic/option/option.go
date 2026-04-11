package option

import (
	"context"
	"strings"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/option/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	permlogic "github.com/nuxtblog/nuxtblog/internal/logic/permission"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// secretKeyPrefixes defines option key prefixes that require admin access to read.
var secretKeyPrefixes = []string{"oauth_", "ai_configs", "ai_active_id", "payment_"}

func isSecretKey(key string) bool {
	for _, prefix := range secretKeyPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

type sOption struct{}

func New() service.IOption { return &sOption{} }

func init() {
	service.RegisterOption(New())
}

func (s *sOption) Get(ctx context.Context, key string) (*v1.OptionItem, error) {
	if isSecretKey(key) {
		role := middleware.GetCurrentUserRole(ctx)
		if role < middleware.RoleAdmin {
			return nil, gerror.NewCode(gcode.CodeNotAuthorized, "permission denied")
		}
	}
	type OptionRow struct {
		Key      string `orm:"key"`
		Value    string `orm:"value"`
		Autoload int    `orm:"autoload"`
	}
	var row OptionRow
	err := dao.Options.Ctx(ctx).Where("key", key).Scan(&row)
	if err != nil {
		return nil, err
	}
	if row.Key == "" {
		return nil, gerror.NewCode(gcode.CodeNotFound, g.I18n().T(ctx, "error.option_not_found"))
	}
	return &v1.OptionItem{
		Key:      row.Key,
		Value:    row.Value,
		Autoload: row.Autoload,
	}, nil
}

func (s *sOption) GetAutoload(ctx context.Context) (map[string]string, error) {
	type OptionRow struct {
		Key   string `orm:"key"`
		Value string `orm:"value"`
	}
	var rows []OptionRow
	err := dao.Options.Ctx(ctx).Where("autoload", 1).Fields("key, value").Scan(&rows)
	if err != nil {
		return nil, err
	}
	isAdmin := middleware.GetCurrentUserRole(ctx) >= middleware.RoleAdmin
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		if !isAdmin && isSecretKey(row.Key) {
			continue
		}
		result[row.Key] = row.Value
	}
	return result, nil
}

func (s *sOption) Set(ctx context.Context, req *v1.OptionSetReq) error {
	role := middleware.GetCurrentUserRole(ctx)
	if !service.Permission().Can(ctx, role, "manage_options") {
		return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: manage_options")
	}
	if req.Key == "role_capabilities" || req.Key == "custom_role_defs" {
		if !service.Permission().Can(ctx, role, "manage_roles") {
			return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: manage_roles")
		}
	}

	autoload := 1
	if req.Autoload != nil {
		autoload = *req.Autoload
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := dao.Options.Ctx(ctx).
		Where("key", req.Key).
		Update(g.Map{
			"value":      req.Value,
			"autoload":   autoload,
			"updated_at": now,
		})
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		_, err = dao.Options.Ctx(ctx).Insert(g.Map{
			"key":        req.Key,
			"value":      req.Value,
			"autoload":   autoload,
			"updated_at": now,
		})
	}
	if err == nil && req.Key == "custom_role_defs" {
		middleware.InvalidateCustomRoleCache(ctx)
	}
	if err == nil && req.Key == "role_capabilities" {
		permlogic.InvalidateCapCache(ctx)
	}
	if err == nil && req.Key == "site_language" {
		middleware.InvalidateSiteLangCache(ctx)
	}
	return err
}

func (s *sOption) Delete(ctx context.Context, key string) error {
	_, err := dao.Options.Ctx(ctx).Where("key", key).Delete()
	return err
}

func (s *sOption) AdminStats(ctx context.Context) (*v1.AdminStatsRes, error) {
	res := &v1.AdminStatsRes{}

	totalPosts, _ := dao.Posts.Ctx(ctx).WhereNull("deleted_at").Count()
	publishedPosts, _ := dao.Posts.Ctx(ctx).WhereNull("deleted_at").Where("status", 2).Count()
	draftPosts, _ := dao.Posts.Ctx(ctx).WhereNull("deleted_at").Where("status", 1).Count()
	res.Posts = v1.AdminStatsPostStats{
		Total:     totalPosts,
		Published: publishedPosts,
		Draft:     draftPosts,
	}

	totalComments, _ := dao.Comments.Ctx(ctx).WhereNull("deleted_at").Count()
	pendingComments, _ := dao.Comments.Ctx(ctx).WhereNull("deleted_at").Where("status", 1).Count()
	res.Comments = v1.AdminStatsCommentStats{
		Total:   totalComments,
		Pending: pendingComments,
	}

	totalUsers, _ := dao.Users.Ctx(ctx).WhereNull("deleted_at").Count()
	activeUsers, _ := dao.Users.Ctx(ctx).WhereNull("deleted_at").Where("status", 1).Count()
	res.Users = v1.AdminStatsUserStats{
		Total:  totalUsers,
		Active: activeUsers,
	}

	type ViewSum struct {
		Total int64 `orm:"total"`
	}
	var vs ViewSum
	_ = dao.PostStats.Ctx(ctx).Fields("SUM(view_count) as total").Scan(&vs)
	res.Views = vs.Total

	return res, nil
}
