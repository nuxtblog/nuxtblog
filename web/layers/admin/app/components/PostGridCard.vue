<template>
  <div class="group flex flex-col rounded-md border border-default bg-default overflow-hidden hover:shadow-md transition-all">
    <!-- 封面 -->
    <div class="relative h-44 bg-muted shrink-0">
      <img
        :src="post.featured_img?.url || defaultCover"
        :alt="post.title"
        class="h-full w-full object-cover group-hover:scale-105 transition-transform duration-300"
        @error="onImgError" />
      <div class="absolute top-2.5 left-2.5">
        <UCheckbox :model-value="isSelected" @update:model-value="emit('toggle-select')" />
      </div>
      <div class="absolute top-2.5 right-2.5 flex flex-col items-end gap-1">
        <UBadge
          v-if="filterStatus !== 'trash' && !isScheduled"
          :label="getStatusLabel(post.status)"
          :color="getStatusColor(post.status)"
          variant="solid"
          size="sm"
          class="backdrop-blur-sm" />
        <UBadge
          v-if="isScheduled"
          :label="$t('admin.posts.editor.schedule_btn')"
          color="primary"
          variant="solid"
          size="sm"
          class="backdrop-blur-sm">
          <template #leading><UIcon name="i-tabler-clock" class="size-3" /></template>
        </UBadge>
        <UBadge v-if="post.metas?.is_featured === '1'" :label="$t('admin.posts.featured')" color="warning" variant="solid" size="sm" class="backdrop-blur-sm" />
        <UBadge v-if="post.metas?.is_banner === '1'" :label="$t('admin.posts.banner')" color="primary" variant="solid" size="sm" class="backdrop-blur-sm" />
      </div>
    </div>

    <!-- 内容 -->
    <div class="flex flex-col flex-1 p-3 gap-2">
      <h3 class="text-sm font-medium text-highlighted group-hover:text-primary transition-colors line-clamp-2 leading-snug">
        {{ post.title }}
      </h3>
      <p v-if="post.excerpt" class="text-xs text-muted line-clamp-2 flex-1">{{ post.excerpt }}</p>
      <div class="flex items-center justify-between pt-2 border-t border-default mt-auto">
        <span class="text-xs text-muted truncate max-w-24">{{ post.author?.nickname || post.author?.username || $t('common.unknown') }}</span>
        <div class="flex items-center gap-2 text-xs text-muted shrink-0">
          <div class="flex items-center gap-0.5">
            <UIcon name="i-tabler-eye" class="size-3" />{{ post.view_count }}
          </div>
          <div class="flex items-center gap-0.5">
            <UIcon name="i-tabler-message" class="size-3" />{{ post.comment_count }}
          </div>
          <UDropdownMenu :items="actions" :popper="{ placement: 'bottom-end' }">
            <UButton color="neutral" variant="ghost" icon="i-tabler-dots-vertical" square size="xs" />
          </UDropdownMenu>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface PostListItem {
  id: number
  title: string
  slug: string
  excerpt: string
  featured_img?: { id: number; url: string }
  status: string
  published_at?: string
  comment_count: number
  view_count: number
  author?: { id: number; username: string; nickname: string; avatar?: string }
  metas?: Record<string, string>
  created_at: string
  updated_at: string
}

const props = defineProps<{
  post: PostListItem
  isSelected: boolean
  actions: any[]
  filterStatus: string
}>()

const emit = defineEmits<{ 'toggle-select': [] }>()

const { defaultCover, onImgError } = usePostCover()
const { t } = useI18n()

const isScheduled = computed(() =>
  props.post.status === 'draft' &&
  !!props.post.published_at &&
  new Date(props.post.published_at).getTime() > Date.now()
)

const getStatusLabel = (status: string) => ({
  published: t('admin.posts.status_published'),
  draft: t('admin.posts.status_draft'),
  pending: t('admin.posts.status_pending'),
  trash: t('admin.posts.status_trash'),
  private: t('admin.posts.status_private'),
  archived: t('admin.posts.status_archived'),
}[status] ?? status)

type BadgeColor = 'success' | 'neutral' | 'warning' | 'error'
const getStatusColor = (status: string): BadgeColor => ({
  published: 'success', draft: 'neutral', pending: 'warning',
  trash: 'error', private: 'warning', archived: 'neutral',
}[status] as BadgeColor ?? 'neutral')
</script>
