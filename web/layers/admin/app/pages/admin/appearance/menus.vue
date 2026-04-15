<template>
  <AdminPageContainer>
    <AdminPageHeader
      :title="$t('admin.appearance.menus.title')"
      :subtitle="$t('admin.appearance.menus.subtitle')" />
    <AdminPageContent>
      <!-- Loading skeleton -->
      <template v-if="isLoading">
        <div class="grid grid-cols-1 lg:grid-cols-6 gap-4">
          <div class="space-y-2">
            <USkeleton v-for="n in 5" :key="n" class="h-12 w-full rounded-md" />
          </div>
          <div class="lg:col-span-2 space-y-3">
            <USkeleton class="h-8 w-48" />
            <USkeleton class="h-48 w-full rounded-md" />
          </div>
          <div class="lg:col-span-3 space-y-4">
            <USkeleton class="h-64 w-full rounded-md" />
          </div>
        </div>
      </template>

      <template v-else>
        <div class="grid grid-cols-1 lg:grid-cols-6 gap-4 md:gap-6">
          <!-- Left: Slot / menu list -->
          <div class="lg:col-span-1 space-y-2">
            <p class="text-xs font-semibold text-muted uppercase tracking-wider px-2 mb-1">
              {{ $t('admin.appearance.menus.theme_positions') }}
            </p>
            <button
              v-for="slot in SLOT_LIST"
              :key="slot.key"
              class="w-full flex items-center gap-3 px-3 py-2.5 rounded-md text-left transition-colors"
              :class="selectedKey === slot.key ? 'bg-primary/10 text-primary' : 'hover:bg-elevated text-highlighted'"
              @click="selectSlot(slot.key)">
              <UIcon name="i-tabler-lock" class="size-3.5 shrink-0 opacity-60" />
              <span class="text-sm font-medium flex-1">{{ slot.label }}</span>
              <UBadge v-if="slotItems(slot.key).length" :label="String(slotItems(slot.key).length)" color="neutral" variant="soft" size="xs" />
            </button>

            <USeparator class="my-3" />

            <p class="text-xs font-semibold text-muted uppercase tracking-wider px-2 mb-1">
              {{ $t('admin.appearance.menus.custom_menus') }}
            </p>
            <button
              v-for="menu in customMenus"
              :key="menu.id"
              class="w-full flex items-center gap-3 px-3 py-2.5 rounded-md text-left transition-colors"
              :class="selectedKey === menu.id ? 'bg-primary/10 text-primary' : 'hover:bg-elevated text-highlighted'"
              @click="selectCustom(menu.id)">
              <UIcon name="i-tabler-menu-2" class="size-3.5 shrink-0 opacity-60" />
              <span class="text-sm font-medium flex-1 truncate">{{ menu.name }}</span>
              <UBadge v-if="menu.items.length" :label="String(menu.items.length)" color="neutral" variant="soft" size="xs" />
            </button>

            <UButton
              color="primary" variant="soft" leading-icon="i-tabler-plus"
              class="w-full" size="sm"
              @click="showCreateModal = true">
              {{ $t('admin.appearance.menus.new_menu') }}
            </UButton>
          </div>

          <!-- Middle: Add items panel -->
          <div class="lg:col-span-2 space-y-3">
            <template v-if="selectedKey">
              <!-- Header row -->
              <div class="flex items-center justify-between">
                <div class="space-y-0.5">
                  <div class="flex items-center gap-3">
                    <h2 class="text-lg font-semibold text-highlighted">{{ selectedLabel }}</h2>
                    <UBadge
                      v-if="isBuiltinSelected"
                      :label="$t('admin.appearance.menus.builtin_badge')"
                      color="primary" variant="soft" size="xs" />
                  </div>
                  <p v-if="selectedHint" class="text-xs text-muted">{{ selectedHint }}</p>
                </div>
                <UButton
                  v-if="!isBuiltinSelected"
                  color="error" variant="ghost" size="sm"
                  leading-icon="i-tabler-trash"
                  :loading="isDeleting"
                  @click="showDeleteModal = true">
                  {{ $t('common.delete') }}
                </UButton>
              </div>
              <NavMenuAddPanel :slot-key="selectedKey" :contribution-slot="currentContribSlot" @add-items="onAddItems" />
            </template>
            <div v-else class="flex items-center justify-center h-64 text-muted text-sm">
              {{ $t('admin.appearance.menus.select_to_edit') }}
            </div>
          </div>

          <!-- Right: Menu structure + preview + save -->
          <div class="lg:col-span-3 space-y-4">
            <template v-if="selectedKey">
              <NavMenuStructure v-model:items="menuItems" :slot-key="selectedKey" />
              <NavMenuPreview :items="menuItems" :slot-key="selectedKey" />

              <!-- Save / Reset -->
              <div class="flex justify-end gap-3">
                <UButton color="neutral" variant="outline" @click="resetCurrentMenu">
                  {{ $t('common.reset') }}
                </UButton>
                <UButton
                  color="primary" :loading="isSaving"
                  leading-icon="i-tabler-device-floppy"
                  @click="saveCurrentMenu">
                  {{ $t('admin.appearance.menus.save_menu') }}
                </UButton>
              </div>
            </template>
          </div>
        </div>
      </template>
    </AdminPageContent>

    <!-- Create custom menu modal -->
    <UModal v-model:open="showCreateModal">
      <template #content>
        <div class="p-6 space-y-4">
          <h3 class="text-lg font-semibold text-highlighted mb-4">
            {{ $t('admin.appearance.menus.create_modal_title') }}
          </h3>
          <UFormField :label="$t('admin.appearance.menus.menu_name_label')" required>
            <UInput v-model="newMenu.name" :placeholder="$t('admin.appearance.menus.menu_name_placeholder')" class="w-full" />
          </UFormField>
          <UFormField :label="$t('admin.appearance.menus.menu_desc_label')">
            <UTextarea v-model="newMenu.description" :placeholder="$t('admin.appearance.menus.menu_desc_placeholder')" class="w-full" :rows="2" />
          </UFormField>
          <div class="flex gap-3 justify-end pt-2">
            <UButton color="neutral" variant="outline" @click="showCreateModal = false">
              {{ $t('common.cancel') }}
            </UButton>
            <UButton color="primary" :loading="isCreating" @click="createCustomMenu">
              {{ $t('common.create') }}
            </UButton>
          </div>
        </div>
      </template>
    </UModal>

    <!-- Delete confirm modal -->
    <UModal v-model:open="showDeleteModal">
      <template #content>
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
            </div>
            <div>
              <h3 class="font-semibold text-highlighted">{{ $t('admin.appearance.menus.delete_title') }}</h3>
              <p class="text-sm text-muted mt-0.5">
                {{ $t('admin.appearance.menus.delete_confirm', { name: selectedCustomMenu?.name }) }}
              </p>
            </div>
          </div>
          <div class="flex justify-end gap-2 mt-6">
            <UButton color="neutral" variant="outline" @click="showDeleteModal = false">
              {{ $t('common.cancel') }}
            </UButton>
            <UButton color="error" :loading="isDeleting" @click="deleteCustomMenu">
              {{ $t('common.confirm_delete') }}
            </UButton>
          </div>
        </div>
      </template>
    </UModal>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { NavMenuItem, NavMenuSlotKey, NavCustomMenu, UiMenuItem, MenuItemType } from '~/types/api/navMenu'
import { NAV_MENU_SLOTS, NAV_MENU_SLOT_HINTS, NAV_MENU_SLOT_CONFIGS } from '~/types/api/navMenu'

const toast = useToast()
const { t } = useI18n()
const navMenuApi = useNavMenuApi()
const optionApi = useOptionApi()
const optionsStore = useOptionsStore()
const { getOption } = useOption()
const contributionsStore = usePluginContributionsStore()

// ─── Plugin view visibility (per-slot) ─────────────────────────────────────

const disabledPluginViews = ref<string[]>([])

const currentSlotConfig = computed(() => NAV_MENU_SLOT_CONFIGS[selectedKey.value as NavMenuSlotKey])
const currentContribSlot = computed(() => currentSlotConfig.value?.contributionSlot)

watch(() => getOption('disabled_plugin_views'), (val) => {
  disabledPluginViews.value = [...val]
}, { immediate: true })

// ─── Loading ──────────────────────────────────────────────────────────────────

const rawLoading = ref(true)
const isLoading = useMinLoading(rawLoading)
const isSaving = ref(false)
const isCreating = ref(false)
const isDeleting = ref(false)

// ─── Slot data ────────────────────────────────────────────────────────────────

const SLOT_LIST = (Object.entries(NAV_MENU_SLOTS) as [NavMenuSlotKey, string][]).map(([key, label]) => ({ key, label }))
const slotData = ref<Record<string, NavMenuItem[]>>({} as Record<NavMenuSlotKey, NavMenuItem[]>)
const customMenus = ref<NavCustomMenu[]>([])

function slotItems(key: NavMenuSlotKey): NavMenuItem[] {
  return slotData.value[key] ?? []
}

// ─── Selection ────────────────────────────────────────────────────────────────

type SelectedKey = NavMenuSlotKey | string
const selectedKey = ref<SelectedKey>('')

const isBuiltinSelected = computed(() => selectedKey.value in NAV_MENU_SLOTS)
const selectedHint = computed(() =>
  isBuiltinSelected.value ? (NAV_MENU_SLOT_HINTS[selectedKey.value as NavMenuSlotKey] ?? '') : '',
)
const selectedLabel = computed(() => {
  if (isBuiltinSelected.value) return NAV_MENU_SLOTS[selectedKey.value as NavMenuSlotKey]
  return customMenus.value.find(m => m.id === selectedKey.value)?.name ?? ''
})
const selectedCustomMenu = computed(() => customMenus.value.find(m => m.id === selectedKey.value) ?? null)

// ─── Current items (editing) ──────────────────────────────────────────────────

const menuItems = ref<UiMenuItem[]>([])
let itemsSnapshot: NavMenuItem[] = []

function loadMenuItems(items: NavMenuItem[]) {
  itemsSnapshot = items
  menuItems.value = fromStored(items)
}

// ─── Convert stored ↔ UI ─────────────────────────────────────────────────────

function fromStored(items: NavMenuItem[]): UiMenuItem[] {
  const depthOf = new Map<string, number>()
  return items.map(it => {
    const d = it.parent_local_id ? (depthOf.get(it.parent_local_id) ?? 0) + 1 : 0
    depthOf.set(it.local_id, d)
    return {
      local_id: it.local_id,
      label: it.label,
      url: it.url,
      type: it.object_type,
      object_id: it.object_id,
      openInNewTab: it.target === '_blank',
      cssClasses: it.css_classes,
      depth: d,
      parent_local_id: it.parent_local_id,
      expanded: false,
    }
  })
}

function toStored(items: UiMenuItem[]): NavMenuItem[] {
  // Recalculate parent_local_id from position + depth so that drag-and-drop
  // reordering produces correct parent references (not stale ones).
  const parentStack: string[] = [] // parentStack[depth] = local_id of parent at that depth
  return items.map(it => {
    // For depth 0 items, no parent
    let parentId = ''
    if (it.depth > 0) {
      // The parent is the last item we saw at depth - 1
      parentId = parentStack[it.depth - 1] ?? ''
    }
    // Record this item as potential parent for deeper items
    parentStack[it.depth] = it.local_id
    return {
      local_id: it.local_id,
      label: it.label,
      url: it.url,
      object_type: it.type,
      object_id: it.object_id,
      target: it.openInNewTab ? '_blank' : '',
      css_classes: it.cssClasses,
      parent_local_id: parentId,
    }
  })
}

// ─── Load ─────────────────────────────────────────────────────────────────────

onMounted(async () => {
  rawLoading.value = true
  try {
    const slotPromises = SLOT_LIST.map(s => navMenuApi.loadSlot(s.key))
    const [allSlots, custom] = await Promise.all([
      Promise.all(slotPromises),
      navMenuApi.loadCustomMenus(),
    ])
    SLOT_LIST.forEach((s, i) => { slotData.value[s.key] = allSlots[i]! })
    customMenus.value = custom
    selectSlot(SLOT_LIST[0]!.key)
  } finally {
    rawLoading.value = false
  }
})

function selectSlot(key: NavMenuSlotKey) {
  selectedKey.value = key
  let items = slotData.value[key] ?? []
  if (!items.length) {
    const defaults = NAV_MENU_SLOT_CONFIGS[key]?.defaultItems
    if (defaults) items = defaults
  }
  loadMenuItems(items)
}

function selectCustom(id: string) {
  selectedKey.value = id
  loadMenuItems(customMenus.value.find(m => m.id === id)?.items ?? [])
}

// ─── Add items (from NavMenuAddPanel) ─────────────────────────────────────────

function makeId() {
  return `${Date.now()}-${Math.random().toString(36).slice(2)}`
}

function onAddItems(items: Array<{ label: string; url: string; type: MenuItemType; object_id: number; local_id?: string; css_classes?: string }>) {
  menuItems.value = [...menuItems.value, ...items.map(it => ({
    local_id: it.local_id || makeId(),
    label: it.label,
    url: it.url,
    type: it.type,
    object_id: it.object_id,
    openInNewTab: false,
    cssClasses: it.css_classes || '',
    depth: 0,
    parent_local_id: '',
    expanded: false,
  }))]
}

// ─── Reset ────────────────────────────────────────────────────────────────────

function resetCurrentMenu() {
  menuItems.value = fromStored(itemsSnapshot)
}

// ─── Save ─────────────────────────────────────────────────────────────────────

async function saveCurrentMenu() {
  if (!selectedKey.value) return
  isSaving.value = true
  try {
    const stored = toStored(menuItems.value)
    if (isBuiltinSelected.value) {
      await navMenuApi.saveSlot(selectedKey.value as NavMenuSlotKey, stored)
      slotData.value[selectedKey.value] = stored
    } else {
      const idx = customMenus.value.findIndex(m => m.id === selectedKey.value)
      if (idx !== -1) customMenus.value[idx]!.items = stored
      await navMenuApi.saveCustomMenus(customMenus.value)
    }
    // Sync disabled_plugin_views: plugin items not present in menu are disabled
    if (currentContribSlot.value) {
      const pluginViews = contributionsStore.getViewItems(currentContribSlot.value).value
      const pluginMenus = contributionsStore.getMenuItems(currentContribSlot.value).value
      const menuPluginIds = new Set(stored.filter(i => i.local_id.startsWith('plugin:')).map(i => i.local_id.replace('plugin:', '')))
      const allPluginIds = [...pluginViews.map(v => v.id), ...pluginMenus.map(m => m.command)]
      const otherDisabled = disabledPluginViews.value.filter(id => !allPluginIds.includes(id))
      const newDisabled = [...otherDisabled, ...allPluginIds.filter(id => !menuPluginIds.has(id))]
      if (JSON.stringify(newDisabled.sort()) !== JSON.stringify(disabledPluginViews.value.sort())) {
        disabledPluginViews.value = newDisabled
        await optionApi.setOption('disabled_plugin_views', newDisabled)
        await optionsStore.reload()
      }
    }
    itemsSnapshot = stored
    toast.add({ title: t('admin.appearance.menus.saved'), color: 'success' })
  } catch {
    toast.add({ title: t('admin.appearance.menus.save_failed'), color: 'error' })
  } finally {
    isSaving.value = false
  }
}

// ─── Create custom menu ───────────────────────────────────────────────────────

const showCreateModal = ref(false)
const newMenu = ref({ name: '', description: '' })

async function createCustomMenu() {
  if (!newMenu.value.name.trim()) {
    toast.add({ title: t('admin.appearance.menus.fill_menu_name'), color: 'warning' })
    return
  }
  isCreating.value = true
  try {
    const menu: NavCustomMenu = {
      id: makeId(),
      name: newMenu.value.name.trim(),
      description: newMenu.value.description,
      items: [],
    }
    customMenus.value.push(menu)
    await navMenuApi.saveCustomMenus(customMenus.value)
    showCreateModal.value = false
    newMenu.value = { name: '', description: '' }
    selectCustom(menu.id)
    toast.add({ title: t('admin.appearance.menus.created'), color: 'success' })
  } catch {
    toast.add({ title: t('admin.appearance.menus.create_failed'), color: 'error' })
  } finally {
    isCreating.value = false
  }
}

// ─── Delete custom menu ───────────────────────────────────────────────────────

const showDeleteModal = ref(false)

async function deleteCustomMenu() {
  if (!selectedKey.value || isBuiltinSelected.value) return
  isDeleting.value = true
  try {
    customMenus.value = customMenus.value.filter(m => m.id !== selectedKey.value)
    await navMenuApi.saveCustomMenus(customMenus.value)
    showDeleteModal.value = false
    selectSlot(SLOT_LIST[0]!.key)
    toast.add({ title: t('admin.appearance.menus.deleted'), color: 'success' })
  } catch {
    toast.add({ title: t('admin.appearance.menus.delete_failed'), color: 'error' })
  } finally {
    isDeleting.value = false
  }
}
</script>
