package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
)

// GetStats returns aggregated statistics for all regular posts (post_type=1).
func (s *sPost) GetStats(ctx context.Context) (*v1.GetStatsRes, error) {
	// Count posts by status (post_type=1, not soft-deleted)
	type StatusCount struct {
		Status int   `orm:"status"`
		Count  int64 `orm:"count"`
	}
	var statusCounts []StatusCount
	err := dao.Posts.Ctx(ctx).
		Fields("status, COUNT(*) AS count").
		WhereNull("deleted_at").
		Where("post_type", int(v1.PostTypePost)).
		WhereNot("status", int(v1.PostStatusTrashed)).
		Group("status").
		Scan(&statusCounts)
	if err != nil {
		return nil, err
	}

	res := &v1.GetStatsRes{}
	for _, sc := range statusCounts {
		res.TotalPosts += sc.Count
		switch v1.PostStatus(sc.Status) {
		case v1.PostStatusDraft:
			res.DraftPosts = sc.Count
		case v1.PostStatusPublished:
			res.PublishedPosts = sc.Count
		case v1.PostStatusPrivate:
			res.PrivatePosts = sc.Count
		case v1.PostStatusArchived:
			res.ArchivedPosts = sc.Count
		}
	}

	// Sum view_count, like_count, comment_count from post_stats joined with posts
	type AggStats struct {
		TotalViews    int64 `orm:"total_views"`
		TotalLikes    int64 `orm:"total_likes"`
		TotalComments int64 `orm:"total_comments"`
	}
	var agg AggStats
	err = dao.Posts.DB().Ctx(ctx).
		Raw(`SELECT
			COALESCE(SUM(ps.view_count), 0)    AS total_views,
			COALESCE(SUM(ps.like_count), 0)    AS total_likes,
			COALESCE(SUM(ps.comment_count), 0) AS total_comments
		FROM post_stats ps
		JOIN posts p ON ps.post_id = p.id
		WHERE p.deleted_at IS NULL
		  AND p.post_type = ?`, int(v1.PostTypePost)).
		Scan(&agg)
	if err != nil {
		return nil, err
	}

	res.TotalViews = agg.TotalViews
	res.TotalLikes = agg.TotalLikes
	res.TotalComments = agg.TotalComments

	return res, nil
}
