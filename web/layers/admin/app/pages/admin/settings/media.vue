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
        <UCard v-for="i in 4" :key="i">
          <template #header><USkeleton class="h-5 w-40" /></template>
          <div class="space-y-3">
            <USkeleton v-for="j in 3" :key="j" class="h-9 w-full" />
          </div>
        </UCard>
      </div>

      <template v-if="!isLoading">
        <!-- Extension Groups -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div>
                <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.media.extension_groups') }}</h3>
                <p class="text-xs text-muted mt-0.5">{{ $t('admin.settings.media.extension_groups_hint') }}</p>
              </div>
              <UButton color="primary" variant="soft" size="xs" icon="i-tabler-plus" @click="openGroupModal(-1)">
                {{ $t('admin.settings.media.add_group') }}
              </UButton>
            </div>
          </template>
          <div class="overflow-x-auto">
            <table class="w-full text-sm min-w-[500px]">
              <thead>
                <tr class="border-b border-default text-left">
                  <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.group_name') }}</th>
                  <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.group_extensions') }}</th>
                  <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.group_max_size') }}</th>
                  <th class="pb-2 w-20"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-default">
                <tr v-for="(grp, idx) in extensionGroups" :key="grp.name">
                  <td class="py-2.5 pr-3">
                    <div class="text-sm font-medium text-highlighted">{{ isZh ? grp.label_zh : grp.label_en }}</div>
                    <div class="text-xs text-muted font-mono">{{ grp.name }}</div>
                  </td>
                  <td class="py-2.5 pr-3">
                    <div class="flex flex-wrap gap-1">
                      <UBadge v-for="ext in grp.extensions" :key="ext" :label="'.' + ext" color="neutral" variant="soft" size="xs" />
                    </div>
                  </td>
                  <td class="py-2.5 pr-3 text-sm text-highlighted">{{ grp.max_size_mb }} MB</td>
                  <td class="py-2.5">
                    <div class="flex items-center gap-1">
                      <UButton color="neutral" variant="ghost" icon="i-tabler-pencil" size="xs" square @click="openGroupModal(idx)" />
                      <UButton color="error" variant="ghost" icon="i-tabler-trash" size="xs" square @click="removeExtensionGroup(idx)" />
                    </div>
                  </td>
                </tr>
                <tr v-if="extensionGroups.length === 0">
                  <td colspan="4" class="py-6 text-center text-xs text-muted">{{ $t('common.no_data') }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </UCard>

        <!-- Format Policies -->
        <UCard>
          <template #header>
            <div class="flex items-center justify-between">
              <div>
                <h3 class="text-lg font-semibold text-highlighted">{{ $t('admin.settings.media.format_policies') }}</h3>
                <p class="text-xs text-muted mt-0.5">{{ $t('admin.settings.media.format_policies_hint') }}</p>
              </div>
              <UButton color="primary" variant="soft" size="xs" icon="i-tabler-plus" @click="openPolicyModal(null)">
                {{ $t('admin.settings.media.add_policy') }}
              </UButton>
            </div>
          </template>
          <div class="overflow-x-auto">
            <table class="w-full text-sm min-w-[500px]">
              <thead>
                <tr class="border-b border-default text-left">
                  <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.group_name') }}</th>
                  <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.policy_groups') }}</th>
                  <th class="pb-2 pr-3 text-xs font-medium text-muted">{{ $t('admin.settings.media.system_policy') }}</th>
                  <th class="pb-2 w-20"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-default">
                <tr v-for="pol in formatPolicies" :key="pol.name">
                  <td class="py-2.5 pr-3">
                    <div class="text-sm font-medium text-highlighted">{{ isZh ? pol.label_zh : pol.label_en }}</div>
                    <div class="text-xs text-muted font-mono">{{ pol.name }}</div>
                  </td>
                  <td class="py-2.5 pr-3">
                    <div class="flex flex-wrap gap-1">
                      <UBadge v-for="g in pol.groups" :key="g" :label="getGroupLabel(g)" color="primary" variant="soft" size="xs" />
                    </div>
                  </td>
                  <td class="py-2.5 pr-3">
                    <UIcon v-if="pol.is_system" name="i-tabler-lock" class="size-4 text-muted" />
                  </td>
                  <td class="py-2.5">
                    <div class="flex items-center gap-1">
                      <UButton color="neutral" variant="ghost" icon="i-tabler-pencil" size="xs" square @click="openPolicyModal(pol)" />
                      <UButton
                        v-if="!pol.is_system"
                        color="error" variant="ghost" icon="i-tabler-trash" size="xs" square
                        @click="removeFormatPolicy(pol.name)" />
                    </div>
                  </td>
                </tr>
                <tr v-if="formatPolicies.length === 0">
                  <td colspan="4" class="py-6 text-center text-xs text-muted">{{ $t('common.no_data') }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </UCard>

        <!-- Thumbnail Config -->
        <MediaThumbConfig
          v-model:thumbnail-thumb="form.thumbnailThumb"
          v-model:cover-thumb="form.coverThumb"
          v-model:content-thumb="form.contentThumb"
        />

        <!-- Upload Path Config -->
        <MediaUploadPathConfig v-model="form.uploadPath" />

        <!-- Categories & Storage -->
        <MediaCategoryConfig
          :categories="categories"
          :format-policies="formatPolicies"
          :storage-backends="storageBackends"
          :storage-default="storageDefault"
          @update:category="onCategoryUpdate"
        />
      </template>
    </AdminPageContent>

    <!-- Extension Group Edit Modal -->
    <MediaExtensionGroupModal
      v-model:open="showGroupModal"
      :edit-index="editingGroupIndex"
      :group="editingGroupIndex >= 0 ? extensionGroups[editingGroupIndex] : null"
      @save="onGroupSaved"
    />

    <!-- Format Policy Edit Modal -->
    <MediaFormatPolicyModal
      v-model:open="showPolicyModal"
      :policy="editingPolicy"
      :extension-groups="extensionGroups"
      @saved="loadFormatPolicies"
    />
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { MediaCategoryItem } from '~/composables/useMediaCategories'
import type { ExtensionGroup, FormatPolicy } from '~/types/api/media'

interface ThumbSize  { width: number; height: number }
interface StorageBackend { name: string; display_name: string; type: string; enabled: boolean }
interface FormState {
  thumbnailThumb: ThumbSize
  coverThumb:     ThumbSize
  contentThumb:   ThumbSize
  uploadPath:     string
}

const { t, locale } = useI18n()
const isZh = computed(() => locale.value.startsWith('zh'))

const DEFAULT_UPLOAD_PATH = '{year}/{month}'

const DEFAULTS = {
  thumbnailThumb: { width: 300, height: 200 } as ThumbSize,
  coverThumb:     { width: 400, height: 300 } as ThumbSize,
  contentThumb:   { width: 1200, height: 0  } as ThumbSize,
}

const PRESET_VALUES = [
  '{year}/{month}',
  '{year}/{month}/{day}',
  '{category}',
  '{category}/{year}/{month}',
  'files',
]

const { apiFetch } = useApiFetch()
const mediaApi = useMediaApi()
const toast        = useToast()
const isSaving     = ref(false)
const rawLoading   = ref(true)
const isLoading    = useMinLoading(rawLoading)

const form = ref<FormState>({
  thumbnailThumb: { ...DEFAULTS.thumbnailThumb },
  coverThumb:     { ...DEFAULTS.coverThumb },
  contentThumb:   { ...DEFAULTS.contentThumb },
  uploadPath:     DEFAULT_UPLOAD_PATH,
})

// ── Extension Groups ────────────────────────────────────────────────────────
const extensionGroups = ref<ExtensionGroup[]>([])

const showGroupModal = ref(false)
const editingGroupIndex = ref(-1)

const openGroupModal = (idx: number) => {
  editingGroupIndex.value = idx
  showGroupModal.value = true
}

const onGroupSaved = (group: ExtensionGroup, index: number) => {
  if (index >= 0) {
    extensionGroups.value[index] = group
  } else {
    extensionGroups.value.push(group)
  }
}

const removeExtensionGroup = (idx: number) => {
  extensionGroups.value.splice(idx, 1)
}

const getGroupLabel = (name: string): string => {
  const grp = extensionGroups.value.find(g => g.name === name)
  if (!grp) return name
  return (isZh.value ? grp.label_zh : grp.label_en) || name
}

// ── Format Policies ─────────────────────────────────────────────────────────
const formatPolicies = ref<FormatPolicy[]>([])

const showPolicyModal = ref(false)
const editingPolicy = ref<FormatPolicy | null>(null)

const openPolicyModal = (pol: FormatPolicy | null) => {
  editingPolicy.value = pol
  showPolicyModal.value = true
}

const removeFormatPolicy = async (name: string) => {
  try {
    await mediaApi.deleteFormatPolicy(name)
    await loadFormatPolicies()
  } catch (err) {
    toast.add({ title: t('admin.settings.media.save_failed'), description: err instanceof Error ? err.message : undefined, color: 'error' })
  }
}

// ── Categories ─────────────────────────────────────────────────────────────
const categories      = ref<MediaCategoryItem[]>([])
const dirtyCategories = ref(new Set<string>())

const onCategoryUpdate = (slug: string, patch: { storage_key?: string; format_policy?: string; path_template?: string }) => {
  const cat = categories.value.find(c => c.slug === slug)
  if (!cat) return
  if (patch.storage_key !== undefined) cat.storage_key = patch.storage_key
  if (patch.format_policy !== undefined) cat.format_policy = patch.format_policy
  if (patch.path_template !== undefined) cat.path_template = patch.path_template
  dirtyCategories.value.add(slug)
}

// ── Storage backends ────────────────────────────────────────────────────────
const storageBackends = ref<StorageBackend[]>([])
const storageDefault  = ref('')

// ── Load / Save ─────────────────────────────────────────────────────────────
const loadSettings = async () => {
  const res  = await apiFetch<{ options: Record<string, string> }>('/options/autoload')
  const opts = res.options ?? {}
  const p    = <T>(key: string, fallback: T): T => {
    try { return JSON.parse(opts[key] ?? 'null') ?? fallback } catch { return fallback }
  }
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

const loadExtensionGroups = async () => {
  const res = await mediaApi.getExtensionGroups()
  extensionGroups.value = res.list ?? []
}

const loadFormatPolicies = async () => {
  const res = await mediaApi.getFormatPolicies()
  formatPolicies.value = res.list ?? []
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
    await Promise.all([loadSettings(), loadCategories(), loadExtensionGroups(), loadFormatPolicies(), fetchStorageBackends()])
  } catch {
    toast.add({ title: t('admin.settings.media.title'), description: t('common.load_failed'), color: 'error' })
  } finally {
    rawLoading.value = false
  }
}

const saveSettings = async () => {
  isSaving.value = true
  try {
    await Promise.all([
      mediaApi.saveExtensionGroups(extensionGroups.value),
      apiFetch('/options/media_thumbnail',     { method: 'PUT', body: { value: JSON.stringify(form.value.thumbnailThumb), autoload: 1 } }),
      apiFetch('/options/media_cover_thumb',   { method: 'PUT', body: { value: JSON.stringify(form.value.coverThumb),     autoload: 1 } }),
      apiFetch('/options/media_content_thumb', { method: 'PUT', body: { value: JSON.stringify(form.value.contentThumb),   autoload: 1 } }),
      apiFetch('/options/media_upload_path',   { method: 'PUT', body: { value: JSON.stringify(form.value.uploadPath),     autoload: 1 } }),
    ])

    if (dirtyCategories.value.size > 0) {
      const dirty = categories.value.filter(c => dirtyCategories.value.has(c.slug))
      await Promise.all(
        dirty.map(c => apiFetch(`/admin/media/categories/${c.slug}`, {
          method: 'PUT',
          body: {
            storage_key: c.storage_key,
            format_policy: c.format_policy || '',
            path_template: c.path_template || '',
          },
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
