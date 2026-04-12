<template>
  <div class="min-h-screen pb-16">
    <PageHeader
      icon="i-tabler-file-text"
      :title="$t('site.articles.title')"
      :subtitle="total ? $t('site.articles.subtitle', { n: total }) : undefined">
      <template #toolbar>
        <div class="flex flex-col sm:flex-row gap-3">
          <UInput
            v-model="searchQuery"
            :placeholder="$t('site.articles.search_placeholder')"
            leading-icon="i-tabler-search"
            class="flex-1"
            size="md"
            @keyup.enter="onSearch">
            <template v-if="searchQuery" #trailing>
              <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="searchQuery = ''; onSearch()" />
            </template>
          </UInput>
          <UButton @click="onSearch" color="primary">{{ $t('site.articles.search_btn') }}</UButton>
        </div>
      </template>
    </PageHeader>

    <!-- Plugin slot: posts list top -->
    <ClientOnly><ContributionSlot name="public:posts-top" /></ClientOnly>

    <PageContent>
      <!-- Active tag/category filter hint -->
      <div v-if="activeTagId || activeCategoryId" class="flex items-center gap-2 mb-4">
        <UBadge color="primary" variant="soft" class="flex items-center gap-1">
          <UIcon name="i-tabler-filter" class="size-3.5" />
          {{ activeTagId ? $t('site.articles.filter_by_tag') : $t('site.articles.filter_by_category') }}
        </UBadge>
        <button class="text-xs text-muted hover:text-primary" @click="clearFilter">
          {{ $t('site.articles.clear_filter') }}
        </button>
      </div>

      <!-- Main card -->
      <div class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">
        <!-- Loading skeleton -->
        <div v-if="loading" class="divide-y divide-default">
          <div v-for="i in 5" :key="i" class="flex gap-4 px-4 py-4">
            <USkeleton class="h-20 w-32 rounded-md shrink-0" />
            <div class="flex-1 space-y-2 pt-1">
              <USkeleton class="h-4 w-3/4" />
              <USkeleton class="h-3 w-full" />
              <USkeleton class="h-3 w-1/2" />
            </div>
          </div>
        </div>

        <!-- Empty state -->
        <div v-else-if="posts.length === 0" class="py-20 flex flex-col items-center gap-3 text-center px-4">
          <div class="size-16 rounded-full bg-muted flex items-center justify-center">
            <UIcon name="i-tabler-file-off" class="size-8 text-muted" />
          </div>
          <p class="font-semibold text-highlighted">{{ $t('site.articles.no_posts') }}</p>
        </div>

        <!-- Post list -->
        <div v-else class="divide-y divide-default">
          <NuxtLink
            v-for="post in posts"
            :key="post.id"
            :to="`/posts/${post.slug}`"
            class="flex gap-4 px-4 py-4 hover:bg-elevated/50 transition-colors group">
            <div class="w-28 h-20 rounded-md overflow-hidden shrink-0 bg-muted">
              <BaseImg
                :src="post.featured_img?.url || defaultCover"
                :alt="post.title"
                class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
            </div>
            <div class="flex-1 min-w-0 py-0.5">
              <h3 class="font-semibold text-highlighted group-hover:text-primary transition-colors line-clamp-2 mb-1">
                {{ post.title }}
              </h3>
              <p v-if="post.excerpt" class="text-sm text-muted line-clamp-2 mb-2">
                {{ post.excerpt }}
              </p>
              <div class="flex items-center gap-4 text-xs text-muted">
                <span v-if="post.author">{{ post.author.nickname }}</span>
                <span v-if="post.published_at">{{ new Date(post.published_at).toLocaleDateString('zh-CN') }}</span>
                <span class="flex items-center gap-1">
                  <UIcon name="i-tabler-eye" class="size-3.5" />{{ post.view_count }}
                </span>
                <span class="flex items-center gap-1">
                  <UIcon name="i-tabler-message" class="size-3.5" />{{ post.comment_count }}
                </span>
              </div>
            </div>
          </NuxtLink>
        </div>
      </div>
    </PageContent>

    <!-- Pagination -->
    <PageFooter v-if="totalPages > 1">
      <div class="flex justify-center">
        <UPagination
          v-model:page="currentPage"
          :total="total"
          :items-per-page="pageSize" />
      </div>
    </PageFooter>

    <!-- Plugin slot: posts list bottom -->
    <ClientOnly><ContributionSlot name="public:posts-bottom" /></ClientOnly>
  </div>
</template>

<script setup lang="ts">
const { t } = useI18n()
import type { PostListItemResponse } from '~/types/api/post'

const route = useRoute()
const postApi = usePostApi()
const { defaultCover } = usePostCover()

const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)
const posts = ref<PostListItemResponse[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 12
const searchQuery = ref('')

const activeTagId = computed(() => {
  const v = route.query.tag
  return v ? Number(v) : null
})
const activeCategoryId = computed(() => {
  const v = route.query.category
  return v ? Number(v) : null
})

const totalPages = computed(() => Math.ceil(total.value / pageSize))

const fetchPosts = async () => {
  rawLoading.value = true
  try {
    const params: Record<string, any> = {
      status: 'published',
      page: currentPage.value,
      page_size: pageSize,
    }
    if (searchQuery.value) params.keyword = searchQuery.value
    if (activeTagId.value) params.term_taxonomy_id = activeTagId.value
    else if (activeCategoryId.value) params.term_taxonomy_id = activeCategoryId.value

    const res = await postApi.getPosts(params)
    posts.value = res.data || []
    total.value = res.total || 0
  } catch (e) {
    console.error('Failed to load posts:', e)
  } finally {
    rawLoading.value = false
  }
}

const onSearch = () => {
  currentPage.value = 1
  fetchPosts()
}

const clearFilter = () => {
  navigateTo('/posts')
}

watch([currentPage, activeTagId, activeCategoryId], () => {
  fetchPosts()
})

onMounted(() => {
  fetchPosts()
})
</script>
