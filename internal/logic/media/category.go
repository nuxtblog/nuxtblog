package media

import (
	"context"
	"encoding/json"
	"errors"
	"slices"

	"github.com/gogf/gf/v2/frame/g"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/consts"
	"github.com/nuxtblog/nuxtblog/internal/dao"
)

// CategoryItem is the JSON-serialisable representation stored in options["media_categories"].
type CategoryItem struct {
	Slug       string `json:"slug"`
	LabelZh    string `json:"label_zh"`
	LabelEn    string `json:"label_en"`
	IsSystem   bool   `json:"is_system"`
	Order      int    `json:"order"`
	StorageKey string `json:"storage_key"` // empty = use system default backend
}

const optKeyCategories = "media_categories"

// getCategories reads all categories from options, sorted by Order.
// Falls back to built-in defaults when the option is absent or empty.
func getCategories(ctx context.Context) ([]CategoryItem, error) {
	val, err := dao.Options.Ctx(ctx).Where("key", optKeyCategories).Value("value")
	if err != nil || val.IsNil() {
		return builtinItems(), nil
	}
	raw := val.String()
	var cats []CategoryItem
	var inner string
	if json.Unmarshal([]byte(raw), &inner) == nil {
		_ = json.Unmarshal([]byte(inner), &cats)
	} else {
		_ = json.Unmarshal([]byte(raw), &cats)
	}
	if len(cats) == 0 {
		return builtinItems(), nil
	}
	slices.SortFunc(cats, func(a, b CategoryItem) int { return a.Order - b.Order })
	return cats, nil
}

// saveCategories persists the category list to options (upsert).
func saveCategories(ctx context.Context, cats []CategoryItem) error {
	slices.SortFunc(cats, func(a, b CategoryItem) int { return a.Order - b.Order })
	data, err := json.Marshal(cats)
	if err != nil {
		return err
	}
	cnt, _ := dao.Options.Ctx(ctx).Where("key", optKeyCategories).Count()
	if cnt == 0 {
		_, err = dao.Options.Ctx(ctx).Data(g.Map{
			"key": optKeyCategories, "value": string(data), "autoload": 1,
		}).Insert()
	} else {
		_, err = dao.Options.Ctx(ctx).Where("key", optKeyCategories).Data(g.Map{"value": string(data)}).Update()
	}
	return err
}

// GetCategoryStorageKey returns the configured backend name for a slug ("" = use default).
func GetCategoryStorageKey(ctx context.Context, slug string) string {
	cats, _ := getCategories(ctx)
	for _, c := range cats {
		if c.Slug == slug {
			return c.StorageKey
		}
	}
	return ""
}

// SyncBuiltinCategories upserts consts.BuiltinMediaCategories into options["media_categories"].
// It is idempotent and safe to call on every startup.
// - New built-in slugs are inserted.
// - Existing built-in slugs: IsSystem/Order/Labels are refreshed from code; StorageKey is preserved.
// - Custom slugs are kept untouched; their IsSystem flag is forced to false.
func SyncBuiltinCategories(ctx context.Context) error {
	existing, _ := getCategories(ctx)

	bySlug := make(map[string]*CategoryItem, len(existing))
	for i := range existing {
		bySlug[existing[i].Slug] = &existing[i]
	}

	builtinSlugs := make(map[string]struct{}, len(consts.BuiltinMediaCategories))
	for _, def := range consts.BuiltinMediaCategories {
		builtinSlugs[def.Slug] = struct{}{}
		if item, ok := bySlug[def.Slug]; ok {
			item.IsSystem = true
			item.Order   = def.Order
			item.LabelZh = def.LabelZh
			item.LabelEn = def.LabelEn
		} else {
			existing = append(existing, CategoryItem{
				Slug: def.Slug, LabelZh: def.LabelZh, LabelEn: def.LabelEn,
				IsSystem: true, Order: def.Order,
			})
		}
	}
	// Ensure custom categories are never marked system
	for i := range existing {
		if _, isBuiltin := builtinSlugs[existing[i].Slug]; !isBuiltin {
			existing[i].IsSystem = false
		}
	}
	return saveCategories(ctx, existing)
}

// builtinItems converts consts into CategoryItem slice (fallback when options is empty).
func builtinItems() []CategoryItem {
	items := make([]CategoryItem, len(consts.BuiltinMediaCategories))
	for i, def := range consts.BuiltinMediaCategories {
		items[i] = CategoryItem{
			Slug: def.Slug, LabelZh: def.LabelZh, LabelEn: def.LabelEn,
			IsSystem: true, Order: def.Order,
		}
	}
	return items
}

// toAPIItem converts a CategoryItem to the API response type.
func toAPIItem(c CategoryItem) v1.MediaCategoryItem {
	return v1.MediaCategoryItem{
		Slug: c.Slug, LabelZh: c.LabelZh, LabelEn: c.LabelEn,
		IsSystem: c.IsSystem, Order: c.Order, StorageKey: c.StorageKey,
	}
}

// ── IMedia method implementations ─────────────────────────────────────────────

func (s *sMedia) GetCategories(ctx context.Context) ([]v1.MediaCategoryItem, error) {
	cats, err := getCategories(ctx)
	if err != nil {
		return nil, err
	}
	items := make([]v1.MediaCategoryItem, len(cats))
	for i, c := range cats {
		items[i] = toAPIItem(c)
	}
	return items, nil
}

func (s *sMedia) UpdateCategory(ctx context.Context, req *v1.MediaCategoryUpdateReq) error {
	cats, err := getCategories(ctx)
	if err != nil {
		return err
	}
	for i := range cats {
		if cats[i].Slug == req.Slug {
			cats[i].StorageKey = req.StorageKey
			return saveCategories(ctx, cats)
		}
	}
	return errors.New("category not found")
}
