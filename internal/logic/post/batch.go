package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

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
