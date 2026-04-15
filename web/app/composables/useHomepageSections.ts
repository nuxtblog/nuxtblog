export interface SectionActionConfig {
  enabled: boolean
  label?: string        // custom button text
  href?: string         // custom URL (overrides default)
  categorySlug?: string // for latest: link to category page
}

export interface SectionConfig {
  id: string
  label: string
  enabled: boolean
  count: number
  title?: string
  layout?: string
  includeCategoryIds?: number[]
  excludeCategoryIds?: number[]
  action?: SectionActionConfig
  loadMoreEnabled?: boolean
  isPlugin?: boolean
  pluginId?: string
  component?: string
  module?: string
  pluginSettings?: Record<string, unknown>
}

/** label values are i18n keys — call $t(layout.label) in templates */
export const LAYOUT_OPTIONS = [
  { label: 'admin.settings.homepage.layout_grid',     value: 'grid' },
  { label: 'admin.settings.homepage.layout_list',     value: 'list' },
  { label: 'admin.settings.homepage.layout_simple',   value: 'simple' },
  { label: 'admin.settings.homepage.layout_masonry',  value: 'masonry' },
  { label: 'admin.settings.homepage.layout_hero',     value: 'hero' },
  { label: 'admin.settings.homepage.layout_timeline', value: 'timeline' },
  { label: 'admin.settings.homepage.layout_ranking',  value: 'ranking' },
]

/** label/title values are i18n keys — call $t(section.label) in templates */
export const SECTION_DEFAULTS: SectionConfig[] = [
  { id: 'latest',   label: 'admin.settings.homepage.section_latest',   enabled: true,  count: 6,  title: 'admin.settings.homepage.section_latest',   layout: 'grid',    includeCategoryIds: [], excludeCategoryIds: [] },
  { id: 'hot',      label: 'admin.settings.homepage.section_hot',      enabled: true,  count: 8,  title: 'admin.settings.homepage.section_hot',      layout: 'ranking' },
  { id: 'featured', label: 'admin.settings.homepage.section_featured', enabled: false, count: 6,  title: 'admin.settings.homepage.section_featured', layout: 'hero' },
  { id: 'random',   label: 'admin.settings.homepage.section_random',   enabled: false, count: 8,  title: 'admin.settings.homepage.section_random',   layout: 'grid' },
  { id: 'timeline', label: 'admin.settings.homepage.section_timeline', enabled: false, count: 10, title: 'admin.settings.homepage.section_timeline', layout: 'timeline' },
  { id: 'masonry',  label: 'admin.settings.homepage.section_masonry',  enabled: false, count: 9,  title: 'admin.settings.homepage.section_masonry',  layout: 'masonry' },
]

export const useHomepageSections = () => {
  const optionStore = useOptionsStore()
  const { t } = useI18n()

  const getSectionConfig = (id: string): SectionConfig => {
    const saved = optionStore.getJSON<SectionConfig[]>('homepage_sections', SECTION_DEFAULTS)
    const found = saved.find(s => s.id === id)
    const def = SECTION_DEFAULTS.find(s => s.id === id)

    // Plugin section: no built-in default, return saved entry directly
    if (!def) {
      return found ?? { id, label: id, enabled: false, count: 5 }
    }

    // Built-in section: merge saved over defaults
    // label: always i18n key (used in admin UI with $t())
    // title: resolved plain text — user custom value or translated default
    const resolvedTitle = found?.title && !found.title.startsWith('admin.') ? found.title : t(def.title ?? def.label)
    return { ...def, ...found, label: def.label, title: resolvedTitle }
  }

  return { getSectionConfig, SECTION_DEFAULTS, LAYOUT_OPTIONS }
}
