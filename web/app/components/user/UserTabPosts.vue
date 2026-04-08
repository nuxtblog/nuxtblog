<script setup lang="ts">
import type { PostListItemResponse } from '~/types/api/post'

const props = defineProps<{ userId: number }>()

const postApi = usePostApi()
const { defaultCover, onImgError } = usePostCover()

const page = ref(1)
const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)
const posts = ref<PostListItemResponse[]>([])
const total = ref(0)

const formatDate = (iso: string) =>
  iso ? new Date(iso).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' }) : ''
const formatNum = (n: number) => n >= 1000 ? (n / 1000).toFixed(1) + 'k' : String(n)

async function load() {
  rawLoading.value = true
  try {
    const res = await postApi.getPosts({ author_id: props.userId, status: '2', page: page.value, page_size: 10 })
    posts.value = res.data
    total.value = res.total
  } finally {
    rawLoading.value = false
  }
}

watch(page, load)
onMounted(load)
</script>

<template>
  <BaseCardList :loading="loading" :empty="posts.length === 0" :skeleton-count="3">

    <template #skeleton>
      <div class="flex gap-4">
        <USkeleton class="w-20 h-20 md:w-28 md:h-28 rounded-md shrink-0" />
        <div class="flex-1 space-y-2">
          <USkeleton class="h-5 w-3/4" />
          <USkeleton class="h-4 w-full" />
          <USkeleton class="h-3 w-1/2" />
        </div>
      </div>
    </template>

    <template #empty>
      <div class="text-center py-12 text-muted">
        <UIcon name="i-tabler-file-off" class="size-12 mx-auto mb-3" />
        <p>{{ $t('site.user.no_posts') }}</p>
      </div>
    </template>

    <NuxtLink
      v-for="post in posts"
      :key="post.id"
      :to="`/posts/${post.slug}`"
      class="group block">
      <UCard class="hover:shadow-md transition-shadow">
        <div class="flex gap-4">
          <BaseImg
            :src="post.featured_img?.url || defaultCover"
            :alt="post.title"
            class="w-20 h-20 md:w-28 md:h-28 rounded-md object-cover shrink-0"
            @error="onImgError" />
          <div class="flex-1 min-w-0">
            <h3 class="font-semibold text-highlighted group-hover:text-primary transition-colors line-clamp-2 mb-1">
              {{ post.title }}
            </h3>
            <p class="text-muted text-sm line-clamp-2 mb-3">{{ post.excerpt }}</p>
            <div class="flex flex-wrap items-center gap-3 text-xs text-muted">
              <span>{{ formatDate(post.published_at || post.created_at) }}</span>
              <span class="flex items-center gap-0.5">
                <UIcon name="i-tabler-eye" class="size-3" />{{ formatNum(post.view_count) }}
              </span>
              <span class="flex items-center gap-0.5">
                <UIcon name="i-tabler-message-circle" class="size-3" />{{ post.comment_count }}
              </span>
            </div>
          </div>
        </div>
      </UCard>
    </NuxtLink>

    <div v-if="total > 10" class="flex justify-center pt-2">
      <UPagination v-model:page="page" :total="total" :items-per-page="10" />
    </div>

  </BaseCardList>
</template>
