<script setup lang="ts">
import type { PostListItemResponse } from '~/types/api/post'

defineProps<{ posts: PostListItemResponse[]; loading?: boolean }>()
const emit = defineEmits<{ viewAll: [] }>()

const { defaultCover, onImgError } = usePostCover()

const formatDate = (iso: string) =>
  iso ? new Date(iso).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' }) : ''
const formatNum = (n: number) => n >= 1000 ? (n / 1000).toFixed(1) + 'k' : String(n)
</script>

<template>
  <div class="space-y-4">
    <!-- Skeleton -->
    <UCard v-if="loading">
      <template #header>
        <div class="flex items-center justify-between">
          <USkeleton class="h-5 w-32" />
          <USkeleton class="h-6 w-16 rounded" />
        </div>
      </template>
      <div class="divide-y divide-default -my-1">
        <div v-for="i in 5" :key="i" class="flex items-center gap-3 py-3">
          <USkeleton class="size-12 rounded-md shrink-0" />
          <div class="flex-1 space-y-2">
            <USkeleton class="h-4 w-3/4" />
            <USkeleton class="h-3 w-1/2" />
          </div>
        </div>
      </div>
    </UCard>

    <UCard v-else-if="posts.length > 0">
      <template #header>
        <div class="flex items-center justify-between">
          <h3 class="font-semibold text-highlighted flex items-center gap-2">
            <UIcon name="i-tabler-file-text" class="size-5 text-primary" />
            {{ $t('site.user.latest_posts') }}
          </h3>
          <UButton color="neutral" variant="ghost" size="xs" @click="emit('viewAll')">
            {{ $t('site.user.view_all') }}
          </UButton>
        </div>
      </template>
      <div class="divide-y divide-default -my-1">
        <NuxtLink
          v-for="post in posts"
          :key="post.id"
          :to="`/posts/${post.slug}`"
          class="group flex items-center gap-3 py-3">
          <BaseImg
            :src="post.featured_img?.url || defaultCover"
            :alt="post.title"
            class="size-12 rounded-md object-cover shrink-0"
            @error="onImgError" />
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-default group-hover:text-primary transition-colors line-clamp-1">
              {{ post.title }}
            </p>
            <div class="flex items-center gap-2 text-xs text-muted mt-0.5">
              <span>{{ formatDate(post.published_at || post.created_at) }}</span>
              <span>·</span>
              <span class="flex items-center gap-0.5">
                <UIcon name="i-tabler-eye" class="size-3" />{{ formatNum(post.view_count) }}
              </span>
              <span>·</span>
              <span class="flex items-center gap-0.5">
                <UIcon name="i-tabler-message-circle" class="size-3" />{{ post.comment_count }}
              </span>
            </div>
          </div>
        </NuxtLink>
      </div>
    </UCard>
    <UCard v-else-if="!loading">
      <div class="text-center py-8 text-muted">
        <UIcon name="i-tabler-file-off" class="size-10 mx-auto mb-2" />
        <p class="text-sm">{{ $t('site.user.no_posts_empty') }}</p>
      </div>
    </UCard>
  </div>
</template>
