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
                :icon="isAllWidgetsExpanded ? 'i-tabler-fold' : 'i-tabler-unfold'"
                @click="toggleExpandAllWidgets">
                {{ isAllWidgetsExpanded ? $t('common.collapse_all') : $t('common.expand_all') }}
              </UButton>
            </div>
          </template>
          <p class="text-sm text-muted mb-4">{{ $t('admin.settings.homepage.widget_hint') }}</p>

          <div class="space-y-2">
            <div
              v-for="(widget, index) in form.widgets"
              :key="widget.id"
              class="rounded-md border border-default overflow-hidden"
            >
              <!-- 行 -->
              <div
                class="flex items-center gap-3 px-3 py-2.5"
                :class="expanded[widget.id] ? 'bg-elevated' : 'bg-default'"
              >
                <UCheckbox v-model="widget.enabled" />
                <span
                  class="flex-1 text-sm select-none"
                  :class="widget.enabled ? 'text-highlighted' : 'text-muted'"
                >{{ $t(widget.label) }}</span>
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
                  <UButton
                    icon="i-tabler-arrow-up"
                    variant="ghost"
                    color="neutral"
                    size="xs"
                    square
                    :disabled="index === 0"
                    @click="moveUp(index)"
                  />
                  <UButton
                    icon="i-tabler-arrow-down"
                    variant="ghost"
                    color="neutral"
                    size="xs"
                    square
                    :disabled="index === form.widgets.length - 1"
                    @click="moveDown(index)"
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
                    :placeholder="$t(widget.label)"
                    size="sm"
                    class="w-40"
                  />
                </div>

                <!-- 搜索组件 -->
                <template v-if="widget.id === 'search'">
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-default">{{ $t('admin.settings.homepage.show_recent_search') }}</span>
                    <UCheckbox v-model="widget.showRecent" />
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-default">{{ $t('admin.settings.homepage.show_hot_search') }}</span>
                    <UCheckbox v-model="widget.showHot" />
                  </div>
                </template>

                <!-- 数量配置 -->
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
          </div>
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
                :icon="isAllSectionsExpanded ? 'i-tabler-fold' : 'i-tabler-unfold'"
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
import { type WidgetConfig, WIDGET_DEFAULTS } from '~/composables/useWidgetConfig'
import { type SectionConfig, type SectionActionConfig, SECTION_DEFAULTS, LAYOUT_OPTIONS as LAYOUT_OPTIONS_RAW } from '~/composables/useHomepageSections'
import { getWidgetsByContext } from '~/config/widgets'

// Widgets filtered to homepage context
const HOMEPAGE_WIDGET_DEFAULTS: WidgetConfig[] = WIDGET_DEFAULTS.filter(w =>
  getWidgetsByContext('homepage').some(def => def.id === w.id),
)

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

const form = ref({
  sidebarEnabled: false,
  widgets: HOMEPAGE_WIDGET_DEFAULTS.map(w => ({ ...w, title: resolveTitle(w.title ?? w.label) })),
  sections: SECTION_DEFAULTS.map(s => ({ ...s, title: resolveTitle(s.title ?? s.label), layout: s.layout ?? 'grid', includeCategoryIds: [...(s.includeCategoryIds ?? [])], excludeCategoryIds: [...(s.excludeCategoryIds ?? [])], action: { enabled: false, ...s.action }, loadMoreEnabled: s.loadMoreEnabled ?? false })),
})

const moveUp = (index: number) => {
  if (index === 0) return
  const arr = form.value.widgets;
  [arr[index - 1], arr[index]] = [arr[index], arr[index - 1]]
}

const moveDown = (index: number) => {
  const arr = form.value.widgets
  if (index === arr.length - 1) return;
  [arr[index], arr[index + 1]] = [arr[index + 1], arr[index]]
}

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
      const savedIds = parsed.map(w => w.id)
      const newIds = HOMEPAGE_WIDGET_DEFAULTS.map(w => w.id).filter(id => !savedIds.includes(id))
      const merged = [
        ...parsed
          .filter(saved => HOMEPAGE_WIDGET_DEFAULTS.some(d => d.id === saved.id))
          .map(saved => {
            const def = HOMEPAGE_WIDGET_DEFAULTS.find(w => w.id === saved.id)
            // always take label from def (it's an i18n key), not from saved (may be old Chinese text)
            return def ? { ...def, ...saved, label: def.label, title: resolveTitle(saved.title ?? def.title ?? def.label) } : saved
          }),
        ...HOMEPAGE_WIDGET_DEFAULTS.filter(w => newIds.includes(w.id)).map(w => ({ ...w, title: resolveTitle(w.title ?? w.label) })),
      ]
      form.value.widgets = merged
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
</script>
