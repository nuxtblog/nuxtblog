<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.settings.homepage.title')" :subtitle="$t('admin.settings.homepage.subtitle')">
      <template #actions>
        <UButton color="neutral" variant="outline" :disabled="isSaving" @click="loadSettings">{{ $t('common.reset') }}</UButton>
        <UButton color="primary" icon="i-tabler-device-floppy" :loading="isSaving" :disabled="isSaving" @click="saveSettings">
          {{ $t('common.save_changes') }}
        </UButton>
      </template>
    </AdminPageHeader>
    <AdminPageContent>
      <!-- 加载骨架 -->
      <div v-if="isLoading" class="space-y-6">
        <UCard v-for="i in 2" :key="i">
          <template #header><USkeleton class="h-5 w-32" /></template>
          <div class="space-y-3">
            <div v-for="j in 3" :key="j" class="h-10 bg-muted rounded animate-pulse" />
          </div>
        </UCard>
      </div>

      <div v-else class="space-y-6">
        <!-- 侧栏开关 -->
        <UCard>
          <template #header>
            <h3 class="font-semibold text-highlighted">{{ $t('admin.settings.homepage.sidebar_settings') }}</h3>
          </template>
          <div class="flex items-center justify-between">
            <div>
              <h4 class="text-sm font-medium text-highlighted mb-1">{{ $t('admin.settings.homepage.enable_sidebar') }}</h4>
              <p class="text-xs text-muted">{{ $t('admin.settings.homepage.enable_sidebar_desc') }}</p>
            </div>
            <UCheckbox v-model="form.sidebarEnabled" />
          </div>
        </UCard>

        <!-- 小部件管理 -->
        <UCard v-if="form.sidebarEnabled">
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="font-semibold text-highlighted">{{ $t('admin.settings.homepage.widget_management') }}</h3>
              <UButton
                variant="ghost"
                color="neutral"
                size="xs"
                :icon="isAllWidgetsExpanded ? 'i-tabler-fold' : 'i-tabler-fold-down'"
                @click="toggleExpandAllWidgets">
                {{ isAllWidgetsExpanded ? $t('common.collapse_all') : $t('common.expand_all') }}
              </UButton>
            </div>
          </template>
          <p class="text-sm text-muted mb-4">{{ $t('admin.settings.homepage.widget_hint') }}</p>

          <VueDraggable
            tag="div"
            class="space-y-2"
            v-model="form.widgets"
            :animation="200"
            handle=".widget-drag-handle"
          >
            <div
              v-for="widget in form.widgets"
              :key="widget.id"
              class="rounded-md border border-default overflow-hidden"
            >
              <!-- 行 -->
              <div
                class="flex items-center gap-3 px-3 py-2.5"
                :class="expanded[widget.id] ? 'bg-elevated' : 'bg-default'"
              >
                <UIcon
                  name="i-tabler-grip-vertical"
                  class="widget-drag-handle size-4 shrink-0 cursor-grab text-muted hover:text-primary active:cursor-grabbing"
                />
                <UCheckbox v-model="widget.enabled" />
                <span
                  class="flex-1 text-sm select-none"
                  :class="widget.enabled ? 'text-highlighted' : 'text-muted'"
                >
                  {{ widget.isPlugin ? widget.label : $t(widget.label) }}
                  <UIcon v-if="widget.isPlugin" name="i-tabler-puzzle" class="size-3.5 text-muted ml-1 inline-block align-text-bottom" />
                </span>
                <div class="flex items-center gap-0.5">
                  <!-- 展开配置按钮（author 无配置项） -->
                  <UButton
                    :icon="expanded[widget.id] ? 'i-tabler-chevron-up' : 'i-tabler-settings'"
                    variant="ghost"
                    :color="expanded[widget.id] ? 'primary' : 'neutral'"
                    size="xs"
                    square
                    @click="toggleExpand(widget.id)"
                  />
                </div>
              </div>

              <!-- 展开的配置面板 -->
              <div
                v-if="expanded[widget.id]"
                class="px-4 py-3 bg-muted/40 border-t border-default space-y-3"
              >
                <!-- 标题配置（所有 widget 通用） -->
                <div class="flex items-center gap-3">
                  <span class="text-sm text-default flex-1">{{ $t('admin.settings.homepage.widget_title_label') }}</span>
                  <UInput
                    v-model="widget.title"
                    :placeholder="widget.isPlugin ? widget.label : $t(widget.label)"
                    size="sm"
                    class="w-40"
                  />
                </div>

                <!-- 插件 widget: 仅自定义设置 -->
                <template v-if="widget.isPlugin">
                  <PluginSettingFields
                    v-if="getWidgetSettings(widget.id)?.length && widget.pluginSettings"
                    :schema="getWidgetSettings(widget.id)!"
                    :model-value="widget.pluginSettings"
                    @update:model-value="widget.pluginSettings = $event"
                  />
                </template>

                <!-- 搜索组件 -->
                <template v-else-if="widget.id === 'search'">
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-default">{{ $t('admin.settings.homepage.show_recent_search') }}</span>
                    <UCheckbox v-model="widget.showRecent" />
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-default">{{ $t('admin.settings.homepage.show_hot_search') }}</span>
                    <UCheckbox v-model="widget.showHot" />
                  </div>
                </template>

                <!-- 内建 widget 数量配置 -->
                <template v-else>
                  <div class="flex items-center gap-3">
                    <span class="text-sm text-default flex-1">{{ $t('admin.settings.homepage.max_count_label') }}</span>
                    <UInput
                      v-model.number="widget.maxCount"
                      type="number"
                      :min="1"
                      :max="20"
                      size="sm"
                      class="w-20"
                    />
                    <span class="text-xs text-muted">{{ $t('admin.settings.homepage.count_unit') }}</span>
                  </div>
                </template>
              </div>
            </div>
          </VueDraggable>
        </UCard>

        <!-- 内容区块 -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <h3 class="font-semibold text-highlighted">{{ $t('admin.settings.homepage.content_sections') }}</h3>
              <UButton
                variant="ghost"
                color="neutral"
                size="xs"
                :icon="isAllSectionsExpanded ? 'i-tabler-fold' : 'i-tabler-fold-down'"
                @click="toggleExpandAllSections">
                {{ isAllSectionsExpanded ? $t('common.collapse_all') : $t('common.expand_all') }}
              </UButton>
            </div>
          </template>
          <p class="text-sm text-muted mb-4">{{ $t('admin.settings.homepage.sections_hint') }}</p>

          <div class="space-y-2">
            <div
              v-for="(section, index) in form.sections"
              :key="section.id"
              class="rounded-md border border-default overflow-hidden"
            >
              <!-- 行 -->
              <div
                class="flex items-center gap-3 px-3 py-2.5"
                :class="expandedSection[section.id] ? 'bg-elevated' : 'bg-default'"
              >
                <UCheckbox v-model="section.enabled" />
                <span
                  class="flex-1 text-sm select-none"
                  :class="section.enabled ? 'text-highlighted' : 'text-muted'"
                >{{ $t(section.label) }}</span>
                <UButton
                  :icon="expandedSection[section.id] ? 'i-tabler-chevron-up' : 'i-tabler-settings'"
                  variant="ghost"
                  :color="expandedSection[section.id] ? 'primary' : 'neutral'"
                  size="xs"
                  square
                  @click="toggleSectionExpand(section.id)"
                />
              </div>

              <!-- 展开的配置面板 -->
              <div
                v-if="expandedSection[section.id]"
                class="px-4 py-3 bg-muted/40 border-t border-default space-y-4"
              >
                <!-- 标题 -->
                <div class="flex items-center gap-3">
                  <span class="text-sm text-default flex-1">{{ $t('admin.settings.homepage.section_title_label') }}</span>
                  <UInput
                    v-model="section.title"
                    :placeholder="$t(section.label)"
                    size="sm"
                    class="w-40"
                  />
                </div>

                <!-- 布局 -->
                <div class="flex items-center gap-3">
                  <span class="text-sm text-default flex-1">{{ $t('admin.settings.homepage.section_layout_label') }}</span>
                  <USelect
                    v-model="section.layout"
                    :items="LAYOUT_OPTIONS"
                    value-key="value"
                    label-key="label"
                    size="sm"
                    class="w-40"
                  />
                </div>

                <!-- 数量 -->
                <div class="flex items-center gap-3">
                  <span class="text-sm text-default flex-1">{{ $t('admin.settings.homepage.section_count_label') }}</span>
                  <UInput
                    v-model.number="section.count"
                    type="number"
                    :min="1"
                    :max="50"
                    size="sm"
                    class="w-20"
                  />
                  <span class="text-xs text-muted">{{ $t('admin.settings.homepage.count_unit2') }}</span>
                </div>

                <!-- 分类过滤（仅最新文章） -->
                <template v-if="section.id === 'latest'">
                  <div class="space-y-2">
                    <p class="text-xs font-medium text-muted uppercase tracking-wide">{{ $t('admin.settings.homepage.include_categories') }}</p>
                    <div v-if="categories.length === 0" class="text-xs text-muted">{{ $t('admin.settings.homepage.no_categories') }}</div>
                    <div v-else class="grid grid-cols-2 gap-1.5">
                      <label
                        v-for="cat in categories"
                        :key="cat.term_taxonomy_id"
                        class="flex items-center gap-2 cursor-pointer"
                      >
                        <UCheckbox
                          :model-value="section.includeCategoryIds?.includes(cat.term_taxonomy_id) ?? false"
                          @update:model-value="toggleCategoryId(section, 'includeCategoryIds', cat.term_taxonomy_id, $event)"
                        />
                        <span class="text-sm text-default">{{ cat.name }}</span>
                      </label>
                    </div>
                    <p class="text-xs text-muted">{{ $t('admin.settings.homepage.all_categories_hint') }}</p>
                  </div>
                  <div class="space-y-2">
                    <p class="text-xs font-medium text-muted uppercase tracking-wide">{{ $t('admin.settings.homepage.exclude_categories') }}</p>
                    <div v-if="categories.length === 0" class="text-xs text-muted">{{ $t('admin.settings.homepage.no_categories') }}</div>
                    <div v-else class="grid grid-cols-2 gap-1.5">
                      <label
                        v-for="cat in categories"
                        :key="cat.term_taxonomy_id"
                        class="flex items-center gap-2 cursor-pointer"
                      >
                        <UCheckbox
                          :model-value="section.excludeCategoryIds?.includes(cat.term_taxonomy_id) ?? false"
                          @update:model-value="toggleCategoryId(section, 'excludeCategoryIds', cat.term_taxonomy_id, $event)"
                        />
                        <span class="text-sm text-default">{{ cat.name }}</span>
                      </label>
                    </div>
                  </div>
                </template>

                <!-- 头部操作按钮 -->
                <div class="space-y-3 pt-3 border-t border-default">
                  <p class="text-xs font-medium text-muted uppercase tracking-wide">{{ $t('admin.settings.homepage.section_action_label') }}</p>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-default">{{ $t('admin.settings.homepage.section_action_enable') }}</span>
                    <UCheckbox v-model="ensureAction(section).enabled" />
                  </div>
                  <template v-if="section.action?.enabled">
                    <!-- 按钮文字 -->
                    <div class="flex items-center gap-3">
                      <span class="text-sm text-default flex-1">{{ $t('admin.settings.homepage.section_action_btn_label') }}</span>
                      <UInput
                        v-model="ensureAction(section).label"
                        :placeholder="sectionDefaultActionLabel(section.id)"
                        size="sm"
                        class="w-40" />
                    </div>
                    <!-- 跳转分类（仅最新文章） -->
                    <div v-if="section.id === 'latest'" class="flex items-center gap-3">
                      <span class="text-sm text-default flex-1">{{ $t('admin.settings.homepage.section_action_category') }}</span>
                      <AdminSearchableSelect
                        v-model="ensureAction(section).categorySlug"
                        :items="[{ label: $t('admin.posts.all_categories'), value: '' }, ...categories.map(c => ({ label: c.name, value: c.slug }))]"
                        :placeholder="$t('admin.posts.all_categories')"
                        :search-placeholder="$t('admin.posts.search_categories')"
                        trigger-class="w-40 justify-between" />
                    </div>
                    <!-- 自定义链接（非随机刷新类型） -->
                    <div v-if="section.id !== 'random'" class="flex items-center gap-3">
                      <span class="text-sm text-default flex-1">{{ $t('admin.settings.homepage.section_action_href') }}</span>
                      <UInput
                        v-model="ensureAction(section).href"
                        :placeholder="section.id === 'latest' ? (section.action?.categorySlug ? `/category/${section.action.categorySlug}` : '/posts') : sectionDefaultActionHref(section.id)"
                        size="sm"
                        class="w-40" />
                    </div>
                  </template>
                </div>

                <!-- 瀑布流：底部加载更多 -->
                <div v-if="section.id === 'masonry'" class="space-y-1 pt-3 border-t border-default">
                  <div class="flex items-center justify-between">
                    <div>
                      <span class="text-sm text-default">{{ $t('admin.settings.homepage.section_load_more_enable') }}</span>
                      <p class="text-xs text-muted mt-0.5">{{ $t('admin.settings.homepage.section_load_more_enable_desc') }}</p>
                    </div>
                    <UCheckbox v-model="section.loadMoreEnabled" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </UCard>

      </div>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import { type WidgetConfig, WIDGET_DEFAULTS } from '~/composables/useWidgetConfig'
import type { PluginSettingField } from '~/composables/usePluginApi'
import { type SectionConfig, type SectionActionConfig, SECTION_DEFAULTS, LAYOUT_OPTIONS as LAYOUT_OPTIONS_RAW } from '~/composables/useHomepageSections'
import { getWidgetsByContext } from '~/config/widgets'
import { usePluginContributionsStore } from '~/stores/plugin-contributions'

// Widgets filtered to homepage context
const HOMEPAGE_WIDGET_DEFAULTS: WidgetConfig[] = WIDGET_DEFAULTS.filter(w =>
  getWidgetsByContext('homepage').some(def => def.id === w.id),
)

// Plugin sidebar widgets from contributions
const contributionsStore = usePluginContributionsStore()
const pluginWidgetViews = contributionsStore.getViewItems('public:sidebar-widget')
const pluginWidgetDefaults = computed((): WidgetConfig[] =>
  pluginWidgetViews.value
    .filter((v: { component?: string; module?: string }) => v.component && v.module)
    .map((v: { pluginId: string; id: string; title: string; component?: string; module?: string; settings?: PluginSettingField[] }) => {
      const pluginSettings: Record<string, unknown> = {}
      if (v.settings) {
        for (const s of v.settings) {
          if (s.default !== undefined) pluginSettings[s.key] = s.default
        }
      }
      return {
        id: `plugin:${v.pluginId}:${v.id}`,
        label: v.title || v.id,
        enabled: false,
        isPlugin: true as const,
        pluginId: v.pluginId,
        component: v.component,
        module: v.module,
        maxCount: 5,
        ...(Object.keys(pluginSettings).length > 0 ? { pluginSettings } : {}),
      }
    }),
)

/** Look up widget settings schema from the plugin contributions store. */
function getWidgetSettings(widgetId: string): PluginSettingField[] | undefined {
  const parts = widgetId.split(':')
  if (parts.length < 3 || parts[0] !== 'plugin') return undefined
  const pluginId = parts[1]
  const viewId = parts.slice(2).join(':')
  const view = pluginWidgetViews.value.find(v => v.pluginId === pluginId && v.id === viewId)
  return view?.settings?.length ? view.settings : undefined
}

/** Initialize pluginSettings on a widget if missing. Must be called OUTSIDE render. */
function initPluginSettings(widget: WidgetConfig) {
  if (widget.pluginSettings) return
  const schema = getWidgetSettings(widget.id)
  if (!schema) return
  const defaults: Record<string, unknown> = {}
  for (const field of schema) {
    if (field.default !== undefined) defaults[field.key] = field.default
  }
  widget.pluginSettings = defaults
}

const { apiFetch } = useApiFetch()
const termApi = useTermApi()
const toast = useToast()
const { t } = useI18n()

// Resolve i18n key to translated string (safe to call on already-translated strings)
const resolveTitle = (key: string | undefined): string => {
  if (!key) return ''
  return (key.startsWith('admin.') || key.startsWith('common.')) ? t(key) : key
}
const LAYOUT_OPTIONS = computed(() => LAYOUT_OPTIONS_RAW.map(o => ({ ...o, label: t(o.label) })))
const isSaving = ref(false)
const isLoading = ref(false)
const expanded = ref<Record<string, boolean>>({})
const expandedSection = ref<Record<string, boolean>>({})

const toggleExpand = (id: string) => {
  // Init pluginSettings before expanding (avoids render-phase mutation)
  const widget = form.value.widgets.find(w => w.id === id)
  if (widget?.isPlugin) initPluginSettings(widget)
  expanded.value[id] = !expanded.value[id]
}
const toggleSectionExpand = (id: string) => {
  expandedSection.value[id] = !expandedSection.value[id]
}

const isAllWidgetsExpanded = computed(() =>
  form.value.widgets.length > 0 && form.value.widgets.every(w => expanded.value[w.id])
)
const toggleExpandAllWidgets = () => {
  const val = !isAllWidgetsExpanded.value
  if (val) form.value.widgets.forEach(w => { if (w.isPlugin) initPluginSettings(w) })
  form.value.widgets.forEach(w => { expanded.value[w.id] = val })
}

const isAllSectionsExpanded = computed(() =>
  form.value.sections.length > 0 && form.value.sections.every(s => expandedSection.value[s.id])
)
const toggleExpandAllSections = () => {
  const val = !isAllSectionsExpanded.value
  form.value.sections.forEach(s => { expandedSection.value[s.id] = val })
}

// 分类列表（用于内容区块的分类过滤 + 头部操作按钮）
const categories = ref<Array<{ term_taxonomy_id: number; name: string; slug: string }>>([])

// Default action button labels/hrefs per section
const sectionDefaultActionLabel = (id: string) => ({
  latest: t('common.view_more'), hot: t('common.view_more'),
  featured: t('common.view_more'), random: t('common.refresh'),
  timeline: t('common.view_more'), masonry: t('common.view_more'),
}[id] ?? t('common.view_more'))

const sectionDefaultActionHref = (id: string) => ({
  hot: '/posts?sort=views', featured: '/posts?featured=1',
  timeline: '/archive', masonry: '/posts',
}[id] ?? '/posts')

const ensureAction = (section: SectionConfig): SectionActionConfig => {
  if (!section.action) section.action = { enabled: false }
  return section.action
}

const toggleCategoryId = (
  section: SectionConfig,
  field: 'includeCategoryIds' | 'excludeCategoryIds',
  id: number,
  checked: boolean,
) => {
  const arr = section[field] ?? []
  if (checked) {
    if (!arr.includes(id)) section[field] = [...arr, id]
  } else {
    section[field] = arr.filter(x => x !== id)
  }
}

// Cache raw saved widgets from API for late-arriving plugin contribution merges
const savedWidgetsRaw = ref<WidgetConfig[]>([])

const form = ref({
  sidebarEnabled: false,
  widgets: [...HOMEPAGE_WIDGET_DEFAULTS, ...pluginWidgetDefaults.value].map(w => ({ ...w, title: w.isPlugin ? w.label : resolveTitle(w.title ?? w.label) })),
  sections: SECTION_DEFAULTS.map(s => ({ ...s, title: resolveTitle(s.title ?? s.label), layout: s.layout ?? 'grid', includeCategoryIds: [...(s.includeCategoryIds ?? [])], excludeCategoryIds: [...(s.excludeCategoryIds ?? [])], action: { enabled: false, ...s.action }, loadMoreEnabled: s.loadMoreEnabled ?? false })),
})

const loadSettings = async () => {
  isLoading.value = true
  try {
    const [result, cats] = await Promise.all([
      apiFetch<{ options: Record<string, string> }>('/options/autoload'),
      termApi.getTerms({ taxonomy: 'category' }).catch(() => []),
    ])
    const opts = result.options ?? {}
    categories.value = cats.map(c => ({ term_taxonomy_id: c.term_taxonomy_id, name: c.name, slug: c.slug }))

    if (opts.homepage_sidebar_enabled !== undefined) {
      form.value.sidebarEnabled = JSON.parse(opts.homepage_sidebar_enabled) as boolean
    }
    if (opts.homepage_sidebar_widgets !== undefined) {
      const parsed = JSON.parse(opts.homepage_sidebar_widgets) as WidgetConfig[]
      savedWidgetsRaw.value = parsed
      const allDefaults = [...HOMEPAGE_WIDGET_DEFAULTS, ...pluginWidgetDefaults.value]
      const savedIds = parsed.map(w => w.id)
      // Filter out uninstalled plugin widgets, keep valid saved entries
      const validSaved = parsed
        .filter(saved =>
          !saved.id.startsWith('plugin:')
            ? HOMEPAGE_WIDGET_DEFAULTS.some(d => d.id === saved.id)
            : pluginWidgetDefaults.value.some(pw => pw.id === saved.id),
        )
        .map(saved => {
          const def = allDefaults.find(w => w.id === saved.id)
          // For built-in: always take label from def (i18n key). For plugin: take label from def (latest title).
          return def ? { ...def, ...saved, label: def.label, title: saved.isPlugin ? (saved.title || def.label) : resolveTitle(saved.title ?? def.title ?? def.label) } : saved
        })
      // Append newly registered widgets (built-in or plugin) not yet in saved config
      const newWidgets = allDefaults
        .filter(w => !savedIds.includes(w.id))
        .map(w => ({ ...w, title: w.isPlugin ? w.label : resolveTitle(w.title ?? w.label) }))
      form.value.widgets = [...validSaved, ...newWidgets] as typeof form.value.widgets
    } else {
      // No saved config — use all defaults including plugin widgets
      const allDefaults = [...HOMEPAGE_WIDGET_DEFAULTS, ...pluginWidgetDefaults.value]
      form.value.widgets = allDefaults.map(w => ({ ...w, title: w.isPlugin ? w.label : resolveTitle(w.title ?? w.label) })) as typeof form.value.widgets
    }
    if (opts.homepage_sections !== undefined) {
      const parsed = JSON.parse(opts.homepage_sections) as SectionConfig[]
      form.value.sections = SECTION_DEFAULTS.map(def => {
        const saved = parsed.find(s => s.id === def.id)
        return saved
          ? { ...def, ...saved, label: def.label, title: resolveTitle(saved.title ?? def.title ?? def.label), layout: saved.layout ?? def.layout ?? 'grid', includeCategoryIds: saved.includeCategoryIds ?? [], excludeCategoryIds: saved.excludeCategoryIds ?? [], action: { enabled: false, ...def.action, ...saved.action }, loadMoreEnabled: saved.loadMoreEnabled ?? def.loadMoreEnabled ?? false }
          : { ...def, title: resolveTitle(def.title ?? def.label), layout: def.layout ?? 'grid', action: { enabled: false, ...def.action }, loadMoreEnabled: def.loadMoreEnabled ?? false }
      })
    }
  } catch (e) {
    console.error(e)
    toast.add({ title: t('common.load_failed'), description: t('common.cannot_read_settings'), color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    isLoading.value = false
  }
}

const saveSettings = async () => {
  isSaving.value = true
  try {
    await Promise.all([
      apiFetch('/options/homepage_sidebar_enabled', {
        method: 'PUT',
        body: { value: JSON.stringify(form.value.sidebarEnabled), autoload: 1 },
      }),
      apiFetch('/options/homepage_sidebar_widgets', {
        method: 'PUT',
        body: { value: JSON.stringify(form.value.widgets), autoload: 1 },
      }),
      apiFetch('/options/homepage_sections', {
        method: 'PUT',
        body: { value: JSON.stringify(form.value.sections), autoload: 1 },
      }),
    ])
    toast.add({ title: t('admin.settings.homepage.saved'), description: t('admin.settings.homepage.saved_desc'), color: 'success', icon: 'i-tabler-circle-check' })
  } catch (e) {
    console.error(e)
    toast.add({ title: t('common.save_failed'), description: t('common.settings_save_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    isSaving.value = false
  }
}

await loadSettings()

// When plugin contributions arrive after initial load, merge new plugin widgets into form.
// Restore saved config for plugin widgets that were filtered out during loadSettings.
watch(pluginWidgetDefaults, (newPluginWidgets) => {
  if (!newPluginWidgets.length) return
  const existingIds = new Set(form.value.widgets.map(w => w.id))
  const toAdd: WidgetConfig[] = []
  for (const pw of newPluginWidgets) {
    if (existingIds.has(pw.id)) continue
    const saved = savedWidgetsRaw.value.find(s => s.id === pw.id)
    if (saved) {
      toAdd.push({ ...pw, ...saved, label: pw.label, title: saved.title || pw.label } as WidgetConfig)
    } else {
      toAdd.push({ ...pw, title: pw.label } as WidgetConfig)
    }
  }
  if (toAdd.length > 0) {
    form.value.widgets = [...form.value.widgets, ...toAdd] as typeof form.value.widgets
  }
})
</script>
