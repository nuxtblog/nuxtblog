<template>
  <div class="flex items-center group gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all">
    <UCheckbox :model-value="isSelected" @update:model-value="emit('toggle-select')" />

    <!-- 封面 -->
    <div class="shrink-0 h-16 w-24 rounded overflow-hidden bg-muted">
      <img
        :src="post.featured_img?.url || defaultCover"
        :alt="post.title"
        class="h-full w-full object-cover group-hover:scale-105 transition-transform duration-300"
        @error="onImgError" />
    </div>

    <!-- 信息 -->
    <div class="flex-1 min-w-0">
      <div class="flex items-start justify-between gap-4">
        <div class="flex-1 min-w-0">
          <div class="flex items-center gap-2 flex-wrap">
            <h3 class="text-base font-medium text-highlighted group-hover:text-primary transition-colors truncate">
              {{ post.title }}
            </h3>
            <UBadge v-if="isScheduled" :label="$t('admin.posts.editor.schedule_btn')" color="primary" variant="soft" size="sm">
              <template #leading><UIcon name="i-tabler-clock" class="size-3" /></template>
            </UBadge>
            <UBadge v-if="post.metas?.is_featured === '1'" :label="$t('admin.posts.featured')" color="warning" variant="soft" size="sm" />
            <UBadge v-if="post.metas?.is_banner === '1'" :label="$t('admin.posts.banner')" color="primary" variant="soft" size="sm" />
          </div>
          <p v-if="post.excerpt" class="text-sm text-muted mt-1 line-clamp-1">{{ post.excerpt }}</p>
          <div class="flex items-center gap-4 mt-2 text-xs text-muted">
            <span>{{ post.author?.nickname || post.author?.username || $t('admin.posts.unknown_author') }}</span>
            <span>{{ formatDate(post.created_at) }}</span>
            <div class="flex items-center gap-1">
              <UIcon name="i-tabler-eye" class="size-3.5" />{{ post.view_count }}
            </div>
            <div class="flex items-center gap-1">
              <UIcon name="i-tabler-message" class="size-3.5" />{{ post.comment_count }}
            </div>
          </div>
        </div>

        <!-- 状态 + 操作 -->
        <div class="flex items-center gap-3 shrink-0">
          <UBadge
            v-if="filterStatus !== 'trash'"
            :label="getStatusLabel(post.status)"
            :color="getStatusColor(post.status)"
            variant="soft"
            size="sm" />
          <UDropdownMenu :items="actions" :popper="{ placement: 'bottom-end' }">
            <UButton
              color="neutral"
              variant="ghost"
              icon="i-tabler-dots-vertical"
              square
              size="xs"
              class="opacity-0 group-hover:opacity-100 transition-opacity" />
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
const { formatRelative } = useFormatDate()

const isScheduled = computed(() =>
  props.post.status === 'draft' &&
  !!props.post.published_at &&
  new Date(props.post.published_at).getTime() > Date.now()
)

const formatDate = (s: string) => formatRelative(s)

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
