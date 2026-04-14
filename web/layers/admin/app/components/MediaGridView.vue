<script setup lang="ts">
import type { MediaResponse } from "~/types/api/media";

const props = defineProps<{
  medias: MediaResponse[]
  selectedIds: number[]
  variantsCache: Map<number, Record<string, string>>
}>()

const emit = defineEmits<{
  (e: 'view', media: MediaResponse): void
  (e: 'copy-url', media: MediaResponse): void
  (e: 'delete', media: MediaResponse): void
  (e: 'toggle-select', id: number): void
}>()

const getMediaType = (mimeType: string) => {
  if (mimeType.startsWith("image/")) return "image"
  if (mimeType.startsWith("video/")) return "video"
  if (mimeType.startsWith("audio/")) return "audio"
  return "other"
}

const getThumbUrl = (media: MediaResponse, prefer = 'thumbnail') => {
  const v = props.variantsCache.get(media.id)
  if (v) {
    if (v[prefer]) return v[prefer]
    for (const k of ['thumbnail', 'cover', 'content']) {
      if (k !== prefer && v[k]) return v[k]
    }
  }
  return media.cdn_url
}

const formatFileSize = (bytes: number) => {
  if (!bytes) return "0 B"
  const k = 1024
  const sizes = ["B", "KB", "MB", "GB"]
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${Math.round((bytes / Math.pow(k, i)) * 100) / 100} ${sizes[i]}`
}
</script>

<template>
  <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 2xl:grid-cols-8 gap-4">
    <div
      v-for="media in medias"
      :key="media.id"
      class="group relative bg-default border border-default rounded-md overflow-hidden hover:shadow-md transition-all cursor-pointer"
      @click="emit('view', media)">
      <div class="absolute top-2 left-2 z-10">
        <UCheckbox
          :model-value="selectedIds.includes(media.id)"
          @click.stop="emit('toggle-select', media.id)" />
      </div>
      <div v-if="media.storage_type === 5" class="absolute top-2 right-2 z-10">
        <span class="text-[10px] font-semibold px-1.5 py-0.5 rounded bg-warning/90 text-warning-foreground leading-none">
          {{ $t('admin.media.badge_external') }}
        </span>
      </div>
      <div class="aspect-square bg-muted relative">
        <img
          v-if="getMediaType(media.mime_type) === 'image'"
          :src="getThumbUrl(media, 'thumbnail')"
          :alt="media.alt_text || media.filename"
          loading="lazy"
          class="w-full h-full object-cover" />
        <div v-else-if="getMediaType(media.mime_type) === 'video'" class="w-full h-full flex items-center justify-center">
          <UIcon name="i-tabler-player-play" class="size-12 text-muted" />
        </div>
        <div v-else-if="getMediaType(media.mime_type) === 'audio'" class="w-full h-full flex items-center justify-center">
          <UIcon name="i-tabler-music" class="size-12 text-muted" />
        </div>
        <div v-else class="w-full h-full flex items-center justify-center">
          <UIcon name="i-tabler-file" class="size-12 text-muted" />
        </div>
      </div>
      <div class="p-3">
        <div class="text-sm font-medium text-highlighted truncate" :title="media.filename">{{ media.filename }}</div>
        <div class="text-xs text-muted mt-1">{{ formatFileSize(media.file_size) }}</div>
      </div>
      <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center gap-2">
        <UButton color="neutral" variant="solid" icon="i-tabler-eye" square size="xs" @click.stop="emit('view', media)" />
        <UButton color="neutral" variant="solid" icon="i-tabler-copy" square size="xs" @click.stop="emit('copy-url', media)" />
        <UButton color="error" variant="solid" icon="i-tabler-trash" square size="xs" @click.stop="emit('delete', media)" />
      </div>
    </div>
  </div>
</template>
