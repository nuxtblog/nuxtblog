<script setup lang="ts">
import type { MomentItem } from '~/types/api/moment'

const momentApi = useMomentApi()
const { containerClass } = useContainerWidth()

const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)
const moments = ref<MomentItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

useHead({
  title: '动态',
  meta: [{ name: 'description', content: '查看最新动态' }],
})

async function fetchMoments() {
  rawLoading.value = true
  try {
    const res = await momentApi.getMoments({ page: currentPage.value, page_size: pageSize.value, visibility: 1 })
    moments.value = res.data ?? []
    total.value = res.total ?? 0
  } catch (e: any) {
    // silently fail on public page
  } finally {
    rawLoading.value = false
  }
}

function getAuthorInitials(item: MomentItem) {
  const name = item.author?.nickname || item.author?.username || '?'
  return name.charAt(0).toUpperCase()
}

function formatRelativeTime(d: string) {
  if (!d) return ''
  const diff = Date.now() - new Date(d).getTime()
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 7) return `${days}天前`
  return new Date(d).toLocaleDateString('zh-CN')
}

function isImageMedia(mime: string) {
  return mime.startsWith('image/')
}

watch(currentPage, fetchMoments)
onMounted(fetchMoments)
</script>

<template>
  <div class="min-h-screen pb-16">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-8']">
      <!-- Page header -->
      <div class="flex items-center gap-3 mb-8">
        <div class="size-10 rounded-lg bg-primary/10 flex items-center justify-center">
          <UIcon name="i-tabler-camera" class="size-6 text-primary" />
        </div>
        <div>
          <h1 class="text-2xl font-bold text-highlighted">动态</h1>
          <p class="text-sm text-muted mt-0.5">查看最新动态</p>
        </div>
      </div>

      <!-- Skeleton -->
      <div v-if="loading" class="space-y-4 max-w-2xl mx-auto">
        <div v-for="i in 5" :key="i" class="rounded-lg bg-default ring-1 ring-default shadow-sm p-4">
          <div class="flex items-center gap-3 mb-3">
            <USkeleton class="size-10 rounded-full shrink-0" />
            <div class="space-y-1.5 flex-1">
              <USkeleton class="h-4 w-1/3" />
              <USkeleton class="h-3 w-1/4" />
            </div>
          </div>
          <USkeleton class="h-4 w-full mb-1.5" />
          <USkeleton class="h-4 w-3/4" />
          <div class="flex gap-4 mt-3">
            <USkeleton class="h-3 w-12" />
            <USkeleton class="h-3 w-12" />
            <USkeleton class="h-3 w-12" />
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="moments.length === 0" class="flex flex-col items-center justify-center py-20 max-w-2xl mx-auto">
        <div class="size-16 rounded-full bg-muted flex items-center justify-center mb-4">
          <UIcon name="i-tabler-camera-off" class="size-8 text-muted" />
        </div>
        <h3 class="text-lg font-medium text-highlighted mb-1">暂无动态</h3>
        <p class="text-sm text-muted">还没有发布任何动态</p>
      </div>

      <!-- Moment feed -->
      <div v-else class="max-w-2xl mx-auto space-y-4">
        <div
          v-for="moment in moments" :key="moment.id"
          class="rounded-lg bg-default ring-1 ring-default shadow-sm p-4">
          <!-- Author row -->
          <div class="flex items-center gap-3 mb-3">
            <div
              v-if="moment.author?.avatar"
              class="size-10 rounded-full overflow-hidden shrink-0">
              <img :src="moment.author.avatar" :alt="moment.author.nickname || moment.author.username" class="size-full object-cover" />
            </div>
            <div v-else class="size-10 rounded-full bg-primary/10 flex items-center justify-center shrink-0 text-sm font-semibold text-primary">
              {{ getAuthorInitials(moment) }}
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-highlighted truncate">
                {{ moment.author?.nickname || moment.author?.username || '匿名用户' }}
              </p>
              <p class="text-xs text-muted">{{ formatRelativeTime(moment.created_at) }}</p>
            </div>
          </div>

          <!-- Content -->
          <p class="text-sm text-default leading-relaxed whitespace-pre-wrap">{{ moment.content }}</p>

          <!-- Media grid -->
          <div
            v-if="moment.media && moment.media.length > 0"
            class="grid gap-1 mt-3 rounded-md overflow-hidden"
            :class="moment.media.length === 1 ? 'grid-cols-1' : moment.media.length === 2 ? 'grid-cols-2' : 'grid-cols-3'">
            <div v-for="media in moment.media" :key="media.id" class="aspect-square overflow-hidden bg-muted">
              <img
                v-if="isImageMedia(media.mime_type)"
                :src="media.url"
                :alt="'图片'"
                class="size-full object-cover" />
              <div v-else class="size-full flex items-center justify-center">
                <UIcon name="i-tabler-file" class="size-8 text-muted" />
              </div>
            </div>
          </div>

          <!-- Stats row -->
          <div v-if="moment.stats" class="flex items-center gap-4 mt-3 pt-3 border-t border-default">
            <span class="flex items-center gap-1 text-xs text-muted">
              <UIcon name="i-tabler-eye" class="size-3.5" />
              {{ moment.stats.view_count }}
            </span>
            <span class="flex items-center gap-1 text-xs text-muted">
              <UIcon name="i-tabler-heart" class="size-3.5" />
              {{ moment.stats.like_count }}
            </span>
            <span class="flex items-center gap-1 text-xs text-muted">
              <UIcon name="i-tabler-message" class="size-3.5" />
              {{ moment.stats.comment_count }}
            </span>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="total > pageSize" class="flex justify-center pt-4">
          <UPagination v-model:page="currentPage" :total="total" :items-per-page="pageSize" />
        </div>
      </div>
    </div>
  </div>
</template>
