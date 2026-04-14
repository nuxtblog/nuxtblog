<template>
  <UModal
    :open="localOpen"
    title="选择媒体"
    :ui="{ content: 'max-w-3xl' }"
    @update:open="onModalUpdate">
    <template #content>
      <div class="p-6 space-y-4">
        <!-- Tab 切换 -->
        <div class="flex border-b border-default -mt-1 mb-2">
          <button
            class="px-4 py-2 text-sm font-medium transition-colors"
            :class="activeTab === 'browse' ? 'text-primary border-b-2 border-primary' : 'text-muted hover:text-highlighted'"
            @click="activeTab = 'browse'">
            媒体库
          </button>
          <button
            class="px-4 py-2 text-sm font-medium transition-colors"
            :class="activeTab === 'upload' ? 'text-primary border-b-2 border-primary' : 'text-muted hover:text-highlighted'"
            @click="activeTab = 'upload'">
            上传
          </button>
        </div>

        <!-- 媒体库 Tab -->
        <div v-if="activeTab === 'browse'" class="space-y-3">
          <!-- 分类筛选（多分类时显示下拉） -->
          <div v-if="categoryList.length > 1" class="flex items-center gap-2">
            <span class="text-xs text-muted shrink-0">分类:</span>
            <USelect
              v-model="browseCategory"
              :items="categorySelectItems"
              size="sm"
              class="w-40" />
          </div>

          <!-- 加载中 -->
          <div v-if="browseLoading" class="flex items-center justify-center py-12">
            <span class="text-muted text-sm">加载中...</span>
          </div>

          <!-- 空状态 -->
          <div v-else-if="mediaItems.length === 0" class="flex flex-col items-center justify-center py-12">
            <UIcon name="i-tabler-photo-off" class="size-10 text-muted mb-2" />
            <span class="text-muted text-sm">暂无媒体文件</span>
          </div>

          <!-- 网格 -->
          <div v-else class="grid grid-cols-4 sm:grid-cols-5 md:grid-cols-6 gap-2 max-h-[50vh] overflow-y-auto pr-1">
            <button
              v-for="item in mediaItems"
              :key="item.id"
              class="relative aspect-square rounded-md overflow-hidden border-2 transition-all cursor-pointer group"
              :class="isSelected(item.id)
                ? 'border-primary ring-2 ring-primary/30'
                : 'border-default hover:border-primary/50'"
              @click="toggleSelect(item)">
              <!-- 图片预览 -->
              <img
                v-if="item.mime_type.startsWith('image/')"
                :src="item.cdn_url"
                :alt="item.title"
                class="w-full h-full object-cover" />
              <!-- 非图片文件图标 -->
              <div v-else class="w-full h-full flex flex-col items-center justify-center bg-elevated gap-1">
                <UIcon :name="getFileIcon(item.mime_type)" class="size-8 text-muted" />
                <span class="text-[10px] text-muted truncate max-w-full px-1">{{ item.title || item.filename }}</span>
              </div>
              <!-- 选中标记 -->
              <div
                v-if="isSelected(item.id)"
                class="absolute top-1 right-1 w-5 h-5 rounded-full bg-primary flex items-center justify-center">
                <UIcon name="i-tabler-check" class="size-3 text-white" />
              </div>
            </button>
          </div>

          <!-- 分页 -->
          <div v-if="browseTotal > browseSize" class="flex justify-center pt-2">
            <UPagination
              v-model="browsePage"
              :total="browseTotal"
              :items-per-page="browseSize"
              size="sm" />
          </div>

          <!-- 确认 -->
          <div class="flex justify-end gap-2 pt-2">
            <UButton color="neutral" variant="outline" @click="localOpen = false">取消</UButton>
            <UButton
              :disabled="selected.length === 0"
              @click="confirmSelection">
              确认 ({{ selected.length }})
            </UButton>
          </div>
        </div>

        <!-- 上传 Tab -->
        <div v-if="activeTab === 'upload'" class="space-y-4">
          <!-- 拖拽区 -->
          <div
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
              :multiple="multiple"
              :accept="accept || '*/*'"
              class="hidden"
              @change="handleFileSelect" />
            <UIcon name="i-tabler-upload" class="size-8 mx-auto text-muted mb-3" />
            <p class="text-sm text-muted mb-1">拖拽文件到此处或点击选择</p>
            <p v-if="accept" class="text-xs text-muted">接受类型: {{ accept }}</p>
          </div>

          <!-- 上传中的文件列表 -->
          <div v-if="uploadingFiles.length > 0" class="space-y-2 max-h-[40vh] overflow-y-auto pr-1">
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
            <UButton color="neutral" variant="outline" @click="localOpen = false">取消</UButton>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { MediaResponse } from '~/types/api/media'

interface MediaPickResult {
  id: number
  url: string
  title: string
  mime_type: string
}

interface Props {
  open: boolean
  categories?: string[]
  defaultCategory?: string
  accept?: string
  multiple?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  categories: () => [],
  defaultCategory: '',
  accept: '',
  multiple: false,
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  'select': [result: MediaPickResult | MediaPickResult[]]
}>()

const localOpen = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val),
})

const mediaApi = useMediaApi()
const toast = useToast()

// ── Tab state ──────────────────────────────────────────────────────────────
const activeTab = ref<'browse' | 'upload'>('browse')

// ── Browse state ───────────────────────────────────────────────────────────
const categoryList = computed(() => props.categories?.length ? props.categories : [])
const browseCategory = ref('')
const browsePage = ref(1)
const browseSize = 24
const browseTotal = ref(0)
const browseLoading = ref(false)
const mediaItems = ref<MediaResponse[]>([])
const selected = ref<MediaPickResult[]>([])

const categorySelectItems = computed(() =>
  categoryList.value.map(c => ({ label: c, value: c }))
)

const effectiveCategory = computed(() => {
  if (categoryList.value.length === 1) return categoryList.value[0]
  return browseCategory.value || ''
})

const uploadCategory = computed(() =>
  props.defaultCategory || (categoryList.value.length ? categoryList.value[0] : '')
)

// Fetch media list
const fetchMedia = async () => {
  browseLoading.value = true
  try {
    const query: Record<string, any> = {
      page: browsePage.value,
      size: browseSize,
    }
    if (effectiveCategory.value) query.category = effectiveCategory.value
    if (props.accept && props.accept !== '*/*') {
      // Convert 'image/*' to mime_type filter — API supports prefix match
      // For exact types like 'image/png', pass directly
      // For wildcards like 'image/*', pass 'image/' prefix
      const mime = props.accept.replace('/*', '/')
      if (mime !== '*/') query.mime_type = mime
    }
    const res = await mediaApi.list(query)
    mediaItems.value = res.list ?? []
    browseTotal.value = res.total ?? 0
  } catch (err) {
    toast.add({ title: '加载媒体失败', color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    browseLoading.value = false
  }
}

// Re-fetch when category/page changes
watch([effectiveCategory, browsePage], () => fetchMedia())

// Fetch on modal open
watch(() => props.open, (val) => {
  if (val) {
    selected.value = []
    browsePage.value = 1
    if (categoryList.value.length > 1) browseCategory.value = categoryList.value[0]
    activeTab.value = 'browse'
    fetchMedia()
  }
})

const isSelected = (id: number) => selected.value.some(s => s.id === id)

const toggleSelect = (item: MediaResponse) => {
  const pick: MediaPickResult = {
    id: item.id,
    url: item.cdn_url,
    title: item.title || item.filename,
    mime_type: item.mime_type,
  }

  if (props.multiple) {
    const idx = selected.value.findIndex(s => s.id === item.id)
    if (idx >= 0) {
      selected.value.splice(idx, 1)
    } else {
      selected.value.push(pick)
    }
  } else {
    // Single mode: toggle or replace
    if (selected.value.length === 1 && selected.value[0].id === item.id) {
      selected.value = []
    } else {
      selected.value = [pick]
    }
  }
}

const confirmSelection = () => {
  if (selected.value.length === 0) return
  if (props.multiple) {
    emit('select', [...selected.value])
  } else {
    emit('select', selected.value[0])
  }
  localOpen.value = false
}

const getFileIcon = (mime: string): string => {
  if (mime.startsWith('video/')) return 'i-tabler-video'
  if (mime.startsWith('audio/')) return 'i-tabler-music'
  if (mime.includes('pdf')) return 'i-tabler-file-type-pdf'
  if (mime.includes('zip') || mime.includes('rar') || mime.includes('7z')) return 'i-tabler-file-zip'
  return 'i-tabler-file'
}

// ── Upload state ───────────────────────────────────────────────────────────
const isDragging = ref(false)
const fileInput = ref<HTMLInputElement>()

interface UploadingFile {
  id: string
  name: string
  progress: number
  previewUrl: string | null
  error: string | null
  result?: MediaPickResult
}
const uploadingFiles = ref<UploadingFile[]>([])

const handleFileSelect = (event: Event) => {
  const files = (event.target as HTMLInputElement).files
  if (files) uploadFiles(Array.from(files))
  if (fileInput.value) fileInput.value.value = ''
}

const handleDrop = (event: DragEvent) => {
  isDragging.value = false
  const files = event.dataTransfer?.files
  if (files) uploadFiles(Array.from(files))
}

const uploadFiles = async (files: File[]) => {
  const entries: UploadingFile[] = files.map(f => ({
    id: `${Date.now()}-${Math.random()}`,
    name: f.name,
    progress: 0,
    previewUrl: f.type.startsWith('image/') ? URL.createObjectURL(f) : null,
    error: null,
  }))
  uploadingFiles.value.push(...entries)

  const uploaded: MediaPickResult[] = []

  for (let i = 0; i < files.length; i++) {
    const file = files[i]
    const entry = entries[i]
    entry.progress = 30
    try {
      const res = await mediaApi.upload(file, {
        title: file.name.replace(/\.[^.]+$/, ''),
        category: uploadCategory.value || undefined,
      })
      entry.progress = 100
      const pick: MediaPickResult = {
        id: res.id,
        url: res.cdn_url,
        title: res.title || res.filename,
        mime_type: res.mime_type,
      }
      entry.result = pick
      uploaded.push(pick)
    } catch (err: any) {
      entry.error = err?.message || '上传失败'
    }
  }

  // Auto-select uploaded files
  if (uploaded.length > 0) {
    if (props.multiple) {
      emit('select', uploaded)
    } else {
      emit('select', uploaded[0])
    }
    toast.add({ title: `上传成功 (${uploaded.length})`, color: 'success', icon: 'i-tabler-circle-check' })
    // Close modal after short delay
    setTimeout(() => {
      localOpen.value = false
    }, 300)
  } else {
    // Collect error messages from failed uploads
    const errors = entries.filter(e => e.error).map(e => `${e.name}: ${e.error}`)
    toast.add({
      title: '上传失败',
      description: errors.join('\n'),
      color: 'error',
      icon: 'i-tabler-alert-circle',
    })
  }
}

// ── Cleanup ────────────────────────────────────────────────────────────────
const onModalUpdate = (open: boolean) => {
  if (!open) {
    uploadingFiles.value.forEach(f => { if (f.previewUrl) URL.revokeObjectURL(f.previewUrl) })
    uploadingFiles.value = []
    selected.value = []
  }
  emit('update:open', open)
}
</script>
