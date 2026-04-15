<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.settings.reading.title')" :subtitle="$t('admin.settings.reading.subtitle')">
      <template #actions>
        <UButton color="neutral" variant="outline" :disabled="isSaving" @click="loadSettings">{{ $t('common.reset') }}</UButton>
        <UButton color="primary" icon="i-tabler-device-floppy" :loading="isSaving" :disabled="isSaving" @click="saveSettings">
          {{ $t('common.save_changes') }}
        </UButton>
      </template>
    </AdminPageHeader>
    <AdminPageContent>
      <div v-if="isLoading" class="space-y-4">
        <UCard v-for="i in 3" :key="i">
          <template #header><USkeleton class="h-5 w-40" /></template>
          <div class="space-y-2">
            <USkeleton class="h-4 w-24" />
            <USkeleton class="h-9 w-48 rounded-md" />
          </div>
        </UCard>
      </div>

      <template v-if="!isLoading">
        <!-- 文章列表 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.reading.post_list') }}</h3>
          </template>
          <UFormField :label="$t('admin.settings.reading.posts_per_page_label')">
            <UInput v-model="form.postsPerPage" type="number" :min="1" :max="50" class="w-full max-w-xs" />
            <p class="text-xs text-muted mt-1">{{ $t('admin.settings.reading.posts_per_page_hint') }}</p>
          </UFormField>
        </UCard>

        <!-- 文章设置 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.reading.post_settings') }}</h3>
          </template>
          <div class="space-y-6">
            <!-- 文章页面 -->
            <div class="space-y-4">
              <p class="text-sm font-medium text-muted uppercase tracking-wide">{{ $t('admin.settings.reading.post_page') }}</p>
              <UFormField :label="$t('admin.settings.reading.default_cover_layout')">
                <USelect
                  v-model="form.postDefaultLayout"
                  :items="[
                    { label: $t('admin.settings.reading.layout_hero'), value: 'hero' },
                    { label: $t('admin.settings.reading.layout_half'), value: 'half' },
                    { label: $t('admin.settings.reading.layout_none'), value: 'none' },
                  ]"
                  class="w-full max-w-xs" />
                <p class="text-xs text-muted mt-1">{{ $t('admin.settings.reading.layout_hint') }}</p>
              </UFormField>
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm font-medium text-highlighted">{{ $t('admin.settings.reading.default_sidebar') }}</p>
                  <p class="text-xs text-muted mt-0.5">{{ $t('admin.settings.reading.post_sidebar_desc') }}</p>
                </div>
                <USwitch v-model="form.postSidebarEnabled" />
              </div>
            </div>

            <!-- 文章侧栏小部件 -->
            <div v-if="form.postSidebarEnabled" class="space-y-3 pt-2 border-t border-default">
              <div class="flex items-center justify-between">
                <p class="text-sm font-medium text-muted uppercase tracking-wide">{{ $t('admin.settings.reading.post_widgets_title') }}</p>
                <UButton variant="ghost" color="neutral" size="xs"
                  :icon="isAllPostWidgetsExpanded ? 'i-tabler-fold' : 'i-tabler-fold-down'"
                  @click="toggleExpandAllPost">
                  {{ isAllPostWidgetsExpanded ? $t('common.collapse_all') : $t('common.expand_all') }}
                </UButton>
              </div>
              <p class="text-xs text-muted">{{ $t('admin.settings.reading.sidebar_widgets_hint') }}</p>
              <WidgetList v-model="form.postWidgets" :expanded="expanded" @toggle="toggleExpand" />
            </div>
          </div>
        </UCard>

        <!-- 页面设置 -->
        <UCard>
          <template #header>
            <h3 class="text-base font-semibold text-highlighted">{{ $t('admin.settings.reading.page_settings') }}</h3>
          </template>
          <div class="space-y-6">
            <!-- 静态页面 -->
            <div class="space-y-4">
              <p class="text-sm font-medium text-muted uppercase tracking-wide">{{ $t('admin.settings.reading.static_pages') }}</p>
              <UFormField :label="$t('admin.settings.reading.page_default_layout')">
                <USelect
                  v-model="form.pageDefaultLayout"
                  :items="[
                    { label: $t('admin.settings.reading.layout_hero'), value: 'hero' },
                    { label: $t('admin.settings.reading.layout_half'), value: 'half' },
                    { label: $t('admin.settings.reading.layout_none'), value: 'none' },
                  ]"
                  class="w-full max-w-xs" />
                <p class="text-xs text-muted mt-1">{{ $t('admin.settings.reading.page_layout_hint') }}</p>
              </UFormField>
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm font-medium text-highlighted">{{ $t('admin.settings.reading.default_sidebar') }}</p>
                  <p class="text-xs text-muted mt-0.5">{{ $t('admin.settings.reading.page_sidebar_desc') }}</p>
                </div>
                <USwitch v-model="form.pageSidebarEnabled" />
              </div>
            </div>

            <!-- 页面侧栏小部件 -->
            <div v-if="form.pageSidebarEnabled" class="space-y-3 pt-2 border-t border-default">
              <div class="flex items-center justify-between">
                <p class="text-sm font-medium text-muted uppercase tracking-wide">{{ $t('admin.settings.reading.post_widgets_title') }}</p>
                <UButton variant="ghost" color="neutral" size="xs"
                  :icon="isAllPageWidgetsExpanded ? 'i-tabler-fold' : 'i-tabler-fold-down'"
                  @click="toggleExpandAllPage">
                  {{ isAllPageWidgetsExpanded ? $t('common.collapse_all') : $t('common.expand_all') }}
                </UButton>
              </div>
              <p class="text-xs text-muted">{{ $t('admin.settings.reading.sidebar_widgets_hint') }}</p>
              <WidgetList v-model="form.pageWidgets" :expanded="expandedPage" @toggle="toggleExpandPage" />
            </div>
          </div>
        </UCard>

      </template>
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import { type WidgetConfig, WIDGET_DEFAULTS } from '~/composables/useWidgetConfig'
import type { PluginSettingField } from '~/composables/usePluginApi'
import { usePluginContributionsStore } from '~/stores/plugin-contributions'

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

const { apiFetch } = useApiFetch();
const toast = useToast();
const { t } = useI18n();

const resolveTitle = (key: string | undefined): string => {
  if (!key) return ''
  return (key.startsWith('admin.') || key.startsWith('common.')) ? t(key) : key
}
const isSaving = ref(false);
const rawLoading = ref(true);
const isLoading = useMinLoading(rawLoading);
const expanded = ref<Record<string, boolean>>({})
const expandedPage = ref<Record<string, boolean>>({})

const toggleExpand = (id: string) => { expanded.value[id] = !expanded.value[id] }
const toggleExpandPage = (id: string) => { expandedPage.value[id] = !expandedPage.value[id] }

const isAllPostWidgetsExpanded = computed(() =>
  form.value.postWidgets.length > 0 && form.value.postWidgets.every(w => expanded.value[w.id])
)
const toggleExpandAllPost = () => {
  const val = !isAllPostWidgetsExpanded.value
  form.value.postWidgets.forEach(w => { expanded.value[w.id] = val })
}

const isAllPageWidgetsExpanded = computed(() =>
  form.value.pageWidgets.length > 0 && form.value.pageWidgets.every(w => expandedPage.value[w.id])
)
const toggleExpandAllPage = () => {
  const val = !isAllPageWidgetsExpanded.value
  form.value.pageWidgets.forEach(w => { expandedPage.value[w.id] = val })
}

// Cache raw saved widgets from API for late-arriving plugin contribution merges
const savedPostWidgetsRaw = ref<WidgetConfig[]>([])
const savedPageWidgetsRaw = ref<WidgetConfig[]>([])

const form = ref({
  postsPerPage: 10,
  postDefaultLayout: 'hero' as 'hero' | 'half' | 'none',
  postSidebarEnabled: false,
  pageDefaultLayout: 'none' as 'hero' | 'half' | 'none',
  pageSidebarEnabled: false,
  postWidgets: [...WIDGET_DEFAULTS, ...pluginWidgetDefaults.value].map(w => ({ ...w, title: w.isPlugin ? w.label : resolveTitle(w.title ?? w.label) })),
  pageWidgets: [...WIDGET_DEFAULTS, ...pluginWidgetDefaults.value].map(w => ({ ...w, title: w.isPlugin ? w.label : resolveTitle(w.title ?? w.label) })),
})

const loadSettings = async () => {
  try {
    const result = await apiFetch<{ options: Record<string, string> }>("/options/autoload");
    const opts = result.options ?? {};
    if (opts.posts_per_page !== undefined) form.value.postsPerPage = parseInt(JSON.parse(opts.posts_per_page));
    if (opts.post_default_layout !== undefined) form.value.postDefaultLayout = JSON.parse(opts.post_default_layout);
    if (opts.post_sidebar_enabled !== undefined) form.value.postSidebarEnabled = JSON.parse(opts.post_sidebar_enabled);
    if (opts.page_default_layout !== undefined) form.value.pageDefaultLayout = JSON.parse(opts.page_default_layout);
    if (opts.page_sidebar_enabled !== undefined) form.value.pageSidebarEnabled = JSON.parse(opts.page_sidebar_enabled);
    const mergeWidgets = (raw: string) => {
      const parsed = JSON.parse(raw) as WidgetConfig[]
      const allDefaults = [...WIDGET_DEFAULTS, ...pluginWidgetDefaults.value]
      const savedIds = parsed.map(w => w.id)
      // Filter out uninstalled plugin widgets, keep valid saved entries
      const validSaved = parsed
        .filter(saved =>
          !saved.id.startsWith('plugin:')
            ? WIDGET_DEFAULTS.some(d => d.id === saved.id)
            : pluginWidgetDefaults.value.some(pw => pw.id === saved.id),
        )
        .map(saved => {
          const def = allDefaults.find(w => w.id === saved.id)
          return def ? { ...def, ...saved, label: def.label, title: saved.isPlugin ? (saved.title || def.label) : resolveTitle(saved.title ?? def.title ?? def.label) } : saved
        })
      // Append newly registered widgets not yet in saved config
      const newWidgets = allDefaults
        .filter(w => !savedIds.includes(w.id))
        .map(w => ({ ...w, title: w.isPlugin ? w.label : resolveTitle(w.title ?? w.label) }))
      return [...validSaved, ...newWidgets]
    }
    if (opts.blog_sidebar_widgets !== undefined) {
      savedPostWidgetsRaw.value = JSON.parse(opts.blog_sidebar_widgets) as WidgetConfig[]
      form.value.postWidgets = mergeWidgets(opts.blog_sidebar_widgets) as typeof form.value.postWidgets
    }
    if (opts.page_sidebar_widgets !== undefined) {
      savedPageWidgetsRaw.value = JSON.parse(opts.page_sidebar_widgets) as WidgetConfig[]
      form.value.pageWidgets = mergeWidgets(opts.page_sidebar_widgets) as typeof form.value.pageWidgets
    }
  } catch (e) {
    console.error(e);
    toast.add({ title: t('common.load_failed'), description: t('common.cannot_read_settings'), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    rawLoading.value = false;
  }
};

const saveSettings = async () => {
  isSaving.value = true;
  try {
    await Promise.all([
      apiFetch("/options/posts_per_page", {
        method: "PUT",
        body: { value: JSON.stringify(form.value.postsPerPage), autoload: 1 },
      }),
      apiFetch("/options/post_default_layout", {
        method: "PUT",
        body: { value: JSON.stringify(form.value.postDefaultLayout), autoload: 1 },
      }),
      apiFetch("/options/post_sidebar_enabled", {
        method: "PUT",
        body: { value: JSON.stringify(form.value.postSidebarEnabled), autoload: 1 },
      }),
      apiFetch("/options/page_default_layout", {
        method: "PUT",
        body: { value: JSON.stringify(form.value.pageDefaultLayout), autoload: 1 },
      }),
      apiFetch("/options/page_sidebar_enabled", {
        method: "PUT",
        body: { value: JSON.stringify(form.value.pageSidebarEnabled), autoload: 1 },
      }),
      apiFetch("/options/blog_sidebar_widgets", {
        method: "PUT",
        body: { value: JSON.stringify(form.value.postWidgets), autoload: 1 },
      }),
      apiFetch("/options/page_sidebar_widgets", {
        method: "PUT",
        body: { value: JSON.stringify(form.value.pageWidgets), autoload: 1 },
      }),
    ]);
    toast.add({ title: t('admin.settings.reading.saved'), description: t('admin.settings.reading.saved_desc'), color: "success", icon: "i-tabler-circle-check" });
  } catch (e) {
    console.error(e);
    toast.add({ title: t('common.save_failed'), description: t('common.settings_save_failed'), color: "error", icon: "i-tabler-alert-circle" });
  } finally {
    isSaving.value = false;
  }
};

await loadSettings();

// When plugin contributions arrive after initial load, merge new plugin widgets into form.
// Restore saved config for plugin widgets that were filtered out during loadSettings
// (because contributions hadn't loaded yet at that point).
watch(pluginWidgetDefaults, (newPluginWidgets) => {
  if (!newPluginWidgets.length) return
  const mergeInto = (widgets: WidgetConfig[], savedRaw: WidgetConfig[]) => {
    const existingIds = new Set(widgets.map(w => w.id))
    const toAdd: WidgetConfig[] = []
    for (const pw of newPluginWidgets) {
      if (existingIds.has(pw.id)) continue
      // Prefer saved config (from initial API load) over defaults
      const saved = savedRaw.find(s => s.id === pw.id)
      if (saved) {
        toAdd.push({ ...pw, ...saved, label: pw.label, title: saved.title || pw.label })
      } else {
        toAdd.push({ ...pw, title: pw.label })
      }
    }
    return toAdd.length > 0 ? [...widgets, ...toAdd] : widgets
  }
  form.value.postWidgets = mergeInto(form.value.postWidgets, savedPostWidgetsRaw.value) as typeof form.value.postWidgets
  form.value.pageWidgets = mergeInto(form.value.pageWidgets, savedPageWidgetsRaw.value) as typeof form.value.pageWidgets
})
</script>
