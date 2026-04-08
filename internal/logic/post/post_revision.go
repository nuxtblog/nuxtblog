package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (s *sPost) RevisionList(ctx context.Context, req *v1.PostRevisionListReq) (*v1.PostRevisionListRes, error) {
	type RevRow struct {
		Id        int64  `orm:"id"`
		PostId    int64  `orm:"post_id"`
		AuthorId  int64  `orm:"author_id"`
		Title     string `orm:"title"`
		RevNote   string `orm:"rev_note"`
		CreatedAt string `orm:"created_at"`
	}

	m := dao.PostRevisions.Ctx(ctx).Where("post_id", req.Id)

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var rows []RevRow
	if total > 0 {
		if err = m.Page(req.Page, req.Size).OrderDesc("created_at").Scan(&rows); err != nil {
			return nil, err
		}
	}

	list := make([]*v1.PostRevisionItem, len(rows))
	for i, r := range rows {
		list[i] = &v1.PostRevisionItem{
			Id:       r.Id,
			PostId:   r.PostId,
			AuthorId: r.AuthorId,
			Title:    r.Title,
			RevNote:  r.RevNote,
		}
	}

	return &v1.PostRevisionListRes{List: list, Total: total}, nil
}

func (s *sPost) RevisionRestore(ctx context.Context, postId, revisionId int64) error {
	type RevRow struct {
		Title   string `orm:"title"`
		Content string `orm:"content"`
	}
	var rev RevRow
	err := dao.PostRevisions.Ctx(ctx).
		Where("id", revisionId).
		Where("post_id", postId).
		Scan(&rev)
	if err != nil {
		return err
	}
	if rev.Title == "" {
		return gerror.NewCode(gcode.CodeNotFound, g.I18n().T(ctx, "error.revision_not_found"))
	}

	_, err = dao.Posts.Ctx(ctx).
		Where("id", postId).
		Data(g.Map{"title": rev.Title, "content": rev.Content}).
		Update()
	return err
}
