package post

import (
	"context"

	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PublishScheduled scans draft posts whose published_at <= now() and transitions
// them to published status, firing post.published for each one.
// Called every minute by the cron job registered in cmd.go.
func (s *sPost) PublishScheduled(ctx context.Context) {
	type row struct {
		Id       int64  `orm:"id"`
		AuthorId int64  `orm:"author_id"`
		Title    string `orm:"title"`
		Slug     string `orm:"slug"`
		PostType int    `orm:"post_type"`
	}
	var rows []row
	if err := dao.Posts.Ctx(ctx).
		Where("status", 1). // draft
		WhereNotNull("published_at").
		WhereLTE("published_at", gtime.Now()).
		WhereNull("deleted_at").
		Fields("id, author_id, title, slug, post_type").
		Scan(&rows); err != nil || len(rows) == 0 {
		return
	}

	ids := make([]int64, len(rows))
	for i, r := range rows {
		ids[i] = r.Id
	}
	if _, err := dao.Posts.Ctx(ctx).WhereIn("id", ids).Data(g.Map{"status": 2}).Update(); err != nil {
		g.Log().Warningf(ctx, "[schedule] publish posts error: %v", err)
		return
	}
	for _, r := range rows {
		_ = event.Emit(ctx, event.PostPublished, payload.PostPublished{
			PostID:   r.Id,
			AuthorID: r.AuthorId,
			Title:    r.Title,
			Slug:     r.Slug,
			PostType: r.PostType,
		})
	}
	g.Log().Infof(ctx, "[schedule] auto-published %d post(s)", len(rows))
}
