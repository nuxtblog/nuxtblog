// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"context"
)

type (
	IPost interface {
		Create(ctx context.Context, req *v1.PostCreateReq) (id int64, err error)
		Update(ctx context.Context, req *v1.PostUpdateReq) error
		// Delete permanently removes a post (force delete).
		Delete(ctx context.Context, id int64) error
		// Trash soft-deletes a post by setting deleted_at.
		Trash(ctx context.Context, id int64) error
		// Restore clears deleted_at, recovering a trashed post.
		Restore(ctx context.Context, id int64) error
		// Batch performs a bulk action on a list of post IDs.
		Batch(ctx context.Context, req *v1.PostBatchReq) (int, error)
		// BatchUpdate applies field-level updates (cover, status, tags, author) to multiple posts.
		BatchUpdate(ctx context.Context, req *v1.PostBatchUpdateReq) (int, error)
		// UpsertMetas 批量写入/删除 metas。
		// value 非空 = upsert；value 为空字符串 = 删除该 key。
		UpsertMetas(ctx context.Context, postId int64, metas map[string]string) error
		// GetMetas 查询单篇文章的所有 metas，无数据返回 nil。
		GetMetas(ctx context.Context, postId int64) (map[string]string, error)
		GetById(ctx context.Context, id int64) (item *v1.PostDetailItem, err error)
		GetBySlug(ctx context.Context, slug string) (*v1.PostDetailEnrichedItem, error)
		GetList(ctx context.Context, req *v1.PostGetListReq) (*v1.PostGetListRes, error)
		RevisionList(ctx context.Context, req *v1.PostRevisionListReq) (*v1.PostRevisionListRes, error)
		RevisionRestore(ctx context.Context, postId int64, revisionId int64) error
		SeoUpdate(ctx context.Context, req *v1.PostSeoUpdateReq) error
		// GetStats returns aggregated statistics for all regular posts (post_type=1).
		GetStats(ctx context.Context) (*v1.GetStatsRes, error)
		// IncrementView increments the post view count and records browse history.
		IncrementView(ctx context.Context, id int64) error
		// VerifyPassword checks whether the given password matches the post's password hash.
		VerifyPassword(ctx context.Context, id int64, password string) (bool, error)
		// PublishScheduled auto-publishes draft posts whose published_at has arrived.
		// Intended to be called from a periodic cron job (e.g. every minute).
		PublishScheduled(ctx context.Context)
	}
)

var (
	localPost IPost
)

func Post() IPost {
	if localPost == nil {
		panic("implement not found for interface IPost, forgot register?")
	}
	return localPost
}

func RegisterPost(i IPost) {
	localPost = i
}
