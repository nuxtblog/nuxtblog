// Singleton module-level state so multiple components share one fetch.
export interface MediaCategoryItem {
  slug:          string
  label_zh:      string
  label_en:      string
  is_system:     boolean
  order:         number
  storage_key:   string
  plugin_id?:    string
  max_per_owner?: number
  format_policy?: string
  path_template?: string
}

const _cats    = ref<MediaCategoryItem[]>([])
const _loading = ref(false)
const _loaded  = ref(false)

export const useMediaCategories = () => {
  const { apiFetch } = useApiFetch()
  const { locale }   = useI18n()

  const load = async () => {
    if (_loaded.value || _loading.value) return
    _loading.value = true
    try {
      const res  = await apiFetch<{ list: MediaCategoryItem[] }>('/admin/media/categories')
      _cats.value   = res.list ?? []
      _loaded.value = true
    } catch {
      // non-fatal: pages fall back to empty list
    } finally {
      _loading.value = false
    }
  }

  /** Force a re-fetch (call after mutations). */
  const refresh = async () => {
    _loaded.value = false
    await load()
  }

  /** Return the localised label for a slug, falling back to the slug itself. */
  const getCategoryLabel = (slug: string): string => {
    const cat = _cats.value.find(c => c.slug === slug)
    if (!cat) return slug
    const isZh = locale.value.startsWith('zh')
    return (isZh ? cat.label_zh : cat.label_en) || cat.label_zh || cat.label_en || slug
  }

  return {
    categories:       _cats      as Readonly<typeof _cats>,
    categoriesLoading: _loading  as Readonly<typeof _loading>,
    load,
    refresh,
    getCategoryLabel,
  }
}
