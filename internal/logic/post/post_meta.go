package post

import (
	"context"

	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/frame/g"
)

// UpsertMetas 批量写入/删除 metas。
// value 非空 = upsert；value 为空字符串 = 删除该 key。
func (s *sPost) UpsertMetas(ctx context.Context, postId int64, metas map[string]string) error {
	for key, val := range metas {
		if val == "" {
			_, err := dao.PostMetas.Ctx(ctx).
				Where("post_id", postId).
				Where("meta_key", key).
				Delete()
			if err != nil {
				return err
			}
		} else {
			_, err := dao.PostMetas.Ctx(ctx).
				Data(g.Map{
					"id":         idgen.New(),
					"post_id":    postId,
					"meta_key":   key,
					"meta_value": val,
				}).
				OnConflict("post_id", "meta_key").
				Save()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetMetas 查询单篇文章的所有 metas，无数据返回 nil。
func (s *sPost) GetMetas(ctx context.Context, postId int64) (map[string]string, error) {
	type MetaRow struct {
		MetaKey   string `orm:"meta_key"`
		MetaValue string `orm:"meta_value"`
	}
	var rows []MetaRow
	err := dao.PostMetas.Ctx(ctx).
		Fields("meta_key, meta_value").
		Where("post_id", postId).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, nil
	}
	m := make(map[string]string, len(rows))
	for _, r := range rows {
		m[r.MetaKey] = r.MetaValue
	}
	return m, nil
}
