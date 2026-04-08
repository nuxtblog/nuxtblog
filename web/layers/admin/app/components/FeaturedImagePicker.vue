<template>
  <!-- 已选图片 -->
  <div v-if="imgUrl" class="relative group rounded-md overflow-hidden">
    <img :src="imgUrl" alt="featured" class="w-full h-40 object-cover" />
    <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
      <UButton color="neutral" variant="solid" icon="i-tabler-pencil" size="xs" square @click="open = true" />
      <UButton color="error" variant="solid" icon="i-tabler-x" size="xs" square @click="clear" />
    </div>
  </div>

  <!-- 占位区 -->
  <div
    v-else
    class="border-2 border-dashed border-default rounded-md p-6 text-center flex flex-col items-center cursor-pointer hover:border-primary/50 transition-colors"
    @click="open = true">
    <UIcon name="i-tabler-photo" class="size-6 mb-2 text-muted" />
    <UButton color="primary" variant="link" size="sm" @click.stop="open = true">
      {{ t('admin.posts.editor.select_image') }}
    </UButton>
  </div>

  <!-- 选择器 Modal -->
  <UModal v-model:open="open" :title="t('admin.posts.editor.featured_image')" :ui="{ content: 'max-w-3xl' }">
    <template #content>
      <div class="p-4">
        <!-- Tabs -->
        <div class="flex gap-1 p-1 bg-default rounded-md ring-1 ring-default w-fit mb-4">
          <UButton
            v-for="tab in tabs" :key="tab.value"
            :variant="activeTab === tab.value ? 'solid' : 'ghost'"
            :color="activeTab === tab.value ? 'primary' : 'neutral'"
            size="sm"
            :class="activeTab === tab.value ? 'shadow' : 'text-muted'"
            @click="activeTab = tab.value">
            {{ tab.label }}
          </UButton>
        </div>

        <!-- ── Tab: Library ── -->
        <div v-if="activeTab === 'library'">
          <!-- Search -->
          <UInput
            v-model="libSearch"
            :placeholder="t('admin.media.search_placeholder')"
            leading-icon="i-tabler-search"
            size="sm"
            class="w-full mb-3" />

          <!-- Loading -->
          <div v-if="libLoading" class="grid grid-cols-4 sm:grid-cols-5 gap-3">
            <USkeleton v-for="i in 10" :key="i" class="aspect-square rounded-md" />
          </div>

          <!-- Grid -->
          <div v-else-if="filteredLib.length > 0" class="grid grid-cols-4 sm:grid-cols-5 gap-3 max-h-80 overflow-y-auto pr-1">
            <div
              v-for="item in filteredLib"
              :key="item.id"
              class="relative aspect-square rounded-md overflow-hidden cursor-pointer group ring-2 transition-all"
              :class="libSelected === item.id ? 'ring-primary' : 'ring-transparent hover:ring-primary/50'"
              @click="libSelected = item.id; libSelectedUrl = item.cdn_url">
              <img :src="item.cdn_url" :alt="item.alt_text || item.filename" class="w-full h-full object-cover" loading="lazy" />
              <div v-if="libSelected === item.id" class="absolute inset-0 bg-primary/20 flex items-center justify-center">
                <UIcon name="i-tabler-check" class="size-5 text-white drop-shadow" />
              </div>
            </div>
          </div>

          <!-- Empty -->
          <div v-else class="py-10 text-center text-muted text-sm">
            <UIcon name="i-tabler-photo-off" class="size-10 mx-auto mb-2" />
            <p>{{ t('admin.media.no_files') }}</p>
          </div>

          <!-- Pagination -->
          <div v-if="libTotal > libPageSize" class="flex justify-center mt-3">
            <UPagination v-model:page="libPage" :total="libTotal" :items-per-page="libPageSize" size="sm" />
          </div>

          <div class="flex justify-end gap-2 mt-4">
            <UButton color="neutral" variant="soft" @click="open = false">{{ t('common.cancel') }}</UButton>
            <UButton color="primary" :disabled="!libSelected" @click="confirmLibrary">{{ t('common.confirm') }}</UButton>
          </div>
        </div>

        <!-- ── Tab: Upload ── -->
        <div v-else-if="activeTab === 'upload'">
          <div
            class="border-2 border-dashed rounded-md p-8 text-center transition-colors cursor-pointer"
            :class="isDragging ? 'border-primary bg-primary/5' : 'border-default hover:border-primary/50'"
            @drop.prevent="handleDrop"
            @dragover.prevent="isDragging = true"
            @dragleave="isDragging = false"
            @click="uploadInput?.click()">
            <input ref="uploadInput" type="file" accept="image/*" class="hidden" @change="handleFileSelect" />
            <UIcon name="i-tabler-upload" class="size-8 mx-auto text-muted mb-3" />
            <p class="text-sm text-muted mb-1">{{ t('admin.media.upload_area_title') }}</p>
            <p class="text-xs text-muted">{{ t('admin.posts.editor.upload_image_hint') }}</p>
          </div>

          <!-- Preview after file selected -->
          <div v-if="uploadPreviewUrl" class="mt-4 flex items-start gap-3">
            <img :src="uploadPreviewUrl" class="w-20 h-20 object-cover rounded-md shrink-0" />
            <div class="flex-1 min-w-0">
              <p class="text-sm text-highlighted truncate">{{ uploadFile?.name }}</p>
              <p class="text-xs text-muted mt-0.5">{{ formatFileSize(uploadFile?.size ?? 0) }}</p>

              <UFormField :label="t('admin.media.alt_text')" class="mt-2">
                <UInput v-model="uploadAlt" size="xs" :placeholder="t('admin.media.alt_placeholder')" class="w-full" />
              </UFormField>
            </div>
          </div>

          <div class="flex justify-end gap-2 mt-4">
            <UButton color="neutral" variant="soft" @click="open = false">{{ t('common.cancel') }}</UButton>
            <UButton color="primary" :loading="uploading" :disabled="!uploadFile" @click="confirmUpload">
              {{ t('admin.media.upload_file') }}
            </UButton>
          </div>
        </div>

        <!-- ── Tab: External URL ── -->
        <div v-else-if="activeTab === 'external'">
          <UFormField :label="t('admin.posts.editor.external_url_label')" class="mb-3">
            <UInput
              v-model="extUrl"
              placeholder="https://example.com/image.jpg"
              leading-icon="i-tabler-link"
              class="w-full"
              @input="extPreviewError = false" />
          </UFormField>

          <!-- Preview -->
          <div v-if="extUrl" class="mb-4">
            <p class="text-xs text-muted mb-2">{{ t('admin.posts.editor.preview') }}</p>
            <div class="w-full h-40 bg-elevated rounded-md overflow-hidden flex items-center justify-center">
              <img
                v-if="!extPreviewError"
                :src="extUrl"
                class="w-full h-full object-cover"
                @error="extPreviewError = true" />
              <div v-else class="text-center text-muted text-sm px-4">
                <UIcon name="i-tabler-photo-off" class="size-8 mx-auto mb-2" />
                <p>{{ t('admin.posts.editor.preview_failed') }}</p>
              </div>
            </div>
          </div>

          <UFormField :label="t('common.title')" class="mb-3">
            <UInput v-model="extTitle" :placeholder="t('admin.media.alt_placeholder')" class="w-full" />
          </UFormField>

          <div class="flex justify-end gap-2 mt-4">
            <UButton color="neutral" variant="soft" @click="open = false">{{ t('common.cancel') }}</UButton>
            <UButton color="primary" :loading="linking" :disabled="!extUrl || extPreviewError" @click="confirmExternal">
              {{ t('common.confirm') }}
            </UButton>
          </div>
        </div>

      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { MediaResponse } from '~/types/api/media'

const imgId  = defineModel<number | undefined>('imgId')
const imgUrl = defineModel<string>('imgUrl', { required: true })

const { t } = useI18n()
const toast = useToast()
const mediaApi = useMediaApi()

const open = ref(false)

// ── Tabs ──────────────────────────────────────────────────────────────────────
const activeTab = ref<'library' | 'upload' | 'external'>('library')
const tabs = computed(() => [
  { label: t('admin.posts.editor.tab_library'),  value: 'library'  },
  { label: t('admin.posts.editor.tab_upload'),   value: 'upload'   },
  { label: t('admin.posts.editor.tab_external'), value: 'external' },
])

// Reset state when modal opens
watch(open, (v) => {
  if (v) {
    activeTab.value = 'library'
    libSelected.value = undefined
    libSelectedUrl.value = ''
    uploadFile.value = null
    uploadPreviewUrl.value = ''
    uploadAlt.value = ''
    extUrl.value = ''
    extTitle.value = ''
    extPreviewError.value = false
    if (libImages.value.length === 0) loadLibrary()
  }
})

const clear = () => {
  imgId.value = undefined
  imgUrl.value = ''
}

// ── Library tab ───────────────────────────────────────────────────────────────
const libImages      = ref<MediaResponse[]>([])
const libLoading     = ref(false)
const libPage        = ref(1)
const libPageSize    = 20
const libTotal       = ref(0)
const libSearch      = ref('')
const libSelected    = ref<number | undefined>(undefined)
const libSelectedUrl = ref('')

const filteredLib = computed(() =>
  libSearch.value
    ? libImages.value.filter(m =>
        m.filename.toLowerCase().includes(libSearch.value.toLowerCase()) ||
        m.title.toLowerCase().includes(libSearch.value.toLowerCase())
      )
    : libImages.value
)

const loadLibrary = async () => {
  libLoading.value = true
  try {
    const res = await mediaApi.list({ page: libPage.value, size: libPageSize, mime_type: 'image' })
    libImages.value = res.list
    libTotal.value  = res.total
  } catch {}
  finally { libLoading.value = false }
}

watch(libPage, loadLibrary)

const confirmLibrary = () => {
  if (!libSelected.value || !libSelectedUrl.value) return
  imgId.value  = libSelected.value
  imgUrl.value = libSelectedUrl.value
  open.value = false
}

// ── Upload tab ────────────────────────────────────────────────────────────────
const uploadInput      = ref<HTMLInputElement | null>(null)
const uploadFile       = ref<File | null>(null)
const uploadPreviewUrl = ref('')
const uploadAlt        = ref('')
const uploading        = ref(false)
const isDragging       = ref(false)

const handleFileSelect = (e: Event) => {
  const file = (e.target as HTMLInputElement).files?.[0]
  if (!file) return
  setUploadFile(file)
}
const handleDrop = (e: DragEvent) => {
  isDragging.value = false
  const file = e.dataTransfer?.files?.[0]
  if (file && file.type.startsWith('image/')) setUploadFile(file)
}
const setUploadFile = (file: File) => {
  uploadFile.value = file
  uploadPreviewUrl.value = URL.createObjectURL(file)
}

const confirmUpload = async () => {
  if (!uploadFile.value) return
  uploading.value = true
  try {
    const media = await mediaApi.upload(uploadFile.value, {
      title:    uploadFile.value.name,
      alt_text: uploadAlt.value || undefined,
      category: 'post',
    })
    imgId.value  = media.id
    imgUrl.value = media.cdn_url
    open.value = false
    // Invalidate library cache so it shows on next open
    libImages.value = []
  } catch (err: any) {
    toast.add({ title: t('admin.posts.editor.upload_failed'), description: err?.message, color: 'error' })
  } finally {
    uploading.value = false
  }
}

// ── External URL tab ──────────────────────────────────────────────────────────
const extUrl          = ref('')
const extTitle        = ref('')
const extPreviewError = ref(false)
const linking         = ref(false)

const confirmExternal = async () => {
  if (!extUrl.value) return
  linking.value = true
  try {
    const media = await mediaApi.link({
      url:      extUrl.value,
      title:    extTitle.value || undefined,
      category: 'post',
    })
    imgId.value  = media.id
    imgUrl.value = media.cdn_url
    open.value = false
  } catch (err: any) {
    toast.add({ title: t('admin.posts.editor.link_failed'), description: err?.message, color: 'error' })
  } finally {
    linking.value = false
  }
}

// ── Helpers ───────────────────────────────────────────────────────────────────
const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1024 / 1024).toFixed(1)} MB`
}
</script>
