<template>
  <div class="min-h-screen pb-16">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-8']">

      <!-- Title -->
      <div class="flex items-center gap-2 mb-6">
        <UIcon name="i-tabler-archive" class="size-6 text-primary shrink-0" />
        <h1 class="text-xl font-bold text-highlighted">{{ $t('site.archive.title') }}</h1>
        <span v-if="total" class="text-sm text-muted ml-1">{{ $t('site.archive.subtitle', { n: total }) }}</span>
      </div>

      <!-- Loading State -->
      <div v-if="displayLoading" class="mt-8 space-y-8">
        <div v-for="i in 3" :key="i">
          <USkeleton class="h-6 w-24 mb-4 rounded" />
          <div class="space-y-3 pl-4 border-l-2 border-default">
            <USkeleton v-for="j in 4" :key="j" class="h-5 w-full rounded" />
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="!Object.keys(grouped).length" class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden mt-2">
        <div class="py-20 flex flex-col items-center gap-3 text-center px-4">
          <div class="size-16 rounded-full bg-muted flex items-center justify-center">
            <UIcon name="i-tabler-inbox" class="size-8 text-muted" />
          </div>
          <p class="font-semibold text-highlighted">{{ $t('site.archive.no_posts') }}</p>
        </div>
      </div>

      <!-- Archive List -->
      <div v-else class="mt-8 space-y-10">
        <div v-for="(months, year) in grouped" :key="year">
          <h2 class="text-xl font-bold text-highlighted mb-4 flex items-center gap-2">
            <UIcon name="i-tabler-calendar" class="size-5 text-primary" />
            {{ year }}
            <UBadge color="neutral" variant="soft" size="sm">{{ yearCount(year) }}</UBadge>
          </h2>

          <div v-for="(posts, month) in months" :key="month" class="mb-6 pl-4 border-l-2 border-default">
            <h3 class="text-sm font-semibold text-muted mb-3 uppercase tracking-wide">{{ month }}</h3>
            <ul class="space-y-2">
              <li v-for="post in posts" :key="post.id" class="flex items-start gap-3 group">
                <span class="text-xs text-muted mt-1 shrink-0 w-8">
                  {{ dayOfMonth(post.published_at || post.created_at) }}
                </span>
                <NuxtLink
                  :to="`/posts/${post.slug}`"
                  class="text-default group-hover:text-primary transition-colors line-clamp-1 flex-1"
                >
                  {{ post.title }}
                </NuxtLink>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { containerClass } = useContainerWidth()
const { t } = useI18n()
import type { PostListItemResponse } from '~/types/api/post'

useHead({ title: t('site.archive.title') })

const postApi = usePostApi()
const rawLoading = ref(true)
const displayLoading = useMinLoading(rawLoading)

const allPosts = ref<PostListItemResponse[]>([])
const total    = ref(0)

onMounted(async () => {
  rawLoading.value = true
  try {
    const PAGE = 200
    let page = 1
    let fetched: PostListItemResponse[] = []
    let tot = 0

    do {
      const res = await postApi.getPosts({
        page,
        page_size: PAGE,
        status: 'published',
        post_type: '1',
        order_by: 'published_at',
        order: 'desc',
      })
      fetched = fetched.concat(res.data)
      tot = res.total
      page++
    } while (fetched.length < tot && fetched.length < 1000)

    allPosts.value = fetched
    total.value = tot
  } finally {
    rawLoading.value = false
  }
})

const { formatMonth, formatYear, formatDay } = useFormatDate()

/** Group posts into { year: { 'Month': posts[] } } */
const grouped = computed(() => {
  const result: Record<string, Record<string, PostListItemResponse[]>> = {}
  for (const post of allPosts.value) {
    const dateStr = post.published_at || post.created_at
    const year  = formatYear(dateStr)
    const month = formatMonth(dateStr)
    if (!result[year]) result[year] = {}
    if (!result[year][month]) result[year][month] = []
    result[year][month].push(post)
  }
  return result
})

const yearCount = (year: string) =>
  Object.values(grouped.value[year] ?? {}).reduce((s, p) => s + p.length, 0)

const dayOfMonth = (dateStr: string) => formatDay(dateStr)
</script>
