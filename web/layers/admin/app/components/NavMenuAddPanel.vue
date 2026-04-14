<script setup lang="ts">
import type { MenuItemType, NavMenuSlotKey } from '~/types/api/navMenu'
import { NAV_MENU_SLOT_CONFIGS } from '~/types/api/navMenu'

const props = defineProps<{
  slotKey: string
  contributionSlot?: string
}>()

const emit = defineEmits<{
  addItems: [items: Array<{ label: string; url: string; type: MenuItemType; object_id: number; local_id?: string; css_classes?: string }>]
}>()

const { t } = useI18n()
const postApi = usePostApi()
const termApi = useTermApi()
const toast = useToast()

const slotConfig = computed(() => NAV_MENU_SLOT_CONFIGS[props.slotKey as NavMenuSlotKey])
const isSocialSelected = computed(() => props.slotKey === 'social_menu')

interface SelectablePage { id: number; title: string; slug: string; selected: boolean }
interface SelectableCategory { id: number; name: string; slug: string; selected: boolean }

const contributionsStore = usePluginContributionsStore()

const pluginViewItems = computed(() =>
  props.contributionSlot ? contributionsStore.getViewItems(props.contributionSlot).value : [],
)
const pluginMenuItemsList = computed(() =>
  props.contributionSlot ? contributionsStore.getMenuItems(props.contributionSlot).value : [],
)

interface SelectablePluginItem { id: string; label: string; icon: string; selected: boolean }
const pluginContribItems = computed<SelectablePluginItem[]>(() => {
  const items: SelectablePluginItem[] = []
  for (const v of pluginViewItems.value) {
    items.push({ id: v.id, label: v.title || v.id, icon: v.icon || 'i-tabler-puzzle', selected: false })
  }
  for (const m of pluginMenuItemsList.value) {
    items.push({ id: m.command, label: m.title || m.command, icon: m.icon || 'i-tabler-puzzle', selected: false })
  }
  return items
})
const selectablePluginContribs = ref<SelectablePluginItem[]>([])
watch(pluginContribItems, (items) => {
  selectablePluginContribs.value = items.map(i => ({ ...i, selected: false }))
}, { immediate: true })

const expandedSections = ref({ builtin: true, pages: false, categories: false, custom: false, builtinActions: false, pluginPages: false, pluginContribs: false })

const builtinActionItems = ref<Array<{ id: string; label: string; icon: string; selected: boolean }>>([])

watch(() => slotConfig.value?.builtinActions, (actions) => {
  builtinActionItems.value = (actions ?? []).map(a => ({ ...a, selected: false }))
}, { immediate: true })

const builtinPages = ref([
  { label: '', url: '/', selected: false },
  { label: '', url: '/categories', selected: false },
  { label: '', url: '/tags', selected: false },
  { label: '', url: '/archive', selected: false },
  { label: '', url: '/docs', selected: false },
  { label: '', url: '/moments', selected: false },
])

const availablePages = ref<SelectablePage[]>([])
const availableCategories = ref<SelectableCategory[]>([])
const customLink = ref({ label: '', customLabel: '', url: '', icon: '' })

interface PluginPageAction {
  pluginId: string
  pluginTitle: string
  pluginIcon: string
  pagePath: string
  pageTitle: string
  selected: boolean
}
const pluginPageActions = ref<PluginPageAction[]>([])

// SOCIAL_PLATFORMS is auto-imported from ~/utils/social

watch(() => isSocialSelected.value, (val) => {
  if (val) expandedSections.value.custom = true
}, { immediate: true })

onMounted(async () => {
  builtinPages.value[0]!.label = t('admin.appearance.menus.builtin_page_home')
  builtinPages.value[1]!.label = t('admin.appearance.menus.builtin_page_categories')
  builtinPages.value[2]!.label = t('admin.appearance.menus.builtin_page_tags')
  builtinPages.value[3]!.label = t('admin.appearance.menus.builtin_page_archive')
  builtinPages.value[4]!.label = t('admin.appearance.menus.builtin_page_docs')
  builtinPages.value[5]!.label = t('admin.appearance.menus.builtin_page_moments')

  const [pagesRes, catsRes] = await Promise.all([
    postApi.getPosts({ post_type: '2', status: '2', page: 1, page_size: 100 }).catch(() => null),
    termApi.getTerms({ taxonomy: 'category' }).catch(() => [] as any[]),
  ])

  // Fetch plugin pages for action-type slots
  if (slotConfig.value?.builtinActions) {
    try {
      const resp = await fetch('/api/v1/plugins/public-client')
      if (resp.ok) {
        const json = await resp.json()
        const items: PluginPageAction[] = []
        for (const plugin of json.data?.items ?? []) {
          const pages = plugin.pages ? JSON.parse(plugin.pages) : []
          for (const pg of pages) {
            if (pg.slot === 'public' && pg.path) {
              items.push({
                pluginId: plugin.id,
                pluginTitle: plugin.title,
                pluginIcon: plugin.icon || 'i-tabler-puzzle',
                pagePath: pg.path,
                pageTitle: pg.title || plugin.title,
                selected: false,
              })
            }
          }
        }
        pluginPageActions.value = items
      }
    } catch { /* ignore */ }
  }
  availablePages.value = (pagesRes?.data ?? []).map((p: any) => ({
    id: p.id, title: p.title, slug: p.slug, selected: false,
  }))
  availableCategories.value = (Array.isArray(catsRes) ? catsRes : []).map((c: any) => ({
    id: c.term_taxonomy_id ?? c.id, name: c.name, slug: c.slug, selected: false,
  }))
})

function toggleSection(s: keyof typeof expandedSections.value) {
  expandedSections.value[s] = !expandedSections.value[s]
}

function addBuiltinPages() {
  const sel = builtinPages.value.filter(p => p.selected)
  if (!sel.length) { toast.add({ title: t('admin.appearance.menus.select_first'), color: 'warning' }); return }
  emit('addItems', sel.map(p => ({ label: p.label, url: p.url, type: 'archive' as MenuItemType, object_id: 0 })))
  builtinPages.value.forEach(p => { p.selected = false })
}

function addSelectedPages() {
  const sel = availablePages.value.filter(p => p.selected)
  if (!sel.length) { toast.add({ title: t('admin.appearance.menus.select_first'), color: 'warning' }); return }
  emit('addItems', sel.map(p => ({ label: p.title, url: `/${p.slug}`, type: 'page' as MenuItemType, object_id: p.id })))
  availablePages.value.forEach(p => { p.selected = false })
}

function addSelectedCategories() {
  const sel = availableCategories.value.filter(c => c.selected)
  if (!sel.length) { toast.add({ title: t('admin.appearance.menus.select_category_first'), color: 'warning' }); return }
  emit('addItems', sel.map(c => ({ label: c.name, url: `/category/${c.slug}`, type: 'category' as MenuItemType, object_id: c.id })))
  availableCategories.value.forEach(c => { c.selected = false })
}

function addCustomLink() {
  const isCustomPlatform = customLink.value.label === '__custom__'
  const label = isCustomPlatform ? customLink.value.customLabel : customLink.value.label
  if (!label || !customLink.value.url) { toast.add({ title: t('admin.appearance.menus.fill_link'), color: 'warning' }); return }
  emit('addItems', [{ label, url: customLink.value.url, type: 'custom' as MenuItemType, object_id: 0, css_classes: customLink.value.icon || undefined }])
  customLink.value = { label: '', customLabel: '', url: '', icon: '' }
}

function addPluginPages() {
  const sel = pluginPageActions.value.filter(a => a.selected)
  if (!sel.length) { toast.add({ title: t('admin.appearance.menus.select_first'), color: 'warning' }); return }
  emit('addItems', sel.map(a => ({
    label: a.pageTitle,
    url: a.pagePath,
    type: 'custom' as MenuItemType,
    object_id: 0,
    css_classes: a.pluginIcon,
  })))
  pluginPageActions.value.forEach(a => { a.selected = false })
}

function addBuiltinActions() {
  const sel = builtinActionItems.value.filter(a => a.selected)
  if (!sel.length) { toast.add({ title: t('admin.appearance.menus.select_first'), color: 'warning' }); return }
  emit('addItems', sel.map(a => ({ label: a.label, url: '', type: 'action' as MenuItemType, object_id: 0, local_id: a.id })))
  builtinActionItems.value.forEach(a => { a.selected = false })
}

function addPluginContribs() {
  const sel = selectablePluginContribs.value.filter(a => a.selected)
  if (!sel.length) { toast.add({ title: t('admin.appearance.menus.select_first'), color: 'warning' }); return }
  emit('addItems', sel.map(a => ({
    label: a.label,
    url: '',
    type: 'action' as MenuItemType,
    object_id: 0,
    local_id: `plugin:${a.id}`,
    css_classes: a.icon,
  })))
  selectablePluginContribs.value.forEach(a => { a.selected = false })
}

function addSeparator() {
  emit('addItems', [{ label: '—', url: '', type: 'separator' as MenuItemType, object_id: 0, local_id: `sep:${Date.now()}` }])
}
</script>

<template>
  <div class="space-y-3">
    <!-- Built-in actions (action-type slots) -->
    <UCard v-if="slotConfig?.builtinActions" :ui="{ header: 'p-2.5 sm:px-4', body: 'p-0 sm:p-0' }">
      <template #header>
        <button class="w-full flex items-center justify-between" @click="toggleSection('builtinActions')">
          <div class="flex items-center gap-2">
            <UIcon name="i-tabler-click" class="size-4 text-primary" />
            <span class="text-sm font-medium text-highlighted">{{ $t('admin.appearance.menus.builtin_actions') }}</span>
          </div>
          <UIcon name="i-tabler-chevron-down" class="size-4 text-muted transition-transform" :class="{ 'rotate-180': expandedSections.builtinActions }" />
        </button>
      </template>
      <div v-if="expandedSections.builtinActions" class="p-3 space-y-1.5">
        <div v-for="action in builtinActionItems" :key="action.id"
          class="flex items-center gap-2 p-1.5 hover:bg-elevated rounded cursor-pointer"
          @click="action.selected = !action.selected">
          <UCheckbox v-model="action.selected" @click.stop />
          <UIcon :name="action.icon" class="size-4 text-primary shrink-0" />
          <span class="text-sm text-highlighted flex-1">{{ action.label }}</span>
        </div>
        <UButton color="primary" size="sm" class="w-full mt-1" @click="addBuiltinActions">
          {{ $t('admin.appearance.menus.add_to_menu') }}
        </UButton>
      </div>
    </UCard>

    <!-- Plugin pages (action-type slots) -->
    <UCard v-if="slotConfig?.builtinActions && pluginPageActions.length" :ui="{ header: 'p-2.5 sm:px-4', body: 'p-0 sm:p-0' }">
      <template #header>
        <button class="w-full flex items-center justify-between" @click="toggleSection('pluginPages')">
          <div class="flex items-center gap-2">
            <UIcon name="i-tabler-puzzle" class="size-4 text-primary" />
            <span class="text-sm font-medium text-highlighted">插件页面</span>
          </div>
          <UIcon name="i-tabler-chevron-down" class="size-4 text-muted transition-transform" :class="{ 'rotate-180': expandedSections.pluginPages }" />
        </button>
      </template>
      <div v-if="expandedSections.pluginPages" class="p-3 space-y-1.5">
        <div v-for="action in pluginPageActions" :key="action.pluginId + action.pagePath"
          class="flex items-center gap-2 p-1.5 hover:bg-elevated rounded cursor-pointer"
          @click="action.selected = !action.selected">
          <UCheckbox v-model="action.selected" @click.stop />
          <UIcon :name="action.pluginIcon" class="size-4 text-primary shrink-0" />
          <span class="text-sm text-highlighted flex-1">{{ action.pageTitle }}</span>
        </div>
        <UButton color="primary" size="sm" class="w-full mt-1" @click="addPluginPages">
          {{ $t('admin.appearance.menus.add_to_menu') }}
        </UButton>
      </div>
    </UCard>

    <!-- Plugin components (slots with contributionSlot) -->
    <UCard v-if="selectablePluginContribs.length" :ui="{ header: 'p-2.5 sm:px-4', body: 'p-0 sm:p-0' }">
      <template #header>
        <button class="w-full flex items-center justify-between" @click="toggleSection('pluginContribs')">
          <div class="flex items-center gap-2">
            <UIcon name="i-tabler-puzzle" class="size-4 text-primary" />
            <span class="text-sm font-medium text-highlighted">{{ $t('admin.appearance.menus.plugin_items') }}</span>
          </div>
          <UIcon name="i-tabler-chevron-down" class="size-4 text-muted transition-transform" :class="{ 'rotate-180': expandedSections.pluginContribs }" />
        </button>
      </template>
      <div v-if="expandedSections.pluginContribs" class="p-3 space-y-1.5">
        <div v-for="item in selectablePluginContribs" :key="item.id"
          class="flex items-center gap-2 p-1.5 hover:bg-elevated rounded cursor-pointer"
          @click="item.selected = !item.selected">
          <UCheckbox v-model="item.selected" @click.stop />
          <UIcon :name="item.icon" class="size-4 text-primary shrink-0" />
          <span class="text-sm text-highlighted flex-1">{{ item.label }}</span>
        </div>
        <UButton color="primary" size="sm" class="w-full mt-1" @click="addPluginContribs">
          {{ $t('admin.appearance.menus.add_to_menu') }}
        </UButton>
      </div>
    </UCard>

    <!-- Built-in pages (normal menu slots) -->
    <UCard v-if="!isSocialSelected && !slotConfig?.builtinActions" :ui="{ header: 'p-2.5 sm:px-4', body: 'p-0 sm:p-0' }">
      <template #header>
        <button class="w-full flex items-center justify-between" @click="toggleSection('builtin')">
          <div class="flex items-center gap-2">
            <UIcon name="i-tabler-archive" class="size-4 text-primary" />
            <span class="text-sm font-medium text-highlighted">{{ $t('admin.appearance.menus.builtin_pages') }}</span>
          </div>
          <UIcon name="i-tabler-chevron-down" class="size-4 text-muted transition-transform" :class="{ 'rotate-180': expandedSections.builtin }" />
        </button>
      </template>
      <div v-if="expandedSections.builtin" class="p-3 space-y-1.5">
        <div v-for="page in builtinPages" :key="page.url"
          class="flex items-center gap-2 p-1.5 hover:bg-elevated rounded cursor-pointer"
          @click="page.selected = !page.selected">
          <UCheckbox v-model="page.selected" @click.stop />
          <div class="flex-1 min-w-0">
            <p class="text-sm text-highlighted">{{ page.label }}</p>
            <p class="text-xs text-muted truncate">{{ page.url }}</p>
          </div>
        </div>
        <UButton color="primary" size="sm" class="w-full mt-1" @click="addBuiltinPages">
          {{ $t('admin.appearance.menus.add_to_menu') }}
        </UButton>
      </div>
    </UCard>

    <!-- Pages (normal menu slots) -->
    <UCard v-if="!isSocialSelected && !slotConfig?.builtinActions" :ui="{ header: 'p-2.5 sm:px-4', body: 'p-0 sm:p-0' }">
      <template #header>
        <button class="w-full flex items-center justify-between" @click="toggleSection('pages')">
          <div class="flex items-center gap-2">
            <UIcon name="i-tabler-file" class="size-4 text-primary" />
            <span class="text-sm font-medium text-highlighted">{{ $t('admin.appearance.menus.pages_section') }}</span>
          </div>
          <UIcon name="i-tabler-chevron-down" class="size-4 text-muted transition-transform" :class="{ 'rotate-180': expandedSections.pages }" />
        </button>
      </template>
      <div v-if="expandedSections.pages" class="p-3 space-y-1.5">
        <div v-for="page in availablePages" :key="page.id"
          class="flex items-center gap-2 p-1.5 hover:bg-elevated rounded cursor-pointer"
          @click="page.selected = !page.selected">
          <UCheckbox v-model="page.selected" @click.stop />
          <span class="text-sm text-highlighted flex-1 truncate">{{ page.title }}</span>
        </div>
        <div v-if="!availablePages.length" class="text-xs text-muted text-center py-2">
          {{ $t('admin.appearance.menus.no_pages') }}
        </div>
        <UButton color="primary" size="sm" class="w-full mt-1" @click="addSelectedPages">
          {{ $t('admin.appearance.menus.add_to_menu') }}
        </UButton>
      </div>
    </UCard>

    <!-- Categories (normal menu slots) -->
    <UCard v-if="!isSocialSelected && !slotConfig?.builtinActions" :ui="{ header: 'p-2.5 sm:px-4', body: 'p-0 sm:p-0' }">
      <template #header>
        <button class="w-full flex items-center justify-between" @click="toggleSection('categories')">
          <div class="flex items-center gap-2">
            <UIcon name="i-tabler-folder" class="size-4 text-primary" />
            <span class="text-sm font-medium text-highlighted">{{ $t('admin.appearance.menus.categories_section') }}</span>
          </div>
          <UIcon name="i-tabler-chevron-down" class="size-4 text-muted transition-transform" :class="{ 'rotate-180': expandedSections.categories }" />
        </button>
      </template>
      <div v-if="expandedSections.categories" class="p-3 space-y-1.5">
        <div v-for="cat in availableCategories" :key="cat.id"
          class="flex items-center gap-2 p-1.5 hover:bg-elevated rounded cursor-pointer"
          @click="cat.selected = !cat.selected">
          <UCheckbox v-model="cat.selected" @click.stop />
          <span class="text-sm text-highlighted flex-1 truncate">{{ cat.name }}</span>
        </div>
        <div v-if="!availableCategories.length" class="text-xs text-muted text-center py-2">
          {{ $t('admin.appearance.menus.no_categories') }}
        </div>
        <UButton color="primary" size="sm" class="w-full mt-1" @click="addSelectedCategories">
          {{ $t('admin.appearance.menus.add_to_menu') }}
        </UButton>
      </div>
    </UCard>

    <!-- Custom Link -->
    <UCard v-if="(slotConfig?.showCustomLink !== false && !slotConfig?.builtinActions) || slotConfig?.showCustomLink === true" :ui="{ header: 'p-2.5 sm:px-4', body: 'p-0 sm:p-0' }">
      <template #header>
        <button class="w-full flex items-center justify-between" @click="toggleSection('custom')">
          <div class="flex items-center gap-2">
            <UIcon name="i-tabler-link" class="size-4 text-primary" />
            <span class="text-sm font-medium text-highlighted">{{ $t('admin.appearance.menus.custom_link') }}</span>
          </div>
          <UIcon name="i-tabler-chevron-down" class="size-4 text-muted transition-transform" :class="{ 'rotate-180': expandedSections.custom }" />
        </button>
      </template>
      <div v-if="expandedSections.custom" class="p-3 space-y-2">
        <UFormField :label="$t('admin.appearance.menus.link_label')">
          <template v-if="isSocialSelected">
            <USelect v-model="customLink.label" :items="SOCIAL_PLATFORMS" value-key="value" label-key="label"
              :placeholder="$t('admin.appearance.menus.select_platform')" class="w-full" size="sm" />
            <UInput v-if="customLink.label === '__custom__'" v-model="customLink.customLabel"
              :placeholder="$t('admin.appearance.menus.custom_name_placeholder')" class="w-full mt-2" size="sm" />
          </template>
          <UInput v-else v-model="customLink.label" :placeholder="$t('admin.appearance.menus.link_label')" class="w-full" size="sm" />
        </UFormField>
        <UFormField :label="$t('admin.appearance.menus.link_url')">
          <UInput v-model="customLink.url"
            :placeholder="isSocialSelected ? 'https://github.com/yourname' : 'https://example.com'"
            class="w-full" size="sm" />
        </UFormField>
        <UFormField v-if="slotConfig?.customLinkHasIcon" label="图标">
          <UInput v-model="customLink.icon" placeholder="i-tabler-link" class="w-full" size="sm" />
          <p class="text-xs text-muted mt-1">Iconify 图标名称，如 i-tabler-link</p>
        </UFormField>
        <UButton color="primary" size="sm" class="w-full" @click="addCustomLink">
          {{ $t('admin.appearance.menus.add_to_menu') }}
        </UButton>
      </div>
    </UCard>

    <!-- Add separator (all slots) -->
    <UButton
      color="neutral" variant="outline" size="sm"
      leading-icon="i-tabler-separator"
      class="w-full"
      @click="addSeparator">
      {{ $t('admin.appearance.menus.add_separator') }}
    </UButton>
  </div>
</template>
