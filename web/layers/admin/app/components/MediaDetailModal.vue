<template>
  <UModal
    :open="open"
    :title="$t('admin.media.file_detail')"
    :ui="{ content: 'max-w-4xl' }"
    @update:open="emit('update:open', $event)">
    <template #content>
      <div v-if="media" class="p-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- 预览区 -->
          <div class="space-y-3">
            <!-- 图片：展示缩略图，支持点击查看原图 -->
            <template v-if="getMediaType(media.mime_type) === 'image'">
              <div
                class="bg-muted rounded-md overflow-hidden cursor-zoom-in relative group"
                @click="emit('open-lightbox', media)">
                <img
                  :src="getThumbUrl('cover')"
                  :alt="media.alt_text"
                  class="w-full h-auto max-h-72 object-contain" />
                <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
                  <UIcon name="i-tabler-zoom-in" class="size-8 text-white" />
                </div>
              </div>
              <div class="flex gap-2">
                <UButton
                  color="neutral" variant="outline" size="sm" icon="i-tabler-zoom-in"
                  class="flex-1" @click="emit('open-lightbox', media)">
                  {{ $t('common.preview') }}
                </UButton>
                <UButton
                  :to="media.cdn_url" target="_blank"
                  color="neutral" variant="outline" size="sm" icon="i-tabler-external-link"
                  class="flex-1">
                  {{ $t('common.view') }}
                </UButton>
              </div>
            </template>

            <!-- 视频 -->
            <template v-else-if="getMediaType(media.mime_type) === 'video'">
              <div class="bg-muted rounded-md overflow-hidden">
                <video :src="media.cdn_url" controls class="w-full h-auto" />
              </div>
              <UButton
                :to="media.cdn_url" target="_blank"
                color="neutral" variant="outline" size="sm" icon="i-tabler-external-link"
                class="w-full">
                {{ $t('common.view') }}
              </UButton>
            </template>

            <!-- 音频 -->
            <template v-else-if="getMediaType(media.mime_type) === 'audio'">
              <div class="bg-muted rounded-md p-4">
                <audio :src="media.cdn_url" controls class="w-full" />
              </div>
              <UButton
                :to="media.cdn_url" target="_blank"
                color="neutral" variant="outline" size="sm" icon="i-tabler-external-link"
                class="w-full">
                {{ $t('common.view') }}
              </UButton>
            </template>

            <!-- 其他文件 -->
            <template v-else>
              <div class="bg-muted rounded-md p-8 flex flex-col items-center justify-center gap-3">
                <UIcon name="i-tabler-file" class="size-16 text-muted" />
                <UButton
                  :to="media.cdn_url" target="_blank"
                  color="neutral" variant="outline" size="sm" icon="i-tabler-download">
                  {{ $t('common.view') }}
                </UButton>
              </div>
            </template>

            <!-- URL 复制行 -->
            <div class="flex gap-2">
              <UInput :model-value="media.cdn_url" readonly class="flex-1" size="sm" />
              <UButton color="neutral" variant="ghost" icon="i-tabler-copy" square size="sm" @click="copyUrl(media)" />
            </div>
          </div>

          <!-- 编辑信息 -->
          <div class="space-y-4">
            <UFormField :label="$t('common.title')">
              <UInput v-model="media.title" :placeholder="$t('admin.media.file_name')" class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.media.alt_text')">
              <UInput v-model="media.alt_text" :placeholder="$t('admin.media.alt_placeholder')" class="w-full" />
            </UFormField>

            <USeparator />

            <div class="grid grid-cols-2 gap-4 text-sm">
              <div>
                <div class="text-xs text-muted mb-1">{{ $t('admin.media.file_name') }}</div>
                <div class="font-medium text-highlighted break-all text-xs">{{ media.filename }}</div>
              </div>
              <div>
                <div class="text-xs text-muted mb-1">{{ $t('admin.media.file_type') }}</div>
                <div class="font-medium text-highlighted">{{ media.mime_type }}</div>
              </div>
              <div>
                <div class="text-xs text-muted mb-1">{{ $t('admin.media.file_size') }}</div>
                <div class="font-medium text-highlighted">{{ formatFileSize(media.file_size) }}</div>
              </div>
              <div v-if="media.width">
                <div class="text-xs text-muted mb-1">{{ $t('admin.media.dimensions') }}</div>
                <div class="font-medium text-highlighted">{{ media.width }} × {{ media.height }}</div>
              </div>
              <div>
                <div class="text-xs text-muted mb-1">{{ $t('admin.media.upload_time') }}</div>
                <div class="font-medium text-highlighted">{{ formatDate(media.created_at) }}</div>
              </div>
              <div>
                <div class="text-xs text-muted mb-1">{{ $t('admin.media.category') }}</div>
                <div class="font-medium text-highlighted">{{ getCategoryLabel(media.category) }}</div>
              </div>
            </div>

            <!-- 外部媒体标记 -->
            <div v-if="media.storage_type === 5" class="flex items-center gap-2 p-2 bg-warning/10 rounded-md">
              <UIcon name="i-tabler-link" class="size-4 text-warning shrink-0" />
              <span class="text-xs text-warning font-medium">{{ $t('admin.media.badge_external') }}</span>
            </div>

            <div class="flex gap-3 pt-2 flex-wrap">
              <UButton color="primary" :loading="updating" class="flex-1" @click="handleUpdateMedia">{{ $t('common.save') }}</UButton>
              <UButton
                v-if="media.storage_type === 5"
                color="neutral" variant="outline"
                icon="i-tabler-download"
                :loading="localizing"
                @click="handleLocalize">
                {{ $t('admin.media.localize') }}
              </UButton>
              <UButton color="error" variant="soft" @click="emit('delete', media)">{{ $t('common.delete') }}</UButton>
            </div>
          </div>
        </div>
      </div>
    </template>
  </UModal>
</template>

<script setup lang="ts">
import type { MediaResponse } from '~/types/api/media'

interface Props {
  open: boolean
  media: MediaResponse | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  'delete': [media: MediaResponse]
  'saved': []
  'open-lightbox': [media: MediaResponse]
  'localized': [media: MediaResponse]
}>()

const mediaStore = useMediaStore()
const mediaApi = useMediaApi()
const toast = useToast()
const { t } = useI18n()
const { getCategoryLabel } = useMediaCategories()

// ── State ──────────────────────────────────────────────────────────────────
const updating = ref(false)
const localizing = ref(false)

// ── Variants / thumbnail ──────────────────────────────────────────────────
const variants = computed(() => {
  if (!props.media?.variants) return {} as Record<string, string>
  try { return JSON.parse(props.media.variants) as Record<string, string> } catch { return {} }
})

const getThumbUrl = (prefer: 'thumbnail' | 'cover' | 'content' = 'thumbnail'): string => {
  const v = variants.value
  if (v[prefer]) return v[prefer]!
  for (const k of ['thumbnail', 'cover', 'content'] as const) {
    if (k !== prefer && v[k]) return v[k]!
  }
  return props.media?.cdn_url ?? ''
}

// ── Helpers ────────────────────────────────────────────────────────────────
const getMediaType = (mimeType: string): string => {
  if (mimeType.startsWith('image/')) return 'image'
  if (mimeType.startsWith('video/')) return 'video'
  if (mimeType.startsWith('audio/')) return 'audio'
  return 'other'
}

// ── Formatters ────────────────────────────────────────────────────────────
const formatFileSize = (bytes: number): string => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${Math.round((bytes / Math.pow(k, i)) * 100) / 100} ${sizes[i]}`
}

const formatDate = (s: string): string => {
  if (!s) return '—'
  const d = new Date(s)
  const days = Math.floor((Date.now() - d.getTime()) / 86400000)
  if (days === 0) return t('common.today')
  if (days === 1) return t('common.yesterday')
  if (days < 7) return t('common.days_ago', { n: days })
  if (days < 30) return t('common.weeks_ago', { n: Math.floor(days / 7) })
  return d.toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' })
}

// ── Copy URL ──────────────────────────────────────────────────────────────
const copyUrl = async (media: MediaResponse) => {
  try {
    await navigator.clipboard.writeText(media.cdn_url)
    toast.add({ title: t('admin.media.url_copied'), icon: 'i-tabler-copy', color: 'success' })
  } catch {
    const input = document.createElement('input')
    input.value = media.cdn_url
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
  }
}

// ── Localize ──────────────────────────────────────────────────────────────
const handleLocalize = async () => {
  if (!props.media) return
  localizing.value = true
  try {
    const updated = await mediaApi.localize(props.media.id)
    toast.add({ title: t('admin.media.localize_success'), icon: 'i-tabler-circle-check', color: 'success' })
    emit('localized', updated)
    emit('saved')
  } catch (err) {
    toast.add({ title: err instanceof Error ? err.message : t('admin.media.localize_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    localizing.value = false
  }
}

// ── Update ────────────────────────────────────────────────────────────────
const handleUpdateMedia = async () => {
  if (!props.media) return
  updating.value = true
  try {
    await mediaStore.updateMedia(props.media.id, {
      alt_text: props.media.alt_text,
      title: props.media.title,
    })
    toast.add({ title: t('admin.media.saved'), icon: 'i-tabler-circle-check', color: 'success' })
    emit('saved')
  } catch (err) {
    toast.add({ title: err instanceof Error ? err.message : t('admin.media.save_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    updating.value = false
  }
}
</script>
