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
                <template v-if="uploadPathPreset === '__custom__'">
                  <UInput
                    v-model="form.uploadPath"
                    placeholder="{year}/{month}"
                    class="w-80 font-mono text-sm" />
                  <p class="text-xs text-muted">
                    {{ $t('admin.settings.media.path_vars_hint') }}:
                    <code class="font-mono text-highlighted">{year}</code>
                    <code class="font-mono text-highlighted">{month}</code>
                    <code class="font-mono text-highlighted">{day}</code>
                    <code class="font-mono text-highlighted">{category}</code>
                  </p>
                </template>
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
                        class="w-36"
                        @update:model-value="setCatStorage(cat.slug, $event)" />
                    </td>
                    <!-- format policy -->
                    <td class="py-2.5 pr-3">
                      <USelect
                        :model-value="getCatFormatPolicy(cat.slug)"
                        :items="formatPolicySelectItems"
                        size="xs"
                        class="w-36"
                        @update:model-value="setCatFormatPolicy(cat.slug, $event)" />
                    </td>
                    <!-- path template -->
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
                            @update:model-value="setCatPathTemplate(cat.slug, $event)" />
                          <p class="text-[10px] text-muted leading-tight">
                            <code>{year}</code> <code>{month}</code> <code>{day}</code> <code>{category}</code>
                          </p>
                        </template>
                      </div>
                    </td>
                    <!-- max_per_owner -->
                    <td class="py-2.5 pr-3 text-center">
                      <span class="text-xs text-muted">{{ cat.max_per_owner || '—' }}</span>
                    </td>
                    <!-- plugin_id -->
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

    <!-- Extension Group Edit Modal -->
    <UModal v-model:open="showGroupModal" :title="editingGroupIndex < 0 ? $t('admin.settings.media.add_group') : $t('admin.settings.media.edit_group')">
      <template #content>
        <div class="p-6 space-y-4">
          <UFormField :label="$t('admin.settings.media.group_name')">
            <UInput v-model="groupForm.name" :disabled="editingGroupIndex >= 0" class="w-full" />
          </UFormField>
          <div class="grid grid-cols-2 gap-3">
            <UFormField :label="$t('admin.settings.media.category_label_zh')">
              <UInput v-model="groupForm.label_zh" class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.settings.media.category_label_en')">
              <UInput v-model="groupForm.label_en" class="w-full" />
            </UFormField>
          </div>
          <UFormField :label="$t('admin.settings.media.group_extensions')">
            <UInput v-model="groupForm.extensionsStr" placeholder="jpg,jpeg,png,webp,gif" class="w-full" />
            <template #hint>{{ $t('admin.settings.media.group_extensions_hint') }}</template>
          </UFormField>
          <UFormField :label="$t('admin.settings.media.group_max_size')">
            <div class="flex items-center gap-2">
              <UInput v-model.number="groupForm.max_size_mb" type="number" min="1" class="w-32" />
              <span class="text-sm text-muted">MB</span>
            </div>
          </UFormField>
          <div class="flex justify-end gap-2">
            <UButton color="neutral" variant="outline" @click="showGroupModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="primary" @click="saveGroupForm">{{ $t('common.confirm') }}</UButton>
          </div>
        </div>
      </template>
    </UModal>

    <!-- Format Policy Edit Modal -->
    <UModal v-model:open="showPolicyModal" :title="editingPolicyIsNew ? $t('admin.settings.media.add_policy') : $t('admin.settings.media.edit_policy')">
      <template #content>
        <div class="p-6 space-y-4">
          <UFormField :label="$t('admin.settings.media.group_name')">
            <UInput v-model="policyForm.name" :disabled="!editingPolicyIsNew" class="w-full" />
          </UFormField>
          <div class="grid grid-cols-2 gap-3">
            <UFormField :label="$t('admin.settings.media.category_label_zh')">
              <UInput v-model="policyForm.label_zh" :disabled="policyForm.is_system" class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.settings.media.category_label_en')">
              <UInput v-model="policyForm.label_en" :disabled="policyForm.is_system" class="w-full" />
            </UFormField>
          </div>
          <UFormField :label="$t('admin.settings.media.policy_groups')">
            <div class="space-y-2">
              <label v-for="grp in extensionGroups" :key="grp.name" class="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  :checked="policyForm.groups.includes(grp.name)"
                  class="rounded border-default"
                  @change="togglePolicyGroup(grp.name)" />
                <span class="text-sm">{{ isZh ? grp.label_zh : grp.label_en }} <span class="text-muted font-mono text-xs">({{ grp.name }})</span></span>
              </label>
            </div>
          </UFormField>
          <div class="flex justify-end gap-2">
            <UButton color="neutral" variant="outline" @click="showPolicyModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="primary" @click="savePolicyForm">{{ $t('common.confirm') }}</UButton>
          </div>
        </div>
      </template>
    </UModal>
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

const thumbPresets = computed(() => [
  { key: 'thumbnailThumb' as const, variantKey: 'thumbnail', label: t('admin.settings.media.thumb_general'),  desc: t('admin.settings.media.thumb_general_desc') },
  { key: 'coverThumb'     as const, variantKey: 'cover',     label: t('admin.settings.media.thumb_cover'),    desc: t('admin.settings.media.thumb_cover_desc') },
  { key: 'contentThumb'   as const, variantKey: 'content',   label: t('admin.settings.media.thumb_content'),  desc: t('admin.settings.media.thumb_content_desc') },
])

const DEFAULT_UPLOAD_PATH = '{year}/{month}'

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

const customPathActive = ref(false)

const uploadPathPreset = computed({
  get: () => {
    if (customPathActive.value) return '__custom__'
    return PRESET_VALUES.includes(form.value.uploadPath) ? form.value.uploadPath : '__custom__'
  },
  set: (val: string) => {
    if (val === '__custom__') {
      customPathActive.value = true
    } else {
      customPathActive.value = false
      form.value.uploadPath = val
    }
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

// ── Extension Groups ────────────────────────────────────────────────────────
const extensionGroups = ref<ExtensionGroup[]>([])

const showGroupModal = ref(false)
const editingGroupIndex = ref(-1) // -1 = new
const groupForm = ref({
  name: '',
  label_zh: '',
  label_en: '',
  extensionsStr: '',
  max_size_mb: 10,
})

const openGroupModal = (idx: number) => {
  editingGroupIndex.value = idx
  if (idx >= 0) {
    const g = extensionGroups.value[idx]
    groupForm.value = {
      name: g.name,
      label_zh: g.label_zh,
      label_en: g.label_en,
      extensionsStr: g.extensions.join(', '),
      max_size_mb: g.max_size_mb,
    }
  } else {
    groupForm.value = { name: '', label_zh: '', label_en: '', extensionsStr: '', max_size_mb: 10 }
  }
  showGroupModal.value = true
}

const saveGroupForm = () => {
  const f = groupForm.value
  if (!f.name.trim()) return
  const extensions = f.extensionsStr.split(/[,\s]+/).map(s => s.trim().replace(/^\./, '').toLowerCase()).filter(Boolean)
  const group: ExtensionGroup = {
    name: f.name.trim(),
    label_zh: f.label_zh.trim(),
    label_en: f.label_en.trim(),
    extensions,
    max_size_mb: f.max_size_mb,
  }
  if (editingGroupIndex.value >= 0) {
    extensionGroups.value[editingGroupIndex.value] = group
  } else {
    extensionGroups.value.push(group)
  }
  showGroupModal.value = false
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
const editingPolicyIsNew = ref(false)
const editingPolicyOriginalName = ref('')
const policyForm = ref({
  name: '',
  label_zh: '',
  label_en: '',
  is_system: false,
  groups: [] as string[],
})

const openPolicyModal = (pol: FormatPolicy | null) => {
  if (pol) {
    editingPolicyIsNew.value = false
    editingPolicyOriginalName.value = pol.name
    policyForm.value = {
      name: pol.name,
      label_zh: pol.label_zh,
      label_en: pol.label_en,
      is_system: pol.is_system,
      groups: [...pol.groups],
    }
  } else {
    editingPolicyIsNew.value = true
    editingPolicyOriginalName.value = ''
    policyForm.value = { name: '', label_zh: '', label_en: '', is_system: false, groups: [] }
  }
  showPolicyModal.value = true
}

const togglePolicyGroup = (name: string) => {
  const idx = policyForm.value.groups.indexOf(name)
  if (idx >= 0) policyForm.value.groups.splice(idx, 1)
  else policyForm.value.groups.push(name)
}

const savePolicyForm = async () => {
  const f = policyForm.value
  if (!f.name.trim()) return
  try {
    if (editingPolicyIsNew.value) {
      await mediaApi.createFormatPolicy({
        name: f.name.trim(),
        label_zh: f.label_zh.trim(),
        label_en: f.label_en.trim(),
        groups: f.groups,
      })
    } else {
      await mediaApi.updateFormatPolicy(editingPolicyOriginalName.value, {
        name: f.name.trim(),
        label_zh: f.label_zh.trim(),
        label_en: f.label_en.trim(),
        groups: f.groups,
      })
    }
    await loadFormatPolicies()
    showPolicyModal.value = false
  } catch (err) {
    toast.add({ title: t('admin.settings.media.save_failed'), description: err instanceof Error ? err.message : undefined, color: 'error' })
  }
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
const dirtyCategories = ref(new Set<string>()) // slugs with pending changes

// Sentinel to avoid empty-string SelectItem warning (Reka UI requirement)
const STORAGE_SENTINEL = '__default__'
const POLICY_SENTINEL = '__default__'
const PATH_SENTINEL = '__global__'

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

const getCatFormatPolicy = (slug: string) => {
  const cat = categories.value.find(c => c.slug === slug)
  return cat?.format_policy || 'default'
}
const setCatFormatPolicy = (slug: string, val: string) => {
  const cat = categories.value.find(c => c.slug === slug)
  if (cat) {
    cat.format_policy = val
    dirtyCategories.value.add(slug)
  }
}

const getCatPathTemplate = (slug: string) => {
  const cat = categories.value.find(c => c.slug === slug)
  return cat?.path_template || PATH_SENTINEL
}
const setCatPathTemplate = (slug: string, val: string) => {
  const cat = categories.value.find(c => c.slug === slug)
  if (cat) {
    cat.path_template = val === PATH_SENTINEL ? '' : val
    dirtyCategories.value.add(slug)
  }
}

const formatPolicySelectItems = computed(() =>
  formatPolicies.value.map(p => ({
    value: p.name,
    label: (isZh.value ? p.label_zh : p.label_en) || p.name,
  })),
)

const CAT_PATH_PRESETS = [PATH_SENTINEL, ...PRESET_VALUES]
const customCatPaths = ref(new Set<string>()) // slugs currently in custom mode

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

const pathTemplateSelectItems = computed(() => [
  { value: PATH_SENTINEL, label: t('admin.settings.media.use_global_path') },
  ...PRESET_VALUES.map(v => ({ value: v, label: v })),
  { value: '__custom__', label: t('admin.settings.media.upload_path_preset_custom') },
])

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
  form.value.thumbnailThumb = p('media_thumbnail',      DEFAULTS.thumbnailThumb)
  form.value.coverThumb     = p('media_cover_thumb',    DEFAULTS.coverThumb)
  form.value.contentThumb   = p('media_content_thumb',  DEFAULTS.contentThumb)
  form.value.uploadPath     = p<string>('media_upload_path', DEFAULT_UPLOAD_PATH) || DEFAULT_UPLOAD_PATH
  customPathActive.value    = !PRESET_VALUES.includes(form.value.uploadPath)
}

const loadCategories = async () => {
  const res         = await apiFetch<{ list: MediaCategoryItem[] }>('/admin/media/categories')
  categories.value  = res.list ?? []
  dirtyCategories.value.clear()
  // sync custom path tracking
  const custom = new Set<string>()
  for (const cat of categories.value) {
    if (cat.path_template && !CAT_PATH_PRESETS.includes(cat.path_template)) {
      custom.add(cat.slug)
    }
  }
  customCatPaths.value = custom
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
    // 1. Save extension groups + thumbnail presets + upload path
    await Promise.all([
      mediaApi.saveExtensionGroups(extensionGroups.value),
      apiFetch('/options/media_thumbnail',     { method: 'PUT', body: { value: JSON.stringify(form.value.thumbnailThumb), autoload: 1 } }),
      apiFetch('/options/media_cover_thumb',   { method: 'PUT', body: { value: JSON.stringify(form.value.coverThumb),     autoload: 1 } }),
      apiFetch('/options/media_content_thumb', { method: 'PUT', body: { value: JSON.stringify(form.value.contentThumb),   autoload: 1 } }),
      apiFetch('/options/media_upload_path',   { method: 'PUT', body: { value: JSON.stringify(form.value.uploadPath),     autoload: 1 } }),
    ])

    // 2. Flush category changes (storage_key, format_policy, path_template)
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
