package taxonomy

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/taxonomy/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sTaxonomy struct{}

func New() service.ITaxonomy { return &sTaxonomy{} }

// computeRecursiveCounts sums each item's PostCount with all its descendants'.
// items is a flat list; parent-child relationship is determined by ParentId.
func computeRecursiveCounts(items []*v1.TaxonomyItem) {
	idMap := make(map[int64]*v1.TaxonomyItem, len(items))
	for _, item := range items {
		idMap[item.Id] = item
	}
	childrenOf := make(map[int64][]int64, len(items))
	for _, item := range items {
		if item.ParentId != nil && *item.ParentId != 0 {
			childrenOf[*item.ParentId] = append(childrenOf[*item.ParentId], item.Id)
		}
	}

	visited := make(map[int64]bool, len(items))

	var dfs func(id int64) int64
	dfs = func(id int64) int64 {
		if visited[id] {
			return 0
		}
		visited[id] = true
		item, ok := idMap[id]
		if !ok {
			return 0
		}
		total := int64(item.PostCount)
		for _, childId := range childrenOf[id] {
			total += dfs(childId)
		}
		item.PostCount = int(total)
		return total
	}

	for _, item := range items {
		if item.ParentId == nil || *item.ParentId == 0 {
			dfs(item.Id)
		}
	}
}

func init() {
	service.RegisterTaxonomy(New())
}

func (s *sTaxonomy) TaxonomyCreate(ctx context.Context, req *v1.TaxonomyCreateReq) (int64, error) {
	result, err := dao.Taxonomies.Ctx(ctx).Insert(g.Map{
		"id":          idgen.New(),
		"term_id":     req.TermId,
		"taxonomy":    req.Taxonomy,
		"description": req.Description,
		"parent_id":   req.ParentId,
		"extra":       req.Extra,
		"post_count":  0,
	})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	// Emit taxonomy.created
	{
		type termRow struct {
			Name string `orm:"name"`
			Slug string `orm:"slug"`
		}
		var term termRow
		_ = dao.Terms.Ctx(ctx).Where("id", req.TermId).Scan(&term)
		_ = event.Emit(ctx, event.TaxonomyCreated, payload.TaxonomyCreated{
			TaxID:    id,
			TermID:   req.TermId,
			TermName: term.Name,
			TermSlug: term.Slug,
			Taxonomy: req.Taxonomy,
		})
	}
	return id, nil
}

func (s *sTaxonomy) TaxonomyDelete(ctx context.Context, id int64) error {
	type snapRow struct {
		TermID   int64  `orm:"term_id"`
		Taxonomy string `orm:"taxonomy"`
	}
	var snap snapRow
	_ = dao.Taxonomies.Ctx(ctx).Where("id", id).Scan(&snap)

	type termRow struct {
		Name string `orm:"name"`
		Slug string `orm:"slug"`
	}
	var term termRow
	if snap.TermID > 0 {
		_ = dao.Terms.Ctx(ctx).Where("id", snap.TermID).Scan(&term)
	}

	_, err := dao.Taxonomies.Ctx(ctx).Where("id", id).Delete()
	if err != nil {
		return err
	}
	_, err = dao.ObjectTaxonomies.Ctx(ctx).Where("taxonomy_id", id).Delete()
	if err != nil {
		return err
	}
	_ = event.Emit(ctx, event.TaxonomyDeleted, payload.TaxonomyDeleted{
		TaxID:    id,
		TermName: term.Name,
		TermSlug: term.Slug,
		Taxonomy: snap.Taxonomy,
	})
	return nil
}

func (s *sTaxonomy) TaxonomyGetList(ctx context.Context, req *v1.TaxonomyGetListReq) (*v1.TaxonomyGetListRes, error) {
	total, err := dao.Taxonomies.Ctx(ctx).Where("taxonomy", req.Taxonomy).Count()
	if err != nil {
		return nil, err
	}

	type JoinRow struct {
		TaxId         int64       `orm:"tax_id"`
		TermId        int64       `orm:"term_id"`
		Taxonomy      string      `orm:"taxonomy"`
		Description   string      `orm:"description"`
		ParentId      *int64      `orm:"parent_id"`
		PostCount     int         `orm:"post_count"`
		Extra         string      `orm:"extra"`
		TermName      string      `orm:"term_name"`
		TermSlug      string      `orm:"term_slug"`
		TermCreatedAt *gtime.Time `orm:"term_created_at"`
	}

	var rows []JoinRow
	if total > 0 {
		// Fetch ALL items (no SQL pagination) so recursive count computation works correctly.
		err = dao.Taxonomies.DB().GetScan(ctx, &rows,
			`SELECT tax.id as tax_id, tax.term_id, tax.taxonomy, tax.description, tax.parent_id,
			        (SELECT COUNT(*) FROM object_taxonomies ot JOIN posts p ON p.id = ot.object_id
			         WHERE ot.taxonomy_id = tax.id AND ot.object_type = 'post' AND p.status = 2) as post_count,
			        tax.extra, t.name as term_name, t.slug as term_slug,
			        t.created_at as term_created_at
			 FROM taxonomies tax
			 LEFT JOIN terms t ON t.id = tax.term_id
			 WHERE tax.taxonomy = ?`,
			req.Taxonomy)
		if err != nil {
			return nil, err
		}
	}

	allItems := make([]*v1.TaxonomyItem, len(rows))
	for i, row := range rows {
		allItems[i] = &v1.TaxonomyItem{
			Id:          row.TaxId,
			TermId:      row.TermId,
			Taxonomy:    row.Taxonomy,
			Description: row.Description,
			ParentId:    row.ParentId,
			PostCount:   row.PostCount,
			Extra:       row.Extra,
			Term: &v1.TermItem{
				Id:        row.TermId,
				Name:      row.TermName,
				Slug:      row.TermSlug,
				CreatedAt: row.TermCreatedAt,
			},
		}
	}

	// Sum child counts into parent counts.
	computeRecursiveCounts(allItems)

	// Paginate in Go.
	offset := (req.Page - 1) * req.Size
	end := offset + req.Size
	if offset > len(allItems) {
		offset = len(allItems)
	}
	if end > len(allItems) {
		end = len(allItems)
	}
	return &v1.TaxonomyGetListRes{List: allItems[offset:end], Total: total}, nil
}

func (s *sTaxonomy) TaxonomyGetTree(ctx context.Context, taxonomy string) (*v1.TaxonomyGetTreeRes, error) {
	type JoinRow struct {
		TaxId         int64       `orm:"tax_id"`
		TermId        int64       `orm:"term_id"`
		Taxonomy      string      `orm:"taxonomy"`
		Description   string      `orm:"description"`
		ParentId      *int64      `orm:"parent_id"`
		PostCount     int         `orm:"post_count"`
		Extra         string      `orm:"extra"`
		TermName      string      `orm:"term_name"`
		TermSlug      string      `orm:"term_slug"`
		TermCreatedAt *gtime.Time `orm:"term_created_at"`
	}

	var rows []JoinRow
	err := dao.Taxonomies.DB().GetScan(ctx, &rows,
		`SELECT tax.id as tax_id, tax.term_id, tax.taxonomy, tax.description, tax.parent_id,
		        (SELECT COUNT(*) FROM object_taxonomies ot JOIN posts p ON p.id = ot.object_id
		         WHERE ot.taxonomy_id = tax.id AND ot.object_type = 'post' AND p.status = 2) as post_count,
		        tax.extra, t.name as term_name, t.slug as term_slug,
		        t.created_at as term_created_at
		 FROM taxonomies tax
		 LEFT JOIN terms t ON t.id = tax.term_id
		 WHERE tax.taxonomy = ?`,
		taxonomy)
	if err != nil {
		return nil, err
	}

	itemMap := make(map[int64]*v1.TaxonomyItem)
	for _, row := range rows {
		itemMap[row.TaxId] = &v1.TaxonomyItem{
			Id:          row.TaxId,
			TermId:      row.TermId,
			Taxonomy:    row.Taxonomy,
			Description: row.Description,
			ParentId:    row.ParentId,
			PostCount:   row.PostCount,
			Extra:       row.Extra,
			Term: &v1.TermItem{
				Id:        row.TermId,
				Name:      row.TermName,
				Slug:      row.TermSlug,
				CreatedAt: row.TermCreatedAt,
			},
			Children: []*v1.TaxonomyItem{},
		}
	}

	// Compute recursive post counts (parent includes all descendants) before building tree.
	flatItems := make([]*v1.TaxonomyItem, 0, len(itemMap))
	for _, item := range itemMap {
		flatItems = append(flatItems, item)
	}
	computeRecursiveCounts(flatItems)

	var roots []*v1.TaxonomyItem
	for _, row := range rows {
		item := itemMap[row.TaxId]
		if row.ParentId == nil || *row.ParentId == 0 {
			roots = append(roots, item)
		} else if parent, ok := itemMap[*row.ParentId]; ok {
			parent.Children = append(parent.Children, item)
		} else {
			roots = append(roots, item)
		}
	}

	if roots == nil {
		roots = []*v1.TaxonomyItem{}
	}
	return &v1.TaxonomyGetTreeRes{List: roots}, nil
}

func (s *sTaxonomy) TaxonomyUpdate(ctx context.Context, req *v1.TaxonomyUpdateReq) error {
	data := g.Map{}
	if req.Description != nil {
		data["description"] = *req.Description
	}
	if req.ParentId != nil {
		if *req.ParentId == 0 {
			data["parent_id"] = nil
		} else {
			data["parent_id"] = *req.ParentId
		}
	}
	if req.Extra != nil {
		data["extra"] = *req.Extra
	}
	if len(data) == 0 {
		return nil
	}
	_, err := dao.Taxonomies.Ctx(ctx).Where("id", req.Id).Update(data)
	return err
}

func (s *sTaxonomy) TermCreate(ctx context.Context, req *v1.TermCreateReq) (int64, error) {
	count, err := dao.Terms.Ctx(ctx).Where("slug", req.Slug).Count()
	if err != nil {
		return 0, err
	}
	if count > 0 {
		return 0, gerror.New(g.I18n().Tf(ctx, "taxonomy.slug_taken", req.Slug))
	}
	result, err := dao.Terms.Ctx(ctx).Insert(g.Map{
		"id":   idgen.New(),
		"name": req.Name,
		"slug": req.Slug,
	})
	if err != nil {
		return 0, err
	}
	id, _ := result.LastInsertId()
	_ = event.Emit(ctx, event.TermCreated, payload.TermCreated{
		TermID: id,
		Name:   req.Name,
		Slug:   req.Slug,
	})
	return id, nil
}

func (s *sTaxonomy) TermDelete(ctx context.Context, id int64) error {
	type snapRow struct {
		Name string `orm:"name"`
		Slug string `orm:"slug"`
	}
	var snap snapRow
	_ = dao.Terms.Ctx(ctx).Where("id", id).Scan(&snap)

	_, _ = dao.ObjectTaxonomies.Ctx(ctx).WhereIn(
		"taxonomy_id",
		dao.Taxonomies.Ctx(ctx).Fields("id").Where("term_id", id),
	).Delete()
	_, _ = dao.Taxonomies.Ctx(ctx).Where("term_id", id).Delete()
	_, err := dao.Terms.Ctx(ctx).Where("id", id).Delete()
	if err != nil {
		return err
	}
	_ = event.Emit(ctx, event.TermDeleted, payload.TermDeleted{
		TermID: id,
		Name:   snap.Name,
		Slug:   snap.Slug,
	})
	return nil
}

func (s *sTaxonomy) TermUpdate(ctx context.Context, req *v1.TermUpdateReq) error {
	data := g.Map{}
	if req.Name != nil {
		data["name"] = *req.Name
	}
	if req.Slug != nil {
		count, err := dao.Terms.Ctx(ctx).Where("slug", *req.Slug).WhereNot("id", req.Id).Count()
		if err != nil {
			return err
		}
		if count > 0 {
			return gerror.New(g.I18n().Tf(ctx, "taxonomy.slug_taken", *req.Slug))
		}
		data["slug"] = *req.Slug
	}
	if len(data) == 0 {
		return nil
	}
	_, err := dao.Terms.Ctx(ctx).Where("id", req.Id).Update(data)
	return err
}

func (s *sTaxonomy) ObjectTaxonomyBind(ctx context.Context, req *v1.ObjectTaxonomyBindReq) error {
	for _, taxId := range req.TaxonomyIds {
		_, _ = dao.ObjectTaxonomies.Ctx(ctx).Insert(g.Map{
			"object_id":   req.ObjectId,
			"object_type": req.ObjectType,
			"taxonomy_id": taxId,
		})
	}
	return nil
}

func (s *sTaxonomy) ObjectTaxonomyUnbind(ctx context.Context, req *v1.ObjectTaxonomyUnbindReq) error {
	_, err := dao.ObjectTaxonomies.Ctx(ctx).
		Where("object_id", req.ObjectId).
		Where("object_type", req.ObjectType).
		WhereIn("taxonomy_id", req.TaxonomyIds).
		Delete()
	return err
}

func (s *sTaxonomy) ObjectTaxonomyGet(ctx context.Context, req *v1.ObjectTaxonomyGetReq) (*v1.ObjectTaxonomyGetRes, error) {
	type JoinRow struct {
		TaxId         int64       `orm:"tax_id"`
		TermId        int64       `orm:"term_id"`
		Taxonomy      string      `orm:"taxonomy"`
		Description   string      `orm:"description"`
		ParentId      *int64      `orm:"parent_id"`
		PostCount     int         `orm:"post_count"`
		Extra         string      `orm:"extra"`
		TermName      string      `orm:"term_name"`
		TermSlug      string      `orm:"term_slug"`
		TermCreatedAt *gtime.Time `orm:"term_created_at"`
	}

	query := `SELECT tax.id as tax_id, tax.term_id, tax.taxonomy, tax.description, tax.parent_id,
	                 (SELECT COUNT(*) FROM object_taxonomies ot2 JOIN posts p ON p.id = ot2.object_id
	                  WHERE ot2.taxonomy_id = tax.id AND ot2.object_type = 'post' AND p.status = 2) as post_count,
	                 tax.extra, t.name as term_name, t.slug as term_slug,
	                 t.created_at as term_created_at
	          FROM object_taxonomies ot
	          JOIN taxonomies tax ON tax.id = ot.taxonomy_id
	          LEFT JOIN terms t ON t.id = tax.term_id
	          WHERE ot.object_id = ? AND ot.object_type = ?`
	args := []interface{}{req.ObjectId, req.ObjectType}
	if req.Taxonomy != "" {
		query += " AND tax.taxonomy = ?"
		args = append(args, req.Taxonomy)
	}

	var rows []JoinRow
	err := dao.Taxonomies.DB().GetScan(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}

	list := make([]*v1.TaxonomyItem, len(rows))
	for i, row := range rows {
		list[i] = &v1.TaxonomyItem{
			Id:          row.TaxId,
			TermId:      row.TermId,
			Taxonomy:    row.Taxonomy,
			Description: row.Description,
			ParentId:    row.ParentId,
			PostCount:   row.PostCount,
			Extra:       row.Extra,
			Term: &v1.TermItem{
				Id:        row.TermId,
				Name:      row.TermName,
				Slug:      row.TermSlug,
				CreatedAt: row.TermCreatedAt,
			},
		}
	}
	return &v1.ObjectTaxonomyGetRes{List: list}, nil
}
