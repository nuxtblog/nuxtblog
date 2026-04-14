<template>
  <UModal
    :open="localOpen"
    :title="activeTab === 'upload' ? $t('admin.media.upload_title') : $t('admin.media.import_url_tab')"
    :ui="{ content: 'max-w-2xl' }"
    @update:open="onModalUpdate">
    <template #content>
      <div class="p-6 space-y-4">
        <!-- Tab 切换 -->
        <div class="flex border-b border-default -mt-1 mb-2">
          <button
            class="px-4 py-2 text-sm font-medium transition-colors"
            :class="activeTab === 'upload' ? 'text-primary border-b-2 border-primary' : 'text-muted hover:text-highlighted'"
            @click="activeTab = 'upload'">
            {{ $t('common.upload') }}
          </button>
          <button
            class="px-4 py-2 text-sm font-medium transition-colors"
            :class="activeTab === 'url' ? 'text-primary border-b-2 border-primary' : 'text-muted hover:text-highlighted'"
            @click="activeTab = 'url'">
            {{ $t('admin.media.import_url_tab') }}
          </button>
        </div>

        <!-- URL 导入面板 -->
        <div v-if="activeTab === 'url'" class="space-y-4">
          <UFormField :label="$t('admin.media.import_url_label')">
            <UInput
              v-model="urlImport.url"
              :placeholder="$t('admin.media.import_url_placeholder')"
              class="w-full" />
          </UFormField>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <UFormField :label="$t('admin.media.import_url_title_label')">
              <UInput v-model="urlImport.title" class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.media.category')">
              <USelect
                v-model="urlImport.category"
                :items="categoryOptions"
                :placeholder="$t('admin.media.category_placeholder')"
                class="w-full" />
            </UFormField>
          </div>
          <UFormField :label="$t('admin.media.import_url_alt_label')">
            <UInput v-model="urlImport.altText" class="w-full" />
          </UFormField>
          <div class="flex justify-end gap-2">
            <UButton color="neutral" variant="outline" @click="localOpen = false">{{ $t('common.cancel') }}</UButton>
            <UButton
              color="primary"
              :disabled="!urlImport.url.trim()"
              :loading="urlImporting"
              @click="submitUrlImport">
              {{ urlImporting ? $t('admin.media.import_url_importing') : $t('admin.media.import_url_submit') }}
            </UButton>
          </div>
        </div>

        <!-- 上传面板 -->
        <template v-if="activeTab === 'upload'">
        <!-- 拖拽区（未开始上传时显示） -->
        <div
          v-if="!isUploading"
          :class="[
            'border-2 border-dashed rounded-md p-8 text-center transition-colors cursor-pointer',
            isDragging ? 'border-primary bg-primary/5' : 'border-default hover:border-primary/50',
          ]"
          @drop.prevent="handleDrop"
          @dragover.prevent="isDragging = true"
          @dragleave="isDragging = false"
          @click="fileInput?.click()">
          <input
            ref="fileInput"
            type="file"
            multiple
            accept="*/*"
            class="hidden"
            @change="handleFileSelect" />
          <UIcon name="i-tabler-upload" class="size-8 mx-auto text-muted mb-3" />
          <p class="text-sm text-muted mb-2">{{ $t('admin.media.upload_area_title') }}</p>
          <p class="text-xs text-muted">{{ $t('admin.media.upload_area_desc') }} · {{ uploadLimitHint }}</p>
        </div>

        <!-- 暂存文件列表（待填写信息） -->
        <div v-if="stagedFiles.length > 0 && !isUploading" class="space-y-3">
          <!-- 标题行 -->
          <div class="flex items-center justify-between">
            <h4 class="text-sm font-medium text-highlighted">{{ $t('admin.media.select_files') }} ({{ stagedFiles.length }})</h4>
            <UButton color="neutral" variant="ghost" size="xs" icon="i-tabler-plus" @click="fileInput?.click()">{{ $t('common.add') }}</UButton>
          </div>

          <!-- 批量工具栏 -->
          <div class="flex flex-wrap items-center gap-2 p-2 bg-elevated rounded-md">
            <span class="text-xs text-muted shrink-0">{{ $t('admin.media.move_to_cat') }}:</span>
            <USelect
              v-model="batchCategory"
              :items="categoryOptions"
              :placeholder="$t('admin.media.category_placeholder')"
              size="xs"
              class="w-36" />
            <UButton
              v-if="stagedFiles.some(f => f.file.type.startsWith('image/'))"
              color="neutral" variant="outline" size="xs" icon="i-tabler-tag"
              @click="applyAltFromTitle">
              Alt = {{ $t('common.title') }}
            </UButton>
          </div>

          <div class="space-y-3 max-h-[40vh] overflow-y-auto pr-1">
            <div
              v-for="sf in stagedFiles"
              :key="sf.id"
              class="rounded-md border p-3 space-y-2"
              :class="sf.error ? 'border-error/50 bg-error/5' : 'border-default bg-elevated'">
              <div class="flex items-start gap-3">
                <!-- 预览 -->
                <div class="shrink-0 w-16 h-16 rounded-md overflow-hidden bg-muted flex items-center justify-center">
                  <img
                    v-if="sf.previewUrl"
                    :src="sf.previewUrl"
                    :alt="sf.file.name"
                    class="w-full h-full object-cover" />
                  <UIcon v-else-if="sf.file.type.startsWith('video/')" name="i-tabler-video" class="size-7 text-muted" />
                  <UIcon v-else-if="sf.file.type.startsWith('audio/')" name="i-tabler-music" class="size-7 text-muted" />
                  <UIcon v-else name="i-tabler-file" class="size-7 text-muted" />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="text-sm font-medium text-highlighted truncate">{{ sf.file.name }}</div>
                  <div class="text-xs text-muted">{{ formatFileSize(sf.file.size) }}</div>
                  <div v-if="sf.error" class="text-xs text-error mt-0.5">{{ sf.error }}</div>
                </div>
                <UButton
                  color="neutral" variant="ghost" icon="i-tabler-x" square size="xs"
                  @click="removeStagedFile(sf.id)" />
              </div>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-2 mt-1">
                <UFormField :label="$t('common.title')">
                  <UInput v-model="sf.title" :placeholder="sf.file.name" size="xs" class="w-full" />
                </UFormField>
                <UFormField :label="$t('admin.media.category')">
                  <USelect
                    v-model="sf.category"
                    :items="categoryOptions"
                    size="xs"
                    class="w-full" />
                </UFormField>
              </div>
              <UFormField v-if="sf.file.type.startsWith('image/')" :label="$t('admin.media.alt_text')">
                <UInput v-model="sf.altText" :placeholder="$t('admin.media.alt_placeholder')" size="xs" class="w-full" />
              </UFormField>
            </div>
          </div>  <!-- end scroll list -->
        </div>  <!-- end staged section -->

        <!-- 上传进度（上传中） -->
        <div v-if="isUploading && uploadingFiles.length > 0" class="space-y-3 max-h-[50vh] overflow-y-auto pr-1">
          <div class="flex items-center justify-between">
            <h4 class="text-sm font-medium text-highlighted">{{ $t('admin.media.uploading') }}</h4>
            <span class="text-sm font-medium text-primary">{{ uploadDoneCount }} / {{ uploadingFiles.length }}</span>
          </div>
          <div
            v-for="file in uploadingFiles"
            :key="file.id"
            class="flex items-center gap-3 p-3 bg-elevated rounded-md">
            <div class="shrink-0 w-10 h-10 rounded overflow-hidden bg-muted flex items-center justify-center">
              <img v-if="file.previewUrl" :src="file.previewUrl" class="w-full h-full object-cover" />
              <UIcon v-else name="i-tabler-file" class="size-5 text-muted" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="text-sm font-medium text-highlighted truncate">{{ file.name }}</div>
              <div class="flex items-center gap-2 mt-1">
                <UProgress :value="file.progress" color="primary" size="xs" class="flex-1" />
                <span class="text-xs text-muted w-8 text-right">{{ file.progress }}%</span>
              </div>
              <div v-if="file.error" class="text-xs text-error mt-0.5">{{ file.error }}</div>
            </div>
            <UIcon
              v-if="file.progress === 100"
              name="i-tabler-circle-check"
              class="size-5 text-success shrink-0" />
            <UIcon
              v-else-if="file.error"
              name="i-tabler-alert-circle"
              class="size-5 text-error shrink-0" />
          </div>
        </div>

        <div class="flex justify-end gap-2">
          <UButton color="neutral" variant="outline" @click="localOpen = false">
            {{ isUploading ? $t('admin.media.upload_complete') : $t('common.cancel') }}
          </UButton>
          <UButton
            v-if="stagedFiles.length > 0 && !isUploading"
            color="primary"
            :disabled="stagedFiles.every(f => !!f.error)"
            @click="startUpload">
            {{ $t('common.upload') }} ({{ stagedFiles.filter(f => !f.error).length }})
          </UButton>
        </div>
        </template> <!-- end upload panel -->
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { ExtensionGroup, FormatPolicy } from '~/types/api/media'

interface Props {
  open: boolean
  categoryOptions: Array<{ value: string; label: string }>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  'uploaded': []
}>()

const localOpen = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val),
})

const mediaStore = useMediaStore()
const mediaApi = useMediaApi()
const optionsStore = useOptionsStore()
const toast = useToast()
const { t } = useI18n()

// ── Tab state ─────────────────────────────────────────────────────────────
const activeTab = ref<'upload' | 'url'>('upload')

// ── URL import state ──────────────────────────────────────────────────────
const urlImport = ref({ url: '', title: '', altText: '', category: '' })
const urlImporting = ref(false)

const submitUrlImport = async () => {
  if (!urlImport.value.url.trim()) return
  urlImporting.value = true
  try {
    await mediaApi.link({
      url: urlImport.value.url.trim(),
      title: urlImport.value.title || undefined,
      alt_text: urlImport.value.altText || undefined,
      category: (urlImport.value.category as any) || undefined,
    })
    toast.add({ title: t('admin.media.import_url_success'), icon: 'i-tabler-circle-check', color: 'success' })
    urlImport.value = { url: '', title: '', altText: '', category: '' }
    emit('update:open', false)
    emit('uploaded')
  } catch (err) {
    toast.add({
      title: t('admin.media.import_url_failed'),
      description: err instanceof Error ? err.message : undefined,
      color: 'error',
      icon: 'i-tabler-alert-circle',
    })
  } finally {
    urlImporting.value = false
  }
}

// ── State ──────────────────────────────────────────────────────────────────
const isDragging = ref(false)
const fileInput = ref<HTMLInputElement>()
const isUploading = ref(false)
const batchCategory = ref('')

interface StagedFile {
  id: string
  file: File
  previewUrl: string | null
  title: string
  altText: string
  category: string
  error: string | null
}
interface UploadingFile { id: string; name: string; progress: number; previewUrl: string | null; error: string | null }

const stagedFiles = ref<StagedFile[]>([])
const uploadingFiles = ref<UploadingFile[]>([])
const uploadDoneCount = computed(() => uploadingFiles.value.filter(f => f.progress === 100 || !!f.error).length)

// ── Extension-based validation ─────────────────────────────────────────────
const extensionGroups = ref<ExtensionGroup[]>([])
const formatPolicies = ref<FormatPolicy[]>([])

const loadFormatData = async () => {
  try {
    const [egRes, fpRes] = await Promise.all([
      mediaApi.getExtensionGroups?.() ?? Promise.resolve({ list: [] }),
      mediaApi.getFormatPolicies?.() ?? Promise.resolve({ list: [] }),
    ])
    extensionGroups.value = egRes.list ?? []
    formatPolicies.value = fpRes.list ?? []
  } catch {
    // fallback: allow all
  }
}

// Load format data on mount
onMounted(() => loadFormatData())

/** Get effective policy for the default category, collecting all allowed extensions & limits */
const getAllowedExtensions = computed(() => {
  const policy = formatPolicies.value.find(p => p.name === 'default') ?? formatPolicies.value[0]
  if (!policy) return new Map<string, number>()

  const extMap = new Map<string, number>() // ext -> maxSizeMB
  const groupMap = new Map(extensionGroups.value.map(g => [g.name, g]))

  for (const gName of policy.groups) {
    const grp = groupMap.get(gName)
    if (!grp) continue
    for (const ext of grp.extensions) {
      extMap.set(ext, grp.max_size_mb)
    }
  }
  return extMap
})

const checkFileSize = (file: File): string | null => {
  const ext = file.name.split('.').pop()?.toLowerCase() ?? ''
  const allowedExts = getAllowedExtensions.value
  if (allowedExts.size > 0 && !allowedExts.has(ext)) {
    return `.${ext} not allowed`
  }
  const limitMB = allowedExts.get(ext) ?? 10
  if (file.size > limitMB * 1024 * 1024) {
    return `${formatFileSize(file.size)} > ${limitMB} MB`
  }
  return null
}

const uploadLimitHint = computed(() => {
  const groups = extensionGroups.value
  if (!groups.length) return ''
  return groups.map(g => `${g.label_en || g.name} ${g.max_size_mb}MB`).join(' · ')
})

// ── Formatters ────────────────────────────────────────────────────────────
const formatFileSize = (bytes: number): string => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${Math.round((bytes / Math.pow(k, i)) * 100) / 100} ${sizes[i]}`
}

// ── File staging ──────────────────────────────────────────────────────────
const handleFileSelect = (event: Event) => {
  const files = (event.target as HTMLInputElement).files
  if (files) stageFiles(Array.from(files))
  if (fileInput.value) fileInput.value.value = ''
}

const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  const files = event.dataTransfer?.files
  if (files) stageFiles(Array.from(files))
}

const stageFiles = (files: File[]) => {
  for (const file of files) {
    const id = `${Date.now()}-${Math.random()}`
    const previewUrl = file.type.startsWith('image/') ? URL.createObjectURL(file) : null
    const error = checkFileSize(file)
    const mimecat = getMimeCategory(file.type)
    const defaultCat = mimecat === 'image' ? 'post_cover' : 'attachment'
    const title = file.name.replace(/\.[^.]+$/, '')
    stagedFiles.value.push({
      id, file, previewUrl,
      title,
      altText: file.type.startsWith('image/') ? title : '',
      category: defaultCat,
      error,
    })
  }
}

watch(batchCategory, (val) => {
  if (!val) return
  stagedFiles.value.forEach(sf => { sf.category = val })
  nextTick(() => { batchCategory.value = '' })
})

const applyAltFromTitle = () => {
  stagedFiles.value.forEach(sf => {
    if (sf.file.type.startsWith('image/')) sf.altText = sf.title
  })
}

const removeStagedFile = (id: string) => {
  const sf = stagedFiles.value.find(f => f.id === id)
  if (sf?.previewUrl) URL.revokeObjectURL(sf.previewUrl)
  stagedFiles.value = stagedFiles.value.filter(f => f.id !== id)
}

// ── Modal open/close ──────────────────────────────────────────────────────
const onModalUpdate = (open: boolean) => {
  if (!open && !isUploading.value) {
    // User closed modal manually — clean up staged files
    stagedFiles.value.forEach(sf => { if (sf.previewUrl) URL.revokeObjectURL(sf.previewUrl) })
    stagedFiles.value = []
    uploadingFiles.value = []
    activeTab.value = 'upload'
    urlImport.value = { url: '', title: '', altText: '', category: '' }
  }
  emit('update:open', open)
}

// ── Upload ────────────────────────────────────────────────────────────────
const startUpload = async () => {
  const valid = stagedFiles.value.filter(f => !f.error)
  if (!valid.length) return
  isUploading.value = true

  uploadingFiles.value = valid.map(sf => ({
    id: sf.id,
    name: sf.file.name,
    progress: 0,
    previewUrl: sf.previewUrl,
    error: null as string | null,
  }))

  try {
    for (const sf of valid) {
      const entry = uploadingFiles.value.find(f => f.id === sf.id)
      if (!entry) continue
      entry.progress = 30
      const result = await mediaStore.uploadMedia(sf.file, {
        title: sf.title || sf.file.name,
        alt_text: sf.altText || undefined,
        category: (sf.category as any) || undefined,
      })
      entry.progress = result ? 100 : 0
      if (!result) entry.error = t('admin.media.upload_failed')
    }
  } finally {
    const failed = uploadingFiles.value.filter(f => f.error)
    const succeeded = uploadingFiles.value.filter(f => f.progress === 100)

    isUploading.value = false
    stagedFiles.value.forEach(sf => { if (sf.previewUrl) URL.revokeObjectURL(sf.previewUrl) })
    stagedFiles.value = []
    uploadingFiles.value = []

    await nextTick()

    // Close modal and notify parent
    emit('update:open', false)
    emit('uploaded')

    if (failed.length === 0) {
      toast.add({
        title: t('admin.media.upload_complete'),
        color: 'success',
        icon: 'i-tabler-circle-check',
      })
    } else if (succeeded.length > 0) {
      toast.add({
        title: `${succeeded.length} / ${valid.length}`,
        description: failed.map(f => f.name).join(', '),
        color: 'warning',
        icon: 'i-tabler-alert-triangle',
      })
    } else {
      toast.add({
        title: t('admin.media.upload_failed'),
        description: failed.length === 1 ? failed[0]?.name : failed.map(f => f.name).join(', '),
        color: 'error',
        icon: 'i-tabler-alert-circle',
      })
    }
  }
}
</script>
