package reaction

import (
	"context"
	"fmt"

	reactionv1 "github.com/nuxtblog/nuxtblog/api/reaction/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sReaction struct{}

func New() service.IReaction { return &sReaction{} }

func init() {
	service.RegisterReaction(New())
}

func (s *sReaction) Like(ctx context.Context, postId, userId int64) error {
	// Check if already liked
	count, _ := dao.UserLikes.Ctx(ctx).
		Where("user_id", userId).Where("object_type", "post").Where("object_id", postId).Count()
	if count > 0 {
		return nil // already liked, idempotent
	}
	now := gtime.Now()
	_, err := dao.UserLikes.Ctx(ctx).Data(g.Map{
		"user_id":     userId,
		"object_type": "post",
		"object_id":   postId,
		"created_at":  now,
	}).Insert()
	if err != nil {
		return err
	}
	// Atomically increment like_count
	_, _ = dao.PostStats.Ctx(ctx).Where("post_id", postId).Increment("like_count", 1)
	_ = event.Emit(ctx, event.ReactionAdded, payload.ReactionAdded{
		UserID: userId, ObjectType: "post", ObjectID: postId, Type: "like",
	})
	return nil
}

func (s *sReaction) Unlike(ctx context.Context, postId, userId int64) error {
	result, err := dao.UserLikes.Ctx(ctx).
		Where("user_id", userId).Where("object_type", "post").Where("object_id", postId).Delete()
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected > 0 {
		// Decrement like_count but not below 0
		_, _ = dao.PostStats.Ctx(ctx).
			Where("post_id", postId).Where("like_count > 0").
			Decrement("like_count", 1)
		_ = event.Emit(ctx, event.ReactionRemoved, payload.ReactionRemoved{
			UserID: userId, ObjectType: "post", ObjectID: postId, Type: "like",
		})
	}
	return nil
}

func (s *sReaction) Bookmark(ctx context.Context, postId, userId int64) error {
	count, _ := dao.UserBookmarks.Ctx(ctx).
		Where("user_id", userId).Where("object_type", "post").Where("object_id", postId).Count()
	if count > 0 {
		return nil
	}
	now := gtime.Now()
	_, err := dao.UserBookmarks.Ctx(ctx).Data(g.Map{
		"user_id":     userId,
		"object_type": "post",
		"object_id":   postId,
		"created_at":  now,
	}).Insert()
	if err != nil {
		return err
	}
	_ = event.Emit(ctx, event.ReactionAdded, payload.ReactionAdded{
		UserID: userId, ObjectType: "post", ObjectID: postId, Type: "bookmark",
	})
	return nil
}

func (s *sReaction) Unbookmark(ctx context.Context, postId, userId int64) error {
	result, err := dao.UserBookmarks.Ctx(ctx).
		Where("user_id", userId).Where("object_type", "post").Where("object_id", postId).Delete()
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected > 0 {
		_ = event.Emit(ctx, event.ReactionRemoved, payload.ReactionRemoved{
			UserID: userId, ObjectType: "post", ObjectID: postId, Type: "bookmark",
		})
	}
	return nil
}

func (s *sReaction) GetStatus(ctx context.Context, postId, userId int64) (*reactionv1.GetReactionRes, error) {
	likeCount, _ := dao.UserLikes.Ctx(ctx).
		Where("user_id", userId).Where("object_type", "post").Where("object_id", postId).Count()
	bookmarkCount, _ := dao.UserBookmarks.Ctx(ctx).
		Where("user_id", userId).Where("object_type", "post").Where("object_id", postId).Count()
	return &reactionv1.GetReactionRes{
		Liked:      likeCount > 0,
		Bookmarked: bookmarkCount > 0,
	}, nil
}

func (s *sReaction) GetBookmarks(ctx context.Context, userId int64, page, size int) (*reactionv1.GetBookmarksRes, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 50 {
		size = 20
	}

	total, err := dao.UserBookmarks.Ctx(ctx).
		Where("user_id", userId).Where("object_type", "post").Count()
	if err != nil {
		return nil, err
	}
	if total == 0 {
		return &reactionv1.GetBookmarksRes{List: []*reactionv1.BookmarkPostItem{}, Total: 0, Page: page, Size: size}, nil
	}

	type BookmarkRow struct {
		ObjectId int64 `orm:"object_id"`
	}
	var bRows []BookmarkRow
	err = dao.UserBookmarks.Ctx(ctx).
		Where("user_id", userId).
		Where("object_type", "post").
		Fields("object_id").
		OrderDesc("created_at").
		Page(page, size).
		Scan(&bRows)
	if err != nil {
		return nil, err
	}

	postIds := make([]int64, 0, len(bRows))
	for _, r := range bRows {
		postIds = append(postIds, r.ObjectId)
	}

	type PostRow struct {
		Id        int64  `orm:"id"         json:"id"`
		Title     string `orm:"title"      json:"title"`
		Slug      string `orm:"slug"       json:"slug"`
		Excerpt   string `orm:"excerpt"    json:"excerpt"`
		CreatedAt string `orm:"created_at" json:"created_at"`
	}
	var posts []PostRow
	err = dao.Posts.Ctx(ctx).
		WhereIn("id", postIds).
		WhereNull("deleted_at").
		Fields("id,title,slug,excerpt,created_at").
		Scan(&posts)
	if err != nil {
		return nil, fmt.Errorf("failed to load posts: %w", err)
	}

	postMap := make(map[int64]*PostRow)
	for i := range posts {
		postMap[posts[i].Id] = &posts[i]
	}

	items := make([]*reactionv1.BookmarkPostItem, 0, len(bRows))
	for _, r := range bRows {
		if p, ok := postMap[r.ObjectId]; ok {
			items = append(items, &reactionv1.BookmarkPostItem{
				Id:        p.Id,
				Title:     p.Title,
				Slug:      p.Slug,
				Excerpt:   p.Excerpt,
				CreatedAt: p.CreatedAt,
			})
		}
	}

	return &reactionv1.GetBookmarksRes{
		List:  items,
		Total: total,
		Page:  page,
		Size:  size,
	}, nil
}
