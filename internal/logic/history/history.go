package history

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/history/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sHistory struct{}

func New() service.IHistory { return &sHistory{} }
func init()                 { service.RegisterHistory(New()) }

func (s *sHistory) List(ctx context.Context, req *v1.HistoryListReq) (*v1.HistoryListRes, error) {
	userID, _ := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, g.I18n().T(ctx, "error.unauthorized"))
	}

	offset := (req.Page - 1) * req.Size

	type row struct {
		PostId    int64  `orm:"post_id"`
		PostTitle string `orm:"post_title"`
		PostSlug  string `orm:"post_slug"`
		PostCover string `orm:"post_cover"`
		ViewedAt  string `orm:"viewed_at"`
	}

	// Get latest unique post views (one entry per post, latest time)
	var rows []row
	err := dao.UserActions.DB().Ctx(ctx).Raw(`
		SELECT ua.object_id AS post_id, p.title AS post_title, p.slug AS post_slug,
		       COALESCE(m.cdn_url, '') AS post_cover, MAX(ua.created_at) AS viewed_at
		FROM user_actions ua
		LEFT JOIN posts p ON p.id = ua.object_id AND p.deleted_at IS NULL
		LEFT JOIN medias m ON m.id = p.featured_img_id
		WHERE ua.user_id = ? AND ua.action = 'view' AND ua.object_type = 'post'
		  AND p.id IS NOT NULL
		GROUP BY ua.object_id
		ORDER BY viewed_at DESC
		LIMIT ? OFFSET ?
	`, userID, req.Size, offset).Scan(&rows)
	if err != nil {
		return nil, err
	}

	// Count total unique posts
	var total int
	_ = dao.UserActions.DB().Ctx(ctx).Raw(`
		SELECT COUNT(DISTINCT object_id) FROM user_actions
		WHERE user_id = ? AND action = 'view' AND object_type = 'post'
	`, userID).Scan(&total)

	items := make([]v1.HistoryItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, v1.HistoryItem{
			PostId:    r.PostId,
			PostTitle: r.PostTitle,
			PostSlug:  r.PostSlug,
			PostCover: r.PostCover,
			ViewedAt:  r.ViewedAt,
		})
	}
	return &v1.HistoryListRes{Items: items, Total: total}, nil
}

func (s *sHistory) Clear(ctx context.Context) error {
	userID, _ := middleware.GetCurrentUserID(ctx)
	if userID == 0 {
		return gerror.NewCode(gcode.CodeNotAuthorized, g.I18n().T(ctx, "error.unauthorized"))
	}
	_, err := dao.UserActions.Ctx(ctx).
		Where("user_id", userID).
		Where("action", "view").
		Where("object_type", "post").
		Delete()
	return err
}
