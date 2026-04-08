<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.settings.media.title')" :subtitle="$t('admin.settings.media.subtitle')">
      <template #actions>
        <UButton color="neutral" variant="outline" :disabled="isSaving" @click="loadAll">{{ $t('common.reset') }}</UButton>
        <UButton color="primary" icon="i-tabler-device-floppy" :loading="isSaving" @click="saveSettings">
          {{ $t('common.save_changes') }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <div v-if="isLoading" class="space-y-4">
        <UCard v-for="i in 3" :key="i">
          <template #header><USkeleton class="h-5 w-40" /></template>
          <div class="space-y-3">
            <USkeleton v-for="j in 3" :key="j" class="h-9 w-full" />
          </div>
        </UCard>
      </div>

      <template v-if="!isLoading">
        <!-- 文件大小限制 -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.media.file_limits') }}</h3>
          </template>
          <div class="space-y-4">
            <p class="text-xs text-muted">{{ $t('admin.settings.media.file_limits_hint') }}</p>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <UFormField v-for="ft in fileTypes" :key="ft.key" :label="ft.label">
                <div class="flex items-center gap-2">
                  <UInput v-model.number="form.limits[ft.key]" type="number" min="1" class="w-32" />
                  <span class="text-sm text-muted">MB</span>
                </div>
              </UFormField>
            </div>
          </div>
        </UCard>

        <!-- 缩略图生成 -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.media.thumb_gen') }}</h3>
          </template>
          <div class="space-y-6">
            <UAlert
              icon="i-tabler-info-circle"
              color="info"
              variant="soft"
              :title="$t('admin.settings.media.thumb_info_title')"
              :description="$t('admin.settings.media.thumb_info_desc')"
            />
            <div v-for="preset in thumbPresets" :key="preset.key" class="space-y-3">
              <div>
                <div class="text-sm font-medium text-highlighted">{{ preset.label }}</div>
                <div class="text-xs text-muted mt-0.5">{{ preset.desc }}</div>
              </div>
              <div class="flex items-center gap-3 flex-wrap">
                <UFormField :label="$t('admin.settings.media.width_label')">
                  <UInput v-model.number="form[preset.key].width" type="number" min="0" class="w-28" />
                </UFormField>
                <span class="text-muted mt-5">×</span>
                <UFormField :label="$t('admin.settings.media.height_label')">
                  <UInput v-model.number="form[preset.key].height" type="number" min="0" class="w-28" />
                </UFormField>
                <UBadge
                  class="mt-5"
                  :color="thumbBadgeColor(form[preset.key])"
                  variant="soft"
                  :label="thumbMode(form[preset.key])" />
              </div>
              <USeparator v-if="preset.key !== 'contentThumb'" />
            </div>

            <!-- variants 预览 -->
            <div class="rounded-md bg-muted p-4 text-xs space-y-1.5 font-mono">
              <p class="text-highlighted font-semibold mb-2 font-sans text-xs">{{ $t('admin.settings.media.variants_preview') }}</p>
              <template v-for="preset in thumbPresets" :key="preset.key">
                <p v-if="form[preset.key].width > 0 || form[preset.key].height > 0" class="text-muted">
                  <span class="text-highlighted">"{{ preset.variantKey }}"</span>: ".../xxx_{{ preset.variantKey }}.jpg"
                  <span class="ml-2 opacity-60">
                    {{ form[preset.key].width || 'auto' }}×{{ form[preset.key].height || 'auto' }}px
                  </span>
                </p>
              </template>
            </div>
          </div>
        </UCard>

        <!-- 上传目录结构 -->
        <UCard>
          <template #header>
            <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.media.upload_path') }}</h3>
          </template>
          <div class="space-y-4">
            <p class="text-xs text-muted">{{ $t('admin.settings.media.upload_path_hint') }}</p>
            <UFormField :label="$t('admin.settings.media.upload_path_label')">
              <div class="space-y-2">
                <USelect
                  v-model="uploadPathPreset"
                  :items="uploadPathPresets"
                  class="w-80" />
                <UInput
                  v-if="uploadPathPreset === '__custom__'"
                  v-model="form.uploadPath"
                  placeholder="{year}/{month}"
                  class="w-80 font-mono text-sm" />
              </div>
            </UFormField>
            <div class="rounded-md bg-muted px-4 py-3 text-xs font-mono">
              <span class="text-muted">{{ $t('admin.settings.media.upload_path_preview') }}: </span>
              <span class="text-highlighted">{{ uploadPathExample }}</span>
            </div>
          </div>
        </UCard>

        <!-- 媒体分类 & 存储 -->
        <UCard>
          <template #header>
            <div>
              <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.media.categories') }}</h3>
              <p class="text-xs text-muted mt-0.5">{{ $t('admin.settings.media.storage_per_category') }}</p>
            </div>
          </template>
          <div class="space-y-4">

            <!-- 分类表格 -->
            <div class="overflow-x-auto">
              <table class="w-full text-sm min-w-[560px]">
                <thead>
                  <tr class="border-b border-default text-left">
                    <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.category_slug') }}</th>
                    <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.category_label_zh') }}</th>
                    <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.category_label_en') }}</th>
                    <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.storage') }}</th>
                    <th class="pb-2 w-16"></th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-default">
                  <tr v-for="cat in categories" :key="cat.slug">
                    <!-- slug -->
                    <td class="py-2.5 pr-3">
                      <div class="flex items-center gap-1.5">
                        <span class="font-mono text-xs text-muted">{{ cat.slug }}</span>
                        <UBadge v-if="cat.is_system" :label="$t('admin.settings.media.category_built_in')" color="neutral" variant="soft" size="xs" />
                      </div>
                    </td>
                    <!-- zh label -->
                    <td class="py-2.5 pr-3 text-highlighted">{{ cat.label_zh }}</td>
                    <!-- en label -->
                    <td class="py-2.5 pr-3 text-muted">{{ cat.label_en }}</td>
                    <!-- storage -->
                    <td class="py-2.5 pr-3">
                      <USelect
                        :model-value="getCatStorage(cat.slug)"
                        :items="storageSelectItems"
                        size="xs"
                        class="w-44"
                        @update:model-value="setCatStorage(cat.slug, $event)" />
                    </td>
                    <!-- actions -->
                    <td class="py-2.5 text-center">
                      <UIcon name="i-tabler-lock" class="size-3.5 text-muted" />
                    </td>
                  </tr>
                  <tr v-if="categories.length === 0">
                    <td colspan="5" class="py-6 text-center text-xs text-muted">{{ $t('common.no_data') }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- 存储后端参考 -->
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
    </AdminPageContent>
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { MediaCategoryItem } from '~/composables/useMediaCategories'

interface FileLimits { image: number; video: number; audio: number; document: number; other: number }
interface ThumbSize  { width: number; height: number }
interface StorageBackend { name: string; display_name: string; type: string; enabled: boolean }
interface FormState {
  limits:         FileLimits
  thumbnailThumb: ThumbSize
  coverThumb:     ThumbSize
  contentThumb:   ThumbSize
  uploadPath:     string
}

const { t } = useI18n()

const fileTypes = computed(() => [
  { key: 'image'    as const, label: t('admin.settings.media.label_image') },
  { key: 'video'    as const, label: t('admin.settings.media.label_video') },
  { key: 'audio'    as const, label: t('admin.settings.media.label_audio') },
  { key: 'document' as const, label: t('admin.settings.media.label_document') },
  { key: 'other'    as const, label: t('admin.settings.media.label_other') },
])

const thumbPresets = computed(() => [
  { key: 'thumbnailThumb' as const, variantKey: 'thumbnail', label: t('admin.settings.media.thumb_general'),  desc: t('admin.settings.media.thumb_general_desc') },
  { key: 'coverThumb'     as const, variantKey: 'cover',     label: t('admin.settings.media.thumb_cover'),    desc: t('admin.settings.media.thumb_cover_desc') },
  { key: 'contentThumb'   as const, variantKey: 'content',   label: t('admin.settings.media.thumb_content'),  desc: t('admin.settings.media.thumb_content_desc') },
])

const DEFAULT_UPLOAD_PATH = '{year}/{month}'

const DEFAULT_LIMITS: FileLimits = { image: 10, video: 100, audio: 20, document: 20, other: 10 }
const DEFAULTS = {
  thumbnailThumb: { width: 300, height: 200 } as ThumbSize,
  coverThumb:     { width: 400, height: 300 } as ThumbSize,
  contentThumb:   { width: 1200, height: 0  } as ThumbSize,
}

const thumbMode = (s: ThumbSize) => {
  if (s.width > 0 && s.height > 0) return t('admin.settings.media.mode_crop')
  if (s.width > 0)                  return t('admin.settings.media.mode_width')
  if (s.height > 0)                 return t('admin.settings.media.mode_height')
  return t('admin.settings.media.mode_none')
}
const thumbBadgeColor = (s: ThumbSize) =>
  (s.width === 0 && s.height === 0) ? 'neutral' : 'primary'

const { apiFetch } = useApiFetch()
const toast        = useToast()
const isSaving     = ref(false)
const rawLoading   = ref(true)
const isLoading    = useMinLoading(rawLoading)

const form = ref<FormState>({
  limits:         { ...DEFAULT_LIMITS },
  thumbnailThumb: { ...DEFAULTS.thumbnailThumb },
  coverThumb:     { ...DEFAULTS.coverThumb },
  contentThumb:   { ...DEFAULTS.contentThumb },
  uploadPath:     DEFAULT_UPLOAD_PATH,
})

// ── Upload path presets ──────────────────────────────────────────────────────
const PRESET_VALUES = [
  '{year}/{month}',
  '{year}/{month}/{day}',
  '{category}',
  '{category}/{year}/{month}',
  'files',
]
const uploadPathPresets = computed(() => [
  { value: '{year}/{month}',          label: t('admin.settings.media.upload_path_preset_year_month') },
  { value: '{year}/{month}/{day}',    label: t('admin.settings.media.upload_path_preset_year_month_day') },
  { value: '{category}',             label: t('admin.settings.media.upload_path_preset_category') },
  { value: '{category}/{year}/{month}', label: t('admin.settings.media.upload_path_preset_category_year_month') },
  { value: 'files',                  label: t('admin.settings.media.upload_path_preset_flat') },
  { value: '__custom__',             label: t('admin.settings.media.upload_path_preset_custom') },
])

const uploadPathPreset = computed({
  get: () => PRESET_VALUES.includes(form.value.uploadPath) ? form.value.uploadPath : '__custom__',
  set: (val: string) => {
    if (val !== '__custom__') form.value.uploadPath = val
  },
})

const uploadPathExample = computed(() => {
  const now = new Date()
  const yyyy = String(now.getFullYear())
  const mm   = String(now.getMonth() + 1).padStart(2, '0')
  const dd   = String(now.getDate()).padStart(2, '0')
  const cat  = 'post_content'
  return (form.value.uploadPath || '{year}/{month}')
    .replace('{year}', yyyy)
    .replace('{month}', mm)
    .replace('{day}', dd)
    .replace('{category}', cat)
    + '/abc123.jpg'
})

// ── Categories ─────────────────────────────────────────────────────────────
// Categories are defined in server consts and synced to DB on startup.
// The only admin-configurable field at runtime is storage_key per category.
const categories      = ref<MediaCategoryItem[]>([])
const dirtyCategories = ref(new Set<string>()) // slugs with pending storage_key changes

// Sentinel to avoid empty-string SelectItem warning (Reka UI requirement)
const STORAGE_SENTINEL = '__default__'
const getCatStorage = (slug: string) => {
  const cat = categories.value.find(c => c.slug === slug)
  return cat?.storage_key || STORAGE_SENTINEL
}
const setCatStorage = (slug: string, val: string) => {
  const cat = categories.value.find(c => c.slug === slug)
  if (cat) {
    cat.storage_key = val === STORAGE_SENTINEL ? '' : val
    dirtyCategories.value.add(slug)
  }
}

// ── Storage backends ────────────────────────────────────────────────────────
const storageBackends = ref<StorageBackend[]>([])
const storageDefault  = ref('')

const storageSelectItems = computed(() => [
  {
    value: STORAGE_SENTINEL,
    label: storageDefault.value
      ? `${t('admin.settings.media.storage_use_default')} (${storageDefault.value})`
      : t('admin.settings.media.storage_use_default'),
  },
  ...storageBackends.value
    .filter(b => b.enabled)
    .map(b => ({ value: b.name, label: b.display_name || b.name })),
])

// ── Load / Save ─────────────────────────────────────────────────────────────
const loadSettings = async () => {
  const res  = await apiFetch<{ options: Record<string, string> }>('/options/autoload')
  const opts = res.options ?? {}
  const p    = <T>(key: string, fallback: T): T => {
    try { return JSON.parse(opts[key] ?? 'null') ?? fallback } catch { return fallback }
  }
  form.value.limits         = p('media_size_limits',   DEFAULT_LIMITS)
  form.value.thumbnailThumb = p('media_thumbnail',      DEFAULTS.thumbnailThumb)
  form.value.coverThumb     = p('media_cover_thumb',    DEFAULTS.coverThumb)
  form.value.contentThumb   = p('media_content_thumb',  DEFAULTS.contentThumb)
  form.value.uploadPath     = p<string>('media_upload_path', DEFAULT_UPLOAD_PATH) || DEFAULT_UPLOAD_PATH
}

const loadCategories = async () => {
  const res         = await apiFetch<{ list: MediaCategoryItem[] }>('/admin/media/categories')
  categories.value  = res.list ?? []
  dirtyCategories.value.clear()
}

const fetchStorageBackends = async () => {
  try {
    const res           = await apiFetch<{ backends: StorageBackend[]; default: string }>('/admin/storage/backends')
    storageBackends.value = res.backends ?? []
    storageDefault.value  = res.default  ?? ''
  } catch { /* non-critical */ }
}

const loadAll = async () => {
  rawLoading.value = true
  try {
    await Promise.all([loadSettings(), loadCategories(), fetchStorageBackends()])
  } catch {
    toast.add({ title: t('admin.settings.media.title'), description: t('common.load_failed'), color: 'error' })
  } finally {
    rawLoading.value = false
  }
}

const saveSettings = async () => {
  isSaving.value = true
  try {
    // 1. Save file limits + thumbnail presets + upload path
    await Promise.all([
      apiFetch('/options/media_size_limits',   { method: 'PUT', body: { value: JSON.stringify(form.value.limits),         autoload: 1 } }),
      apiFetch('/options/media_thumbnail',     { method: 'PUT', body: { value: JSON.stringify(form.value.thumbnailThumb), autoload: 1 } }),
      apiFetch('/options/media_cover_thumb',   { method: 'PUT', body: { value: JSON.stringify(form.value.coverThumb),     autoload: 1 } }),
      apiFetch('/options/media_content_thumb', { method: 'PUT', body: { value: JSON.stringify(form.value.contentThumb),   autoload: 1 } }),
      apiFetch('/options/media_upload_path',   { method: 'PUT', body: { value: JSON.stringify(form.value.uploadPath),     autoload: 1 } }),
    ])

    // 2. Flush storage_key changes for categories
    if (dirtyCategories.value.size > 0) {
      const dirty = categories.value.filter(c => dirtyCategories.value.has(c.slug))
      await Promise.all(
        dirty.map(c => apiFetch(`/admin/media/categories/${c.slug}`, {
          method: 'PUT',
          body: { storage_key: c.storage_key },
        })),
      )
      dirtyCategories.value.clear()
    }

    toast.add({ title: t('admin.settings.media.saved'), color: 'success' })
  } catch {
    toast.add({ title: t('admin.settings.media.save_failed'), color: 'error' })
  } finally {
    isSaving.value = false
  }
}

await loadAll()
</script>
