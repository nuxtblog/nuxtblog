package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/model/entity"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"
	"github.com/nuxtblog/nuxtblog/internal/util/password"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sPost struct{}

func New() *sPost { return &sPost{} }

func init() {
	service.RegisterPost(New())
}

// ----------------------------------------------------------------
//  Create
// ----------------------------------------------------------------

func (s *sPost) Create(ctx context.Context, req *v1.PostCreateReq) (id int64, err error) {
	role := middleware.GetCurrentUserRole(ctx)
	// Downgrade to draft if the user lacks publish_posts capability.
	if int(req.Status) == 2 && !service.Permission().Can(ctx, role, "publish_posts") {
		req.Status = v1.PostStatus(1)
	}

	// Run plugin filters — allows plugins to modify title/slug/content before save
	filtered, filterErr := eng.Filter(ctx, eng.FilterPostCreate, map[string]any{
		"title":   req.Title,
		"slug":    req.Slug,
		"content": req.Content,
		"excerpt": req.Excerpt,
		"status":  int(req.Status),
	})
	if filterErr != nil {
		return 0, filterErr
	}
	if v, ok := filtered["title"].(string); ok && v != "" {
		req.Title = v
	}
	if v, ok := filtered["slug"].(string); ok && v != "" {
		req.Slug = v
	}
	if v, ok := filtered["content"].(string); ok {
		req.Content = v
	}
	if v, ok := filtered["excerpt"].(string); ok {
		req.Excerpt = v
	}

	uid, _ := middleware.GetCurrentUserID(ctx)

	post := &entity.Posts{
		PostType:      int(req.PostType),
		Status:        int(req.Status),
		Title:         req.Title,
		Slug:          req.Slug,
		Content:       req.Content,
		Excerpt:       req.Excerpt,
		CommentStatus: int(req.CommentStatus),
		Locale:        req.Locale,
		PublishedAt:   req.PublishedAt,
		AuthorId:      int(uid),
	}
	if req.FeaturedImgId != nil {
		post.FeaturedImgId = int(*req.FeaturedImgId)
	}
	// Admin can set a different author
	if req.AuthorId != nil && role >= middleware.RoleAdmin {
		post.AuthorId = int(*req.AuthorId)
	}

	id, err = dao.Posts.Create(ctx, post)
	if err != nil {
		return
	}
	_ = event.Emit(ctx, event.PostCreated, payload.PostCreated{
		PostID:   id,
		AuthorID: uid,
		Title:    post.Title,
		Slug:     post.Slug,
		Excerpt:  post.Excerpt,
		PostType: post.PostType,
		Status:   post.Status,
	})

	// 关联分类（去重后插入）
	if len(req.TermTaxonomyIds) > 0 {
		seen := make(map[int64]struct{})
		rows := make([]g.Map, 0)
		for _, taxId := range req.TermTaxonomyIds {
			if _, dup := seen[taxId]; dup {
				continue
			}
			seen[taxId] = struct{}{}
			rows = append(rows, g.Map{
				"object_id":   id,
				"object_type": "post",
				"taxonomy_id": taxId,
				"sort_order":  0,
			})
		}
		if len(rows) > 0 {
			if _, err = dao.ObjectTaxonomies.Ctx(ctx).Data(rows).Insert(); err != nil {
				return
			}
		}
	}

	// 写入 metas
	if len(req.Metas) > 0 {
		err = s.UpsertMetas(ctx, id, req.Metas)
	}
	return
}

// ----------------------------------------------------------------
//  Update
// ----------------------------------------------------------------

func (s *sPost) Update(ctx context.Context, req *v1.PostUpdateReq) error {
	role := middleware.GetCurrentUserRole(ctx)
	uid, _ := middleware.GetCurrentUserID(ctx)

	// For admin-level users (bypassed OwnershipCheck), verify edit_others_posts.
	if role >= middleware.RoleAdmin {
		type ownerRow struct{ AuthorId int64 `orm:"author_id"` }
		var row ownerRow
		if e := dao.Posts.Ctx(ctx).Where("id", req.Id).Scan(&row); e == nil && row.AuthorId != uid {
			if !service.Permission().Can(ctx, role, "edit_others_posts") {
				return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: edit_others_posts")
			}
		}
	}

	// Downgrade to draft if the user lacks publish_posts capability.
	if req.Status != nil && int(*req.Status) == 2 && !service.Permission().Can(ctx, role, "publish_posts") {
		s := v1.PostStatus(1)
		req.Status = &s
	}

	// Run plugin filters for fields that are being updated
	if req.Title != nil || req.Slug != nil || req.Content != nil || req.Excerpt != nil {
		filterIn := map[string]any{}
		if req.Title != nil {
			filterIn["title"] = *req.Title
		}
		if req.Slug != nil {
			filterIn["slug"] = *req.Slug
		}
		if req.Content != nil {
			filterIn["content"] = *req.Content
		}
		if req.Excerpt != nil {
			filterIn["excerpt"] = *req.Excerpt
		}
		if filtered, filterErr := eng.Filter(ctx, eng.FilterPostUpdate, filterIn); filterErr != nil {
			return filterErr
		} else {
			if v, ok := filtered["title"].(string); ok && req.Title != nil {
				req.Title = &v
			}
			if v, ok := filtered["slug"].(string); ok {
				req.Slug = &v
			}
			if v, ok := filtered["content"].(string); ok && req.Content != nil {
				req.Content = &v
			}
			if v, ok := filtered["excerpt"].(string); ok && req.Excerpt != nil {
				req.Excerpt = &v
			}
		}
	}

	data := g.Map{}
	if req.Title != nil {
		data["title"] = *req.Title
	}
	if req.Slug != nil {
		data["slug"] = *req.Slug
	}
	if req.Content != nil {
		data["content"] = *req.Content
	}
	if req.Excerpt != nil {
		data["excerpt"] = *req.Excerpt
	}
	if req.Status != nil {
		data["status"] = int(*req.Status)
	}
	if req.CommentStatus != nil {
		data["comment_status"] = int(*req.CommentStatus)
	}
	if req.FeaturedImgId != nil {
		data["featured_img_id"] = *req.FeaturedImgId
	}
	if req.Locale != nil {
		data["locale"] = *req.Locale
	}
	if req.PublishedAt != nil {
		data["published_at"] = req.PublishedAt
	}
	if req.AuthorId != nil && role >= middleware.RoleAdmin {
		data["author_id"] = *req.AuthorId
	}

	if len(data) > 0 {
		// Snapshot current content before overwriting
		if req.Content != nil || req.Title != nil {
			type snapRow struct {
				Title    string `orm:"title"`
				Content  string `orm:"content"`
				AuthorId int64  `orm:"author_id"`
			}
			var snap snapRow
			if e := dao.Posts.Ctx(ctx).Where("id", req.Id).Scan(&snap); e == nil && snap.Title != "" {
				go saveRevisionAndPrune(ctx, req.Id, snap.Title, snap.Content, snap.AuthorId)
			}
		}
		if err := dao.Posts.UpdateById(ctx, req.Id, data); err != nil {
			return err
		}
		// Emit post.updated and (if newly published) post.published.
		if post, _ := dao.Posts.GetById(ctx, req.Id); post != nil {
			_ = event.Emit(ctx, event.PostUpdated, payload.PostUpdated{
				PostID:   req.Id,
				AuthorID: int64(post.AuthorId),
				Title:    post.Title,
				Slug:     post.Slug,
				Excerpt:  post.Excerpt,
				PostType: post.PostType,
				Status:   post.Status,
			})
			if req.Status != nil && int(*req.Status) == 2 {
				if _, ferr := eng.Filter(ctx, eng.FilterPostPublish, map[string]any{
					"id":    req.Id,
					"title": post.Title,
					"slug":  post.Slug,
				}); ferr != nil {
					return ferr
				}
				_ = event.Emit(ctx, event.PostPublished, payload.PostPublished{
					PostID:   req.Id,
					AuthorID: int64(post.AuthorId),
					Title:    post.Title,
					Slug:     post.Slug,
					Excerpt:  post.Excerpt,
					PostType: post.PostType,
				})
			}
		}
	}

	// 更新分类关联（整体替换）
	if req.TermTaxonomyIds != nil {
		if _, err := dao.ObjectTaxonomies.Ctx(ctx).
			Where("object_id", req.Id).
			Where("object_type", "post").
			Delete(); err != nil {
			return err
		}
		if len(*req.TermTaxonomyIds) > 0 {
			seen := make(map[int64]struct{})
			rows := make([]g.Map, 0)
			for _, taxId := range *req.TermTaxonomyIds {
				if _, dup := seen[taxId]; dup {
					continue
				}
				seen[taxId] = struct{}{}
				rows = append(rows, g.Map{
					"object_id":   req.Id,
					"object_type": "post",
					"taxonomy_id": taxId,
					"sort_order":  0,
				})
			}
			if len(rows) > 0 {
				if _, err := dao.ObjectTaxonomies.Ctx(ctx).Data(rows).Insert(); err != nil {
					return err
				}
			}
		}
	}

	// 更新 metas
	if req.Metas != nil {
		return s.UpsertMetas(ctx, req.Id, req.Metas)
	}
	return nil
}

// ----------------------------------------------------------------
//  Delete / Trash / Restore
// ----------------------------------------------------------------

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

// Batch performs a bulk action on a list of post IDs.
//
// Ownership & capability rules:
//   - Editors (role < Admin) are always scoped to their own posts (author_id = uid).
//   - Admins are also scoped to own posts unless they hold the relevant cross-user cap:
//     "publish" → publish_posts,  "trash"/"delete" → delete_others_posts.
//   - publish_posts is checked for both editors and admins on the "publish" action.
func (s *sPost) Batch(ctx context.Context, req *v1.PostBatchReq) (int, error) {
	if len(req.Ids) == 0 {
		return 0, nil
	}

	role := middleware.GetCurrentUserRole(ctx)
	uid, _ := middleware.GetCurrentUserID(ctx)

	// "publish" requires the publish_posts capability regardless of role.
	if req.Action == "publish" && !service.Permission().Can(ctx, role, "publish_posts") {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: publish_posts")
	}

	// Determine whether this user may affect other users' posts.
	// Editors never can; admins need the delete_others_posts cap for destructive actions.
	canAffectOthers := role >= middleware.RoleAdmin &&
		service.Permission().Can(ctx, role, "delete_others_posts")

	// Build base query, scoping to own posts when needed.
	base := dao.Posts.Ctx(ctx).WhereIn("id", req.Ids)
	if !canAffectOthers || role < middleware.RoleAdmin {
		base = base.Where("author_id", uid)
	}

	var (
		affected int64
		err      error
	)
	switch req.Action {
	case "publish":
		res, e := base.WhereNull("deleted_at").Data(g.Map{"status": 2}).Update()
		if e != nil {
			return 0, e
		}
		affected, _ = res.RowsAffected()
	case "draft":
		res, e := base.WhereNull("deleted_at").Data(g.Map{"status": 1}).Update()
		if e != nil {
			return 0, e
		}
		affected, _ = res.RowsAffected()
	case "trash":
		res, e := base.WhereNot("status", 5).Data(g.Map{"status": 5}).Update()
		if e != nil {
			return 0, e
		}
		affected, _ = res.RowsAffected()
	case "restore":
		res, e := base.Where("status", 5).Data(g.Map{"status": 1}).Update()
		if e != nil {
			return 0, e
		}
		affected, _ = res.RowsAffected()
	case "delete":
		res, e := base.Delete()
		if e != nil {
			return 0, e
		}
		affected, _ = res.RowsAffected()
	}
	return int(affected), err
}

// ----------------------------------------------------------------
//  BatchUpdate
// ----------------------------------------------------------------

// BatchUpdate applies field-level updates to multiple posts in one call.
// It resolves ownership scope first, then updates scalar fields and taxonomy
// associations separately. Returns the number of posts actually affected.
func (s *sPost) BatchUpdate(ctx context.Context, req *v1.PostBatchUpdateReq) (int, error) {
	if len(req.Ids) == 0 {
		return 0, nil
	}

	role := middleware.GetCurrentUserRole(ctx)
	uid, _ := middleware.GetCurrentUserID(ctx)

	// Require publish_posts when setting status to published.
	if req.Status != nil && *req.Status == v1.PostStatusPublished &&
		!service.Permission().Can(ctx, role, "publish_posts") {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: publish_posts")
	}
	// Only admins may reassign author.
	if req.AuthorId != nil && role < middleware.RoleAdmin {
		return 0, gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: change author requires admin")
	}

	canAffectOthers := role >= middleware.RoleAdmin &&
		service.Permission().Can(ctx, role, "edit_others_posts")

	// Resolve the actual post IDs this user may modify (respects soft-delete and ownership).
	type idRow struct {
		Id int64 `orm:"id"`
	}
	var idRows []idRow
	q := dao.Posts.Ctx(ctx).WhereIn("id", req.Ids).WhereNull("deleted_at").Fields("id")
	if !canAffectOthers {
		q = q.Where("author_id", uid)
	}
	if err := q.Scan(&idRows); err != nil {
		return 0, err
	}
	if len(idRows) == 0 {
		return 0, nil
	}
	actualIds := make([]int64, len(idRows))
	for i, r := range idRows {
		actualIds[i] = r.Id
	}

	// Update scalar fields on posts table.
	data := g.Map{}
	if req.FeaturedImgId != nil {
		if *req.FeaturedImgId == 0 {
			data["featured_img_id"] = nil // clear
		} else {
			data["featured_img_id"] = *req.FeaturedImgId
		}
	}
	if req.Status != nil {
		data["status"] = int(*req.Status)
	}
	if req.AuthorId != nil {
		data["author_id"] = *req.AuthorId
	}
	if len(data) > 0 {
		if _, err := dao.Posts.Ctx(ctx).WhereIn("id", actualIds).Data(data).Update(); err != nil {
			return 0, err
		}
	}

	// Replace taxonomy associations (nil = no change; non-nil = full replace).
	if req.TermTaxonomyIds != nil {
		if _, err := dao.ObjectTaxonomies.Ctx(ctx).
			WhereIn("object_id", actualIds).
			Where("object_type", "post").
			Delete(); err != nil {
			return 0, err
		}
		if len(*req.TermTaxonomyIds) > 0 {
			// Deduplicate the requested taxonomy IDs.
			seen := make(map[int64]struct{})
			var taxIds []int64
			for _, tid := range *req.TermTaxonomyIds {
				if _, dup := seen[tid]; dup {
					continue
				}
				seen[tid] = struct{}{}
				taxIds = append(taxIds, tid)
			}
			insertRows := make([]g.Map, 0, len(actualIds)*len(taxIds))
			for _, postId := range actualIds {
				for _, taxId := range taxIds {
					insertRows = append(insertRows, g.Map{
						"object_id":   postId,
						"object_type": "post",
						"taxonomy_id": taxId,
						"sort_order":  0,
					})
				}
			}
			if _, err := dao.ObjectTaxonomies.Ctx(ctx).Data(insertRows).Insert(); err != nil {
				return 0, err
			}
		}
	}

	return len(actualIds), nil
}

// ----------------------------------------------------------------
//  View
// ----------------------------------------------------------------

func (s *sPost) IncrementView(ctx context.Context, id int64) error {
	_, err := dao.PostStats.DB().Exec(
		ctx,
		`INSERT INTO post_stats (post_id, view_count)
		 VALUES (?, 1)
		 ON CONFLICT(post_id) DO UPDATE SET view_count = view_count + 1`,
		id,
	)
	if err != nil {
		g.Log().Warningf(ctx, "[post] increment view error: %v", err)
	}
	var userID int64
	if uid, ok := middleware.GetCurrentUserID(ctx); ok && uid > 0 {
		userID = uid
		_, _ = dao.UserActions.Ctx(ctx).Data(g.Map{
			"id":          idgen.New(),
			"user_id":     uid,
			"action":      "view",
			"object_type": "post",
			"object_id":   id,
			"extra":       "{}",
		}).Insert()
	}
	_ = event.Emit(ctx, event.PostViewed, payload.PostViewed{PostID: id, UserID: userID})
	return nil
}

// ----------------------------------------------------------------
//  VerifyPassword
// ----------------------------------------------------------------

func (s *sPost) VerifyPassword(ctx context.Context, id int64, plain string) (bool, error) {
	type hashRow struct {
		PasswordHash string `orm:"password_hash"`
	}
	var row hashRow
	if err := dao.Posts.Ctx(ctx).Where("id", id).Fields("password_hash").Scan(&row); err != nil {
		return false, err
	}
	if row.PasswordHash == "" {
		return false, nil
	}
	return password.Verify(plain, row.PasswordHash), nil
}

// saveRevisionAndPrune saves a snapshot of the current post content before an update,
// then deletes all but the most recent maxRevisions entries for the post.
func saveRevisionAndPrune(ctx context.Context, postId int64, title, content string, authorId int64) {
	const maxRevisions = 20

	// Insert snapshot
	_, err := dao.PostRevisions.Ctx(ctx).Data(g.Map{
		"id":        idgen.New(),
		"post_id":   postId,
		"author_id": authorId,
		"title":     title,
		"content":   content,
		"rev_note":  "auto",
	}).Insert()
	if err != nil {
		g.Log().Warningf(ctx, "[post] save revision error: %v", err)
		return
	}

	// Prune: keep only the most recent maxRevisions rows
	type idRow struct{ Id int64 }
	var rows []idRow
	if err = dao.PostRevisions.Ctx(ctx).
		Where("post_id", postId).
		OrderDesc("id").
		Offset(maxRevisions).
		Fields("id").
		Scan(&rows); err != nil || len(rows) == 0 {
		return
	}
	ids := make([]int64, len(rows))
	for i, r := range rows {
		ids[i] = r.Id
	}
	_, _ = dao.PostRevisions.Ctx(ctx).WhereIn("id", ids).Delete()
}
