<script setup lang="ts">
interface DownloadItem {
  name: string
  url: string
  size?: string
  desc?: string
}

defineProps<{
  title: string
  collapsed: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle'): void
}>()

const postDownloads = defineModel<DownloadItem[]>('downloads', { required: true })

const { t } = useI18n()

const addDownload = () => postDownloads.value.push({ name: '', url: '', size: '', desc: '' })
</script>

<template>
  <SidebarCard :title="title" :collapsed="collapsed" @toggle="emit('toggle')">
    <template #header-actions>
      <UButton color="primary" variant="link" size="xs" icon="i-tabler-plus" @click.stop="addDownload">
        {{ t('common.add') }}
      </UButton>
    </template>
    <div class="space-y-2">
      <div
        v-for="(dl, i) in postDownloads"
        :key="i"
        class="rounded-md border border-default p-3 space-y-2">
        <div class="flex items-center gap-2">
          <UInput v-model="dl.name" :placeholder="t('admin.posts.editor.download_name')" size="xs" class="flex-1" />
          <UButton icon="i-tabler-trash" color="error" variant="ghost" size="xs" square @click="postDownloads.splice(i, 1)" />
        </div>
        <UInput v-model="dl.url" :placeholder="t('admin.posts.editor.download_url')" size="xs" class="w-full" />
        <div class="flex gap-2">
          <UInput v-model="dl.size" :placeholder="t('admin.posts.editor.download_size')" size="xs" class="flex-1" />
          <UInput v-model="dl.desc" :placeholder="t('admin.posts.editor.download_desc')" size="xs" class="flex-1" />
        </div>
      </div>
      <p v-if="postDownloads.length === 0" class="text-xs text-muted text-center py-2">
        {{ t('admin.posts.editor.no_downloads') }}
      </p>
    </div>
  </SidebarCard>
</template>
