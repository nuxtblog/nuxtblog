<script setup lang="ts">
import { useCurrentPost } from '~/composables/useCurrentPost'

import type { WidgetConfig } from '~/composables/useWidgetConfig'
const props = defineProps<{ config: WidgetConfig }>()
const { t } = useI18n()
const title = computed(() => props.config.title || t('site.widget.downloads.default_title'))

const { downloads } = useCurrentPost()

const extIcon = (url: string) => {
  const ext = url.split('.').pop()?.toLowerCase() ?? ''
  if (ext === 'pdf') return 'i-tabler-file-type-pdf'
  if (['zip', 'rar', '7z', 'tar', 'gz'].includes(ext)) return 'i-tabler-file-zip'
  if (['doc', 'docx'].includes(ext)) return 'i-tabler-file-type-doc'
  if (['xls', 'xlsx'].includes(ext)) return 'i-tabler-file-type-xls'
  if (['ppt', 'pptx'].includes(ext)) return 'i-tabler-file-type-ppt'
  if (['jpg', 'jpeg', 'png', 'gif', 'webp', 'svg'].includes(ext)) return 'i-tabler-photo'
  if (['mp4', 'avi', 'mov', 'mkv'].includes(ext)) return 'i-tabler-video'
  if (['mp3', 'wav', 'flac', 'aac'].includes(ext)) return 'i-tabler-music'
  return 'i-tabler-file-download'
}

const extColor = (url: string) => {
  const ext = url.split('.').pop()?.toLowerCase() ?? ''
  if (ext === 'pdf') return 'text-red-500 bg-red-50 dark:bg-red-950/30'
  if (['zip', 'rar', '7z'].includes(ext)) return 'text-yellow-500 bg-yellow-50 dark:bg-yellow-950/30'
  if (['doc', 'docx'].includes(ext)) return 'text-blue-500 bg-blue-50 dark:bg-blue-950/30'
  if (['xls', 'xlsx'].includes(ext)) return 'text-green-500 bg-green-50 dark:bg-green-950/30'
  if (['ppt', 'pptx'].includes(ext)) return 'text-orange-500 bg-orange-50 dark:bg-orange-950/30'
  return 'text-primary bg-primary/10'
}
</script>

<template>
  <UCard v-if="downloads.length > 0" :ui="{ body: 'p-0 sm:p-0' }">
    <template #header>
      <div class="flex items-center gap-2">
        <UIcon name="i-tabler-download" class="size-4 text-primary" />
        <h3 class="text-sm font-semibold text-highlighted">{{ title }}</h3>
        <UBadge :label="String(downloads.length)" color="primary" variant="subtle" size="xs" class="ml-auto" />
      </div>
    </template>
    <div class="divide-y divide-default">
      <a
        v-for="(item, i) in downloads"
        :key="i"
        :href="item.url"
        download
        target="_blank"
        rel="noopener"
        class="flex items-center gap-3 px-4 py-3 hover:bg-muted transition-colors group"
      >
        <div
          class="size-9 rounded-md flex items-center justify-center shrink-0"
          :class="extColor(item.url)"
        >
          <UIcon :name="extIcon(item.url)" class="size-5" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-highlighted truncate group-hover:text-primary transition-colors">
            {{ item.name }}
          </p>
          <p v-if="item.size || item.desc" class="text-xs text-muted mt-0.5 truncate">
            {{ item.size || item.desc }}
          </p>
        </div>
        <UIcon
          name="i-tabler-download"
          class="size-4 text-muted group-hover:text-primary transition-colors shrink-0"
        />
      </a>
    </div>
  </UCard>
</template>
