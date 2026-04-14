<script setup lang="ts">
import type { MediaCategoryItem } from '~/composables/useMediaCategories'
import type { FormatPolicy } from '~/types/api/media'

interface StorageBackend { name: string; display_name: string; type: string; enabled: boolean }

const props = defineProps<{
  categories: MediaCategoryItem[]
  formatPolicies: FormatPolicy[]
  storageBackends: StorageBackend[]
  storageDefault: string
}>()

const emit = defineEmits<{
  (e: 'update:category', slug: string, patch: { storage_key?: string; format_policy?: string; path_template?: string }): void
}>()

const { t, locale } = useI18n()
const isZh = computed(() => locale.value.startsWith('zh'))

const STORAGE_SENTINEL = '__default__'
const PATH_SENTINEL = '__global__'

const PRESET_VALUES = [
  '{year}/{month}',
  '{year}/{month}/{day}',
  '{category}',
  '{category}/{year}/{month}',
  'files',
]
const CAT_PATH_PRESETS = [PATH_SENTINEL, ...PRESET_VALUES]

const customCatPaths = ref(new Set<string>())

// Initialize custom path tracking
watch(() => props.categories, (cats) => {
  const custom = new Set<string>()
  for (const cat of cats) {
    if (cat.path_template && !CAT_PATH_PRESETS.includes(cat.path_template)) {
      custom.add(cat.slug)
    }
  }
  customCatPaths.value = custom
}, { immediate: true })

const getCatStorage = (slug: string) => {
  const cat = props.categories.find(c => c.slug === slug)
  return cat?.storage_key || STORAGE_SENTINEL
}
const setCatStorage = (slug: string, val: string) => {
  emit('update:category', slug, { storage_key: val === STORAGE_SENTINEL ? '' : val })
}

const getCatFormatPolicy = (slug: string) => {
  const cat = props.categories.find(c => c.slug === slug)
  return cat?.format_policy || 'default'
}
const setCatFormatPolicy = (slug: string, val: string) => {
  emit('update:category', slug, { format_policy: val })
}

const getCatPathTemplate = (slug: string) => {
  const cat = props.categories.find(c => c.slug === slug)
  return cat?.path_template || PATH_SENTINEL
}
const setCatPathTemplate = (slug: string, val: string) => {
  emit('update:category', slug, { path_template: val === PATH_SENTINEL ? '' : val })
}

const getCatPathPreset = (slug: string) => {
  if (customCatPaths.value.has(slug)) return '__custom__'
  const val = getCatPathTemplate(slug)
  return CAT_PATH_PRESETS.includes(val) ? val : '__custom__'
}

const setCatPathPreset = (slug: string, val: string) => {
  if (val === '__custom__') {
    customCatPaths.value.add(slug)
  } else {
    customCatPaths.value.delete(slug)
    setCatPathTemplate(slug, val)
  }
}

const formatPolicySelectItems = computed(() =>
  props.formatPolicies.map(p => ({
    value: p.name,
    label: (isZh.value ? p.label_zh : p.label_en) || p.name,
  })),
)

const storageSelectItems = computed(() => [
  {
    value: STORAGE_SENTINEL,
    label: props.storageDefault
      ? `${t('admin.settings.media.storage_use_default')} (${props.storageDefault})`
      : t('admin.settings.media.storage_use_default'),
  },
  ...props.storageBackends
    .filter(b => b.enabled)
    .map(b => ({ value: b.name, label: b.display_name || b.name })),
])

const pathTemplateSelectItems = computed(() => [
  { value: PATH_SENTINEL, label: t('admin.settings.media.use_global_path') },
  ...PRESET_VALUES.map(v => ({ value: v, label: v })),
  { value: '__custom__', label: t('admin.settings.media.upload_path_preset_custom') },
])
</script>

<template>
  <UCard>
    <template #header>
      <div>
        <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.media.categories') }}</h3>
        <p class="text-xs text-muted mt-0.5">{{ $t('admin.settings.media.storage_per_category') }}</p>
      </div>
    </template>
    <div class="space-y-4">
      <!-- Category table -->
      <div class="overflow-x-auto">
        <table class="w-full text-sm min-w-[900px]">
          <thead>
            <tr class="border-b border-default text-left">
              <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.category_slug') }}</th>
              <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.category_label_zh') }}</th>
              <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.category_label_en') }}</th>
              <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.storage') }}</th>
              <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.format_policy') }}</th>
              <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.path_template') }}</th>
              <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.max_per_owner') }}</th>
              <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.plugin_id') }}</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-default">
            <tr v-for="cat in categories" :key="cat.slug">
              <td class="py-2.5 pr-3">
                <div class="flex items-center gap-1.5">
                  <span class="font-mono text-xs text-muted">{{ cat.slug }}</span>
                  <UBadge v-if="cat.is_system" :label="$t('admin.settings.media.category_built_in')" color="neutral" variant="soft" size="xs" />
                </div>
              </td>
              <td class="py-2.5 pr-3 text-highlighted">{{ cat.label_zh }}</td>
              <td class="py-2.5 pr-3 text-muted">{{ cat.label_en }}</td>
              <td class="py-2.5 pr-3">
                <USelect
                  :model-value="getCatStorage(cat.slug)"
                  :items="storageSelectItems"
                  size="xs"
                  class="w-36"
                  @update:model-value="setCatStorage(cat.slug, $event)" />
              </td>
              <td class="py-2.5 pr-3">
                <USelect
                  :model-value="getCatFormatPolicy(cat.slug)"
                  :items="formatPolicySelectItems"
                  size="xs"
                  class="w-36"
                  @update:model-value="setCatFormatPolicy(cat.slug, $event)" />
              </td>
              <td class="py-2.5 pr-3">
                <div class="flex flex-col gap-1">
                  <USelect
                    :model-value="getCatPathPreset(cat.slug)"
                    :items="pathTemplateSelectItems"
                    size="xs"
                    class="w-44"
                    @update:model-value="setCatPathPreset(cat.slug, $event)" />
                  <template v-if="getCatPathPreset(cat.slug) === '__custom__'">
                    <UInput
                      :model-value="getCatPathTemplate(cat.slug)"
                      size="xs"
                      class="w-44 font-mono text-xs"
                      placeholder="{category}/{year}/{month}"
                      @update:model-value="setCatPathTemplate(cat.slug, String($event))" />
                    <p class="text-[10px] text-muted leading-tight">
                      <code>{year}</code> <code>{month}</code> <code>{day}</code> <code>{category}</code>
                    </p>
                  </template>
                </div>
              </td>
              <td class="py-2.5 pr-3 text-center">
                <span class="text-xs text-muted">{{ cat.max_per_owner || '—' }}</span>
              </td>
              <td class="py-2.5 pr-3">
                <UBadge v-if="cat.plugin_id" :label="cat.plugin_id" color="info" variant="soft" size="xs" />
                <span v-else class="text-xs text-muted">—</span>
              </td>
            </tr>
            <tr v-if="categories.length === 0">
              <td colspan="8" class="py-6 text-center text-xs text-muted">{{ $t('common.no_data') }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Storage backends reference -->
      <div class="rounded-md bg-muted p-3 space-y-1.5">
        <p class="text-xs font-medium text-highlighted mb-2">{{ $t('admin.settings.media.storage') }}</p>
        <p class="text-xs text-muted mb-2">{{ $t('admin.settings.media.storage_desc') }}</p>
        <div v-if="storageBackends.length > 0" class="flex flex-wrap gap-2">
          <div
            v-for="b in storageBackends"
            :key="b.name"
            class="flex items-center gap-1.5 px-2 py-1 rounded border border-default bg-default text-xs">
            <span class="font-medium text-highlighted">{{ b.display_name || b.name }}</span>
            <span class="text-muted font-mono">{{ b.name }}</span>
            <UBadge
              v-if="b.name === storageDefault"
              :label="$t('admin.settings.media.storage_default_label')"
              color="primary" variant="soft" size="xs" />
            <UBadge
              :label="$t('admin.settings.media.storage_type') + ': ' + b.type"
              color="neutral" variant="soft" size="xs" />
            <UBadge v-if="!b.enabled" label="Disabled" color="error" variant="soft" size="xs" />
          </div>
        </div>
        <p v-else class="text-xs text-muted">—</p>
      </div>
    </div>
  </UCard>
</template>
