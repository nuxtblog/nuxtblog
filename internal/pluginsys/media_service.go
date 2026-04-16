package pluginsys

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"sync"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/storage"
	"github.com/nuxtblog/nuxtblog/sdk"
)

// pluginAdapterNames tracks which storage adapter name each plugin registered,
// so we can unregister by the correct name (not the plugin ID).
var (
	pluginAdapterMu    sync.Mutex
	pluginAdapterNames = make(map[string]string) // pluginID → adapter name
)

// mediaService implements sdk.MediaService, bridging the plugin SDK
// to the internal media and storage logic.
type mediaService struct {
	pluginID string
}

func newMediaService(pluginID string) sdk.MediaService {
	return &mediaService{pluginID: pluginID}
}

func (ms *mediaService) RegisterStorageAdapter(name string, storageType int, adapter sdk.StorageAdapter) error {
	g.Log().Infof(context.Background(), "[pluginmgr] plugin %s registering storage adapter: %s (type=%d)", ms.pluginID, name, storageType)
	if err := storage.RegisterAdapter(name, storageType, adapter); err != nil {
		return err
	}
	pluginAdapterMu.Lock()
	pluginAdapterNames[ms.pluginID] = name
	pluginAdapterMu.Unlock()
	return nil
}

func (ms *mediaService) RegisterCategory(def sdk.CategoryDef) error {
	ctx := context.Background()
	g.Log().Infof(ctx, "[pluginmgr] plugin %s registering media category: %s", ms.pluginID, def.Slug)
	// Use options table directly to avoid import cycle with logic/media.
	return registerPluginCategoryViaOptions(ctx, ms.pluginID, def)
}

func (ms *mediaService) Upload(ctx context.Context, data []byte, filename string, opts sdk.UploadOpts) (*sdk.UploadResult, error) {
	up := storage.Default(ctx)
	result, err := storage.UploadFromBytes(ctx, up, data, filename, storage.UploadOptions{
		Category:     opts.Category,
		PathTemplate: opts.PathTemplate,
	})
	if err != nil {
		return nil, fmt.Errorf("plugin upload: %w", err)
	}
	return &sdk.UploadResult{
		StorageType: result.StorageType,
		StorageKey:  result.StorageKey,
		CdnUrl:      result.CdnUrl,
		MimeType:    result.MimeType,
		FileSize:    result.FileSize,
		Width:       result.Width,
		Height:      result.Height,
		Duration:    result.Duration,
		Variants:    result.Variants,
	}, nil
}

func (ms *mediaService) Delete(ctx context.Context, mediaID int64) error {
	return service.Media().Delete(ctx, mediaID)
}

// unregisterPluginMedia cleans up storage adapters and categories registered by a plugin.
func unregisterPluginMedia(pluginID string) {
	pluginAdapterMu.Lock()
	adapterName, ok := pluginAdapterNames[pluginID]
	if ok {
		delete(pluginAdapterNames, pluginID)
	}
	pluginAdapterMu.Unlock()
	if ok {
		storage.UnregisterAdapter(adapterName)
	}
	_ = unregisterPluginCategoriesViaOptions(context.Background(), pluginID)
}

// ── Options-based category manipulation (avoids import cycle) ─────────────────

type pluginCatItem struct {
	Slug         string `json:"slug"`
	LabelZh      string `json:"label_zh"`
	LabelEn      string `json:"label_en"`
	IsSystem     bool   `json:"is_system"`
	Order        int    `json:"order"`
	StorageKey   string `json:"storage_key"`
	PluginID     string `json:"plugin_id,omitempty"`
	MaxPerOwner  int    `json:"max_per_owner,omitempty"`
	FormatPolicy string `json:"format_policy,omitempty"`
}

func readCategoriesFromOptions(ctx context.Context) ([]pluginCatItem, error) {
	val, err := g.DB().Ctx(ctx).Model("options").Where("key", "media_categories").Value("value")
	if err != nil || val.IsNil() {
		return nil, err
	}
	raw := val.String()
	var cats []pluginCatItem
	// Historical data may be double-encoded (JSON string wrapping JSON array)
	// due to an earlier serialization bug. Try unwrapping first, then direct parse.
	var inner string
	if json.Unmarshal([]byte(raw), &inner) == nil {
		_ = json.Unmarshal([]byte(inner), &cats)
	} else {
		_ = json.Unmarshal([]byte(raw), &cats)
	}
	return cats, nil
}

func saveCategoriesViaOptions(ctx context.Context, cats []pluginCatItem) error {
	slices.SortFunc(cats, func(a, b pluginCatItem) int { return a.Order - b.Order })
	data, err := json.Marshal(cats)
	if err != nil {
		return err
	}
	cnt, _ := g.DB().Ctx(ctx).Model("options").Where("key", "media_categories").Count()
	if cnt == 0 {
		_, err = g.DB().Ctx(ctx).Model("options").Data(g.Map{
			"key": "media_categories", "value": string(data), "autoload": 1,
		}).Insert()
	} else {
		_, err = g.DB().Ctx(ctx).Model("options").Where("key", "media_categories").Data(g.Map{"value": string(data)}).Update()
	}
	return err
}

func registerPluginCategoryViaOptions(ctx context.Context, pluginID string, def sdk.CategoryDef) error {
	cats, _ := readCategoriesFromOptions(ctx)
	found := false
	for i := range cats {
		if cats[i].Slug == def.Slug {
			cats[i].LabelZh = def.ResolvedZh
			cats[i].LabelEn = def.ResolvedEn
			cats[i].Order = def.Order
			cats[i].MaxPerOwner = def.MaxPerOwner
			cats[i].PluginID = pluginID
			cats[i].IsSystem = false
			found = true
			break
		}
	}
	if !found {
		cats = append(cats, pluginCatItem{
			Slug: def.Slug, LabelZh: def.ResolvedZh, LabelEn: def.ResolvedEn,
			Order: def.Order, MaxPerOwner: def.MaxPerOwner, PluginID: pluginID,
		})
	}
	return saveCategoriesViaOptions(ctx, cats)
}

func unregisterPluginCategoriesViaOptions(ctx context.Context, pluginID string) error {
	cats, _ := readCategoriesFromOptions(ctx)
	filtered := cats[:0]
	for _, c := range cats {
		if c.PluginID != pluginID {
			filtered = append(filtered, c)
		}
	}
	return saveCategoriesViaOptions(ctx, filtered)
}
