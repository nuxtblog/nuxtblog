package post

import (
	"context"

	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Delete permanently removes a post (force delete).
func (s *sPost) Delete(ctx context.Context, id int64) error {
	if err := checkDeletePost(ctx, id); err != nil {
		return err
	}
	if _, ferr := eng.Filter(ctx, eng.FilterPostDelete, map[string]any{"id": id}); ferr != nil {
		return ferr
	}
	type titleRow struct {
		AuthorId int64  `orm:"author_id"`
		Title    string `orm:"title"`
		Slug     string `orm:"slug"`
		PostType int    `orm:"post_type"`
	}
	var row titleRow
	_ = dao.Posts.Ctx(ctx).Where("id", id).Scan(&row)
	if err := dao.Posts.DeleteById(ctx, id); err != nil {
		return err
	}
	_ = event.Emit(ctx, event.PostDeleted, payload.PostDeleted{
		PostID: id, AuthorID: row.AuthorId, Title: row.Title,
		Slug: row.Slug, PostType: row.PostType,
	})
	return nil
}

// Trash soft-deletes a post by setting deleted_at.
func (s *sPost) Trash(ctx context.Context, id int64) error {
	if err := checkDeletePost(ctx, id); err != nil {
		return err
	}
	if _, ferr := eng.Filter(ctx, eng.FilterPostDelete, map[string]any{"id": id}); ferr != nil {
		return ferr
	}
	type titleRow struct {
		AuthorId int64  `orm:"author_id"`
		Title    string `orm:"title"`
		Slug     string `orm:"slug"`
		PostType int    `orm:"post_type"`
	}
	var row titleRow
	_ = dao.Posts.Ctx(ctx).Where("id", id).Scan(&row)
	if err := dao.Posts.TrashById(ctx, id); err != nil {
		return err
	}
	_ = event.Emit(ctx, event.PostDeleted, payload.PostDeleted{
		PostID: id, AuthorID: row.AuthorId, Title: row.Title,
		Slug: row.Slug, PostType: row.PostType,
	})
	return nil
}

// checkDeletePost verifies delete_others_posts capability for admin-level users.
func checkDeletePost(ctx context.Context, id int64) error {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return nil // editors are already guarded by OwnershipCheck middleware
	}
	uid, _ := middleware.GetCurrentUserID(ctx)
	type ownerRow struct{ AuthorId int64 `orm:"author_id"` }
	var row ownerRow
	if e := dao.Posts.Ctx(ctx).Where("id", id).Scan(&row); e == nil && row.AuthorId != uid {
		if !service.Permission().Can(ctx, role, "delete_others_posts") {
			return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: delete_others_posts")
		}
	}
	return nil
}

// Restore clears deleted_at, recovering a trashed post.
func (s *sPost) Restore(ctx context.Context, id int64) error {
	if _, ferr := eng.Filter(ctx, eng.FilterPostRestore, map[string]any{"id": id}); ferr != nil {
		return ferr
	}
	return dao.Posts.RestoreById(ctx, id)
}
