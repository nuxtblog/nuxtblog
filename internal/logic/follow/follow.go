package follow

import (
	"context"
	"errors"

	v1 "github.com/nuxtblog/nuxtblog/api/follow/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

type sFollow struct{}

func New() service.IFollow { return &sFollow{} }

func init() {
	service.RegisterFollow(New())
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

func optionalUserID(ctx context.Context) int64 {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	return r.GetCtxVar("user_id").Int64()
}

// Follow — POST /users/{id}/follow
func (s *sFollow) Follow(ctx context.Context, req *v1.UserFollowReq) (*v1.UserFollowRes, error) {
	followerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if req.Id == followerID {
		return nil, errors.New(g.I18n().T(ctx, "follow.cannot_follow_self"))
	}

	cnt, _ := dao.Users.Ctx(ctx).Where("id", req.Id).WhereNull("deleted_at").Count()
	if cnt == 0 {
		return nil, errors.New(g.I18n().T(ctx, "error.user_not_found"))
	}

	// Idempotent
	existing, _ := dao.UserFollows.Ctx(ctx).
		Where("follower_id", followerID).Where("following_id", req.Id).Count()
	if existing > 0 {
		return &v1.UserFollowRes{Following: true}, nil
	}

	if _, err := dao.UserFollows.Ctx(ctx).Data(g.Map{
		"follower_id":  followerID,
		"following_id": req.Id,
		"created_at":   gtime.Now(),
	}).Insert(); err != nil {
		return nil, err
	}

	// Emit user.followed — notification delivery is handled by the listener.
	type userRow struct {
		DisplayName string `orm:"display_name"`
		Avatar      string `orm:"avatar"`
	}
	var actor userRow
	_ = dao.Users.Ctx(ctx).Where("id", followerID).Fields("display_name,avatar").Scan(&actor)
	_ = event.Emit(ctx, event.UserFollowed, payload.UserFollowed{
		FollowerID:     followerID,
		FollowerName:   actor.DisplayName,
		FollowerAvatar: actor.Avatar,
		FollowingID:    req.Id,
	})

	return &v1.UserFollowRes{Following: true}, nil
}

// Unfollow — DELETE /users/{id}/follow
func (s *sFollow) Unfollow(ctx context.Context, req *v1.UserUnfollowReq) (*v1.UserUnfollowRes, error) {
	followerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := dao.UserFollows.Ctx(ctx).
		Where("follower_id", followerID).Where("following_id", req.Id).Delete(); err != nil {
		return nil, err
	}
	return &v1.UserUnfollowRes{Following: false}, nil
}

// RemoveFollower — DELETE /users/{id}/follower
func (s *sFollow) RemoveFollower(ctx context.Context, req *v1.UserRemoveFollowerReq) (*v1.UserRemoveFollowerRes, error) {
	viewerID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := dao.UserFollows.Ctx(ctx).
		Where("follower_id", req.Id).Where("following_id", viewerID).Delete(); err != nil {
		return nil, err
	}
	return &v1.UserRemoveFollowerRes{Removed: true}, nil
}

// FollowStatus — GET /users/{id}/follow-status
func (s *sFollow) FollowStatus(ctx context.Context, req *v1.UserFollowStatusReq) (*v1.UserFollowStatusRes, error) {
	followerCount, _ := dao.UserFollows.Ctx(ctx).Where("following_id", req.Id).Count()
	followingCount, _ := dao.UserFollows.Ctx(ctx).Where("follower_id", req.Id).Count()

	following := false
	if selfID := optionalUserID(ctx); selfID > 0 {
		cnt, _ := dao.UserFollows.Ctx(ctx).
			Where("follower_id", selfID).Where("following_id", req.Id).Count()
		following = cnt > 0
	}
	return &v1.UserFollowStatusRes{
		Following:      following,
		FollowerCount:  followerCount,
		FollowingCount: followingCount,
	}, nil
}

type followUserRow struct {
	ID            int64  `orm:"id"`
	Username      string `orm:"username"`
	DisplayName   string `orm:"display_name"`
	AvatarId      int64  `orm:"avatar_id"`
	Bio           string `orm:"bio"`
	FollowedAt    string `orm:"followed_at"`
	ArticleCount  int    `orm:"article_count"`
	FollowerCount int    `orm:"follower_count"`
}

// Followers — GET /users/{id}/followers
func (s *sFollow) Followers(ctx context.Context, req *v1.UserFollowersReq) (*v1.UserFollowersRes, error) {
	page, size := normPage(req.Page, req.Size)

	total, _ := dao.UserFollows.Ctx(ctx).Where("following_id", req.Id).Count()
	rows := make([]followUserRow, 0)
	if total > 0 {
		offset := (page - 1) * size
		_ = dao.UserFollows.DB().Ctx(ctx).Raw(`
			SELECT u.id, u.username, u.display_name, u.avatar_id, u.bio,
			       uf.created_at AS followed_at,
			       (SELECT COUNT(*) FROM posts WHERE author_id = u.id AND status = 2 AND deleted_at IS NULL) AS article_count,
			       (SELECT COUNT(*) FROM user_follows WHERE following_id = u.id) AS follower_count
			FROM user_follows uf
			JOIN users u ON u.id = uf.follower_id
			WHERE uf.following_id = ? AND u.deleted_at IS NULL
			ORDER BY uf.created_at DESC
			LIMIT ? OFFSET ?`, req.Id, size, offset).Scan(&rows)
	}

	list := toFollowUserItems(ctx, rows)
	return &v1.UserFollowersRes{List: list, Total: total, Page: page, Size: size}, nil
}

// Following — GET /users/{id}/following
func (s *sFollow) Following(ctx context.Context, req *v1.UserFollowingReq) (*v1.UserFollowingRes, error) {
	page, size := normPage(req.Page, req.Size)

	total, _ := dao.UserFollows.Ctx(ctx).Where("follower_id", req.Id).Count()
	rows := make([]followUserRow, 0)
	if total > 0 {
		offset := (page - 1) * size
		_ = dao.UserFollows.DB().Ctx(ctx).Raw(`
			SELECT u.id, u.username, u.display_name, u.avatar_id, u.bio,
			       uf.created_at AS followed_at,
			       (SELECT COUNT(*) FROM posts WHERE author_id = u.id AND status = 2 AND deleted_at IS NULL) AS article_count,
			       (SELECT COUNT(*) FROM user_follows WHERE following_id = u.id) AS follower_count
			FROM user_follows uf
			JOIN users u ON u.id = uf.following_id
			WHERE uf.follower_id = ? AND u.deleted_at IS NULL
			ORDER BY uf.created_at DESC
			LIMIT ? OFFSET ?`, req.Id, size, offset).Scan(&rows)
	}

	list := toFollowUserItems(ctx, rows)
	return &v1.UserFollowingRes{List: list, Total: total, Page: page, Size: size}, nil
}

// toFollowUserItems converts rows to API items and marks is_following_back.
func toFollowUserItems(ctx context.Context, rows []followUserRow) []v1.FollowUserItem {
	items := make([]v1.FollowUserItem, 0, len(rows))
	if len(rows) == 0 {
		return items
	}

	// Batch-resolve avatar URLs from medias table
	avatarIdSet := map[int64]bool{}
	for _, r := range rows {
		if r.AvatarId > 0 {
			avatarIdSet[r.AvatarId] = true
		}
	}
	avatarUrlMap := map[int64]string{}
	if len(avatarIdSet) > 0 {
		ids := make([]int64, 0, len(avatarIdSet))
		for id := range avatarIdSet {
			ids = append(ids, id)
		}
		type mediaRow struct {
			Id     int64  `orm:"id"`
			CdnUrl string `orm:"cdn_url"`
		}
		var mediaRows []mediaRow
		_ = dao.Medias.Ctx(ctx).Fields("id, cdn_url").WhereIn("id", ids).Scan(&mediaRows)
		for _, mr := range mediaRows {
			if mr.CdnUrl != "" {
				avatarUrlMap[mr.Id] = mr.CdnUrl
			}
		}
	}

	// Build is_following_back set from current viewer's follows
	followingSet := make(map[int64]bool)
	if selfID := optionalUserID(ctx); selfID > 0 {
		type idRow struct{ FollowingID int64 `orm:"following_id"` }
		var fwRows []idRow
		_ = dao.UserFollows.Ctx(ctx).
			Where("follower_id", selfID).Fields("following_id").Scan(&fwRows)
		for _, f := range fwRows {
			followingSet[f.FollowingID] = true
		}
	}

	for _, r := range rows {
		items = append(items, v1.FollowUserItem{
			ID:              r.ID,
			Username:        r.Username,
			DisplayName:     r.DisplayName,
			Avatar:          avatarUrlMap[r.AvatarId],
			Bio:             r.Bio,
			FollowedAt:      r.FollowedAt,
			ArticleCount:    r.ArticleCount,
			FollowerCount:   r.FollowerCount,
			IsFollowingBack: followingSet[r.ID],
		})
	}
	return items
}

func normPage(page, size int) (int, int) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}
	return page, size
}
