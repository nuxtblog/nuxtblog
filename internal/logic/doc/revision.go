package doc

import (
	"context"

	docv1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func (s *sDoc) DocRevisionList(ctx context.Context, req *docv1.DocRevisionListReq) (*docv1.DocRevisionListRes, error) {
	m := dao.DocRevisions.Ctx(ctx).Where("doc_id", req.Id)
	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	size := req.Size
	if size < 1 {
		size = 10
	}

	type revRow struct {
		Id        int64       `orm:"id"`
		DocId     int64       `orm:"doc_id"`
		AuthorId  int64       `orm:"author_id"`
		Title     string      `orm:"title"`
		RevNote   string      `orm:"rev_note"`
		CreatedAt *gtime.Time `orm:"created_at"`
	}
	var rows []revRow
	if total > 0 {
		if err = m.Fields("id,doc_id,author_id,title,rev_note,created_at").
			OrderDesc("id").Page(req.Page, size).Scan(&rows); err != nil {
			return nil, err
		}
	}

	list := make([]*docv1.DocRevisionItem, 0, len(rows))
	for _, r := range rows {
		list = append(list, &docv1.DocRevisionItem{
			Id:        r.Id,
			DocId:     r.DocId,
			AuthorId:  r.AuthorId,
			Title:     r.Title,
			RevNote:   r.RevNote,
			CreatedAt: r.CreatedAt,
		})
	}
	return &docv1.DocRevisionListRes{List: list, Total: total}, nil
}

func (s *sDoc) DocRevisionRestore(ctx context.Context, docId, revisionId int64) error {
	type revRow struct {
		Title   string `orm:"title"`
		Content string `orm:"content"`
	}
	var rev revRow
	if err := dao.DocRevisions.Ctx(ctx).Where("id", revisionId).Where("doc_id", docId).Scan(&rev); err != nil {
		return err
	}
	if rev.Title == "" {
		return gerror.NewCode(gcode.CodeNotFound, "revision not found")
	}

	uid, _ := middleware.GetCurrentUserID(ctx)

	// Save current state as a revision before restoring
	type snapRow struct {
		Title   string `orm:"title"`
		Content string `orm:"content"`
	}
	var snap snapRow
	if e := dao.Docs.Ctx(ctx).Where("id", docId).WhereNull("deleted_at").Scan(&snap); e == nil && snap.Title != "" {
		_, _ = dao.DocRevisions.Ctx(ctx).Data(g.Map{
			"id":        idgen.New(),
			"doc_id":    docId,
			"author_id": uid,
			"title":     snap.Title,
			"content":   snap.Content,
			"rev_note":  "pre-restore",
		}).Insert()
	}

	_, err := dao.Docs.Ctx(ctx).Where("id", docId).WhereNull("deleted_at").Data(g.Map{
		"title":   rev.Title,
		"content": rev.Content,
	}).Update()
	return err
}
