<template>
  <div class="min-h-screen pb-16">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-8']">

      <!-- Loading skeleton -->
      <div v-if="isLoading" class="space-y-6">
        <div class="flex items-center gap-3 mb-6">
          <USkeleton class="size-6 rounded-full" />
          <div class="space-y-2">
            <USkeleton class="h-6 w-40" />
            <USkeleton class="h-4 w-24" />
          </div>
        </div>
        <div class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden divide-y divide-default">
          <div v-for="i in 5" :key="i" class="flex gap-4 px-4 py-4">
            <USkeleton class="w-28 h-20 rounded-md shrink-0" />
            <div class="flex-1 space-y-2 pt-1">
              <USkeleton class="h-4 w-3/4" />
              <USkeleton class="h-3 w-full" />
              <USkeleton class="h-3 w-1/2" />
            </div>
          </div>
        </div>
      </div>

      <!-- Not found -->
      <div v-else-if="!category" class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">
        <div class="py-20 flex flex-col items-center gap-3 text-center px-4">
          <div class="size-16 rounded-full bg-muted flex items-center justify-center">
            <UIcon name="i-tabler-folder-off" class="size-8 text-muted" />
          </div>
          <p class="font-semibold text-highlighted">{{ $t('site.categories.not_found') }}</p>
          <NuxtLink to="/categories">
            <UButton color="neutral" variant="outline" size="sm">{{ $t('site.categories.back_to_list') }}</UButton>
          </NuxtLink>
        </div>
      </div>

      <template v-else>
        <!-- Plugin slot: category top -->
        <ClientOnly><ContributionSlot name="public:category-top" /></ClientOnly>

        <!-- Title -->
        <div class="flex items-center gap-2 mb-2">
          <UIcon name="i-tabler-stack-2" class="size-6 text-primary shrink-0" />
          <h1 class="text-xl font-bold text-highlighted">{{ category.name }}</h1>
          <span class="text-sm text-muted ml-1">{{ $t('site.categories.post_count', { n: total || category.count }) }}</span>
        </div>
        <p v-if="category.description" class="text-sm text-muted mb-6 ml-8">{{ category.description }}</p>
        <div v-if="!category.description" class="mb-6" />

        <!-- Posts card -->
        <div class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">

          <!-- Pagination loading -->
          <div v-if="postsLoading" class="divide-y divide-default">
            <div v-for="i in 5" :key="i" class="flex gap-4 px-4 py-4">
              <USkeleton class="w-28 h-20 rounded-md shrink-0" />
              <div class="flex-1 space-y-2 pt-1">
                <USkeleton class="h-4 w-3/4" />
                <USkeleton class="h-3 w-full" />
                <USkeleton class="h-3 w-1/2" />
              </div>
            </div>
          </div>

          <!-- Empty -->
          <div v-else-if="posts.length === 0" class="py-20 flex flex-col items-center gap-3 text-center px-4">
            <div class="size-16 rounded-full bg-muted flex items-center justify-center">
              <UIcon name="i-tabler-file-off" class="size-8 text-muted" />
            </div>
            <p class="font-semibold text-highlighted">{{ $t('site.categories.no_posts_in') }}</p>
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
                <p v-if="post.excerpt" class="text-sm text-muted line-clamp-2 mb-2">{{ post.excerpt }}</p>
                <div class="flex items-center gap-4 text-xs text-muted">
                  <span v-if="post.author">{{ post.author.nickname }}</span>
                  <span v-if="post.published_at">{{ new Date(post.published_at).toLocaleDateString('zh-CN') }}</span>
                  <span class="flex items-center gap-1">
                    <UIcon name="i-tabler-eye" class="size-3.5" />{{ post.view_count }}
                  </span>
                </div>
              </div>
            </NuxtLink>
          </div>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex justify-center mt-6">
          <UPagination v-model:page="currentPage" :total="total" :items-per-page="pageSize" />
        </div>

        <!-- Plugin slot: category bottom -->
        <ClientOnly><ContributionSlot name="public:category-bottom" /></ClientOnly>
      </template>

    </div>
  </div>
</template>

<script setup lang="ts">
const { containerClass } = useContainerWidth()
import type { TermDetailResponse } from '~/types/api/term'
import type { PostListItemResponse } from '~/types/api/post'

const route = useRoute()
const slug = route.params.categoryId as string

const termApi = useTermApi()
const postApi = usePostApi()
const { defaultCover } = usePostCover()

const rawLoading = ref(true)
const isLoading = useMinLoading(rawLoading)
const postsLoading = ref(false)
const category = ref<TermDetailResponse | null>(null)
const allCategoryIds = ref<number[]>([])
const posts = ref<PostListItemResponse[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 12

const totalPages = computed(() => Math.ceil(total.value / pageSize))

const collectDescendantIds = (terms: TermDetailResponse[], rootTaxId: number): number[] => {
  const result: number[] = [rootTaxId]
  const queue = [rootTaxId]
  while (queue.length > 0) {
    const parentId = queue.shift()!
    for (const t of terms) {
      if (t.parent_id === parentId) {
        result.push(t.term_taxonomy_id)
        queue.push(t.term_taxonomy_id)
      }
    }
  }
  return result
}

const fetchPosts = async (ids: number[]) => {
  postsLoading.value = true
  try {
    const res = await postApi.getPosts({
      include_category_ids: ids.join(','),
      status: 'published',
      page: currentPage.value,
      page_size: pageSize,
    })
    posts.value = res.data || []
    total.value = res.total || 0
  } catch (e) {
    console.error('Failed to load posts:', e)
  } finally {
    postsLoading.value = false
  }
}

watch(currentPage, () => {
  if (allCategoryIds.value.length) fetchPosts(allCategoryIds.value)
})

onMounted(async () => {
  try {
    const terms = await termApi.getTerms({ taxonomy: 'category' })
    category.value = terms.find(t => t.slug === slug) || null
    if (category.value) {
      allCategoryIds.value = collectDescendantIds(terms, category.value.term_taxonomy_id)
      await fetchPosts(allCategoryIds.value)
    }
  } catch (e) {
    console.error('Failed to load category:', e)
  } finally {
    rawLoading.value = false
  }
})
</script>
