<script setup lang="ts">
import type { MediaResponse } from "~/types/api/media";

const props = defineProps<{
  medias: MediaResponse[]
  selectedIds: number[]
  variantsCache: Map<number, Record<string, string>>
  actions: (media: MediaResponse) => any[][]
}>()

const emit = defineEmits<{
  (e: 'toggle-select', id: number): void
}>()

const { t } = useI18n()

const getMediaType = (mimeType: string) => {
  if (mimeType.startsWith("image/")) return "image"
  return "other"
}

const getThumbUrl = (media: MediaResponse) => {
  const v = props.variantsCache.get(media.id)
  if (v) {
    for (const k of ['thumbnail', 'cover', 'content']) {
      if (v[k]) return v[k]
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

const formatDate = (s: string) => {
  if (!s) return "—"
  const d = new Date(s)
  const days = Math.floor((Date.now() - d.getTime()) / 86400000)
  if (days === 0) return t('common.today')
  if (days === 1) return t('common.yesterday')
  if (days < 7) return t('common.days_ago', { n: days })
  if (days < 30) return t('common.weeks_ago', { n: Math.floor(days / 7) })
  return d.toLocaleDateString(undefined, { year: "numeric", month: "long", day: "numeric" })
}
</script>

<template>
  <div class="space-y-3">
    <div
      v-for="media in medias"
      :key="media.id"
      class="flex items-center group gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all">
      <UCheckbox :model-value="selectedIds.includes(media.id)" @update:model-value="emit('toggle-select', media.id)" />

      <div class="h-10 w-10 rounded overflow-hidden bg-muted shrink-0">
        <img
          v-if="getMediaType(media.mime_type) === 'image'"
          :src="getThumbUrl(media)"
          :alt="media.alt_text"
          loading="lazy"
          class="w-full h-full object-cover" />
        <div v-else class="w-full h-full flex items-center justify-center">
          <UIcon name="i-tabler-file" class="size-5 text-muted" />
        </div>
      </div>

      <div class="flex-1 min-w-0">
        <div class="font-medium text-highlighted text-sm truncate">{{ media.filename }}</div>
        <div class="text-xs text-muted mt-0.5">
          {{ media.mime_type }} · {{ formatFileSize(media.file_size) }}<template v-if="media.width"> · {{ media.width }} × {{ media.height }}</template>
        </div>
        <div class="text-xs text-muted mt-0.5">{{ formatDate(media.created_at) }}</div>
      </div>

      <div class="flex items-center gap-3 shrink-0">
        <UBadge
          v-if="media.storage_type === 5"
          :label="$t('admin.media.badge_external')"
          color="warning" variant="soft" size="sm" />
        <UDropdownMenu :items="actions(media)" :popper="{ placement: 'bottom-end' }">
          <UButton
            color="neutral" variant="ghost" icon="i-tabler-dots-vertical"
            square size="xs"
            class="opacity-0 group-hover:opacity-100 transition-opacity" />
        </UDropdownMenu>
      </div>
    </div>
  </div>
</template>
