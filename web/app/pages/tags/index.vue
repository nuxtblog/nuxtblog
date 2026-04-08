<template>
  <div class="min-h-screen pb-16">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-8']">

      <!-- Title -->
      <div class="flex items-center gap-2 mb-6">
        <UIcon name="i-tabler-tags" class="size-6 text-primary shrink-0" />
        <h1 class="text-xl font-bold text-highlighted">{{ $t('site.tags.title') }}</h1>
        <span class="text-sm text-muted ml-1">{{ $t('site.tags.subtitle', { n: tags.length }) }}</span>
      </div>

      <!-- Stats row -->
      <div class="grid grid-cols-4 rounded-md overflow-hidden ring-1 ring-default bg-default shadow-sm mb-6">
        <div class="flex flex-col items-center py-4">
          <span class="text-2xl font-bold text-primary">{{ tags.length }}</span>
          <span class="text-xs text-muted mt-0.5">{{ $t('site.tags.total_tags') }}</span>
        </div>
        <div class="flex flex-col items-center py-4 border-l border-default">
          <span class="text-2xl font-bold text-primary">{{ totalArticles }}</span>
          <span class="text-xs text-muted mt-0.5">{{ $t('site.tags.total_articles') }}</span>
        </div>
        <div class="flex flex-col items-center py-4 border-l border-default">
          <span class="text-2xl font-bold text-highlighted">{{ Math.round(averageArticles * 10) / 10 }}</span>
          <span class="text-xs text-muted mt-0.5">{{ $t('site.tags.avg_articles') }}</span>
        </div>
        <div class="flex flex-col items-center py-4 border-l border-default">
          <span class="text-2xl font-bold text-highlighted">{{ popularTagCount }}</span>
          <span class="text-xs text-muted mt-0.5">{{ $t('site.tags.popular_tags') }}</span>
        </div>
      </div>

      <!-- Search and Filter -->
      <div class="flex flex-col sm:flex-row gap-3 mb-6">
        <UInput
          v-model="searchQuery"
          :placeholder="$t('site.tags.search_placeholder')"
          leading-icon="i-tabler-search"
          class="flex-1" />
        <USelect v-model="sortBy" :items="sortOptions" value-key="value" label-key="label" class="w-auto" />
        <USelect v-model="viewMode" :items="viewModeOptions" value-key="value" label-key="label" class="w-auto" />
      </div>

      <!-- Loading skeleton -->
      <div v-if="isLoading" class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">
        <div class="divide-y divide-default">
          <div v-for="i in 6" :key="i" class="flex items-center gap-3 px-4 py-3.5">
            <USkeleton class="size-10 rounded-md shrink-0" />
            <div class="flex-1 space-y-1.5">
              <USkeleton class="h-4 w-24" />
              <USkeleton class="h-3 w-16" />
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredTags.length === 0" class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">
        <div class="py-20 flex flex-col items-center gap-3 text-center px-4">
          <div class="size-16 rounded-full bg-muted flex items-center justify-center">
            <UIcon name="i-tabler-tags" class="size-8 text-muted" />
          </div>
          <p class="font-semibold text-highlighted">{{ $t('site.tags.no_tags') }}</p>
          <p class="text-sm text-muted">
            {{ searchQuery ? $t('site.tags.no_results') : $t('site.tags.no_any') }}
          </p>
        </div>
      </div>

      <template v-else>
        <!-- Tag Cloud View -->
        <div v-if="viewMode === 'cloud'" class="rounded-md bg-default ring-1 ring-default shadow-sm p-8">
          <div class="flex flex-wrap gap-3 justify-center items-center">
            <NuxtLink
              v-for="tag in pagedTags"
              :key="tag.term_taxonomy_id"
              :to="`/posts?tag=${tag.term_taxonomy_id}`"
              class="inline-flex items-center gap-2 px-4 py-2 rounded-full border border-default hover:border-primary hover:bg-primary/10 transition-all"
              :style="{ fontSize: getTagSize(tag.count) + 'rem' }">
              <UIcon name="i-tabler-tag" class="size-5" />
              <span class="font-medium text-default">{{ tag.name }}</span>
              <UBadge color="primary" variant="soft" size="sm">{{ tag.count }}</UBadge>
            </NuxtLink>
          </div>
        </div>

        <!-- Card View -->
        <div
          v-else-if="viewMode === 'card'"
          class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3 sm:gap-4">
          <NuxtLink
            v-for="(tag, index) in pagedTags"
            :key="tag.term_taxonomy_id"
            :to="`/posts?tag=${tag.term_taxonomy_id}`"
            class="rounded-md ring-1 ring-default bg-default shadow-sm hover:shadow-md transition-shadow group">
            <div class="p-5">
              <div class="flex items-start gap-3 mb-2">
                <div
                  class="w-12 h-12 rounded-md flex items-center justify-center shrink-0"
                  :style="{ backgroundColor: getColor(index) + '20', color: getColor(index) }">
                  <UIcon name="i-tabler-tag" class="size-5" />
                </div>
                <div class="flex-1 min-w-0">
                  <h3 class="font-semibold text-highlighted text-base truncate group-hover:text-primary transition-colors">
                    {{ tag.name }}
                  </h3>
                  <p class="text-sm text-muted">{{ $t('site.tags.post_count', { n: tag.count }) }}</p>
                </div>
              </div>
              <div class="flex items-center justify-between pt-3 border-t border-default">
                <span class="text-xs text-muted">{{ tag.slug }}</span>
                <UBadge v-if="tag.count > 3" color="primary" variant="soft" size="sm">{{ $t('site.tags.hot') }}</UBadge>
              </div>
            </div>
          </NuxtLink>
        </div>

        <!-- List View -->
        <div v-else class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">
          <NuxtLink
            v-for="(tag, index) in pagedTags"
            :key="tag.term_taxonomy_id"
            :to="`/posts?tag=${tag.term_taxonomy_id}`"
            class="flex items-center justify-between px-4 py-3.5 hover:bg-elevated/50 transition-colors group"
            :class="index > 0 ? 'border-t border-default' : ''">
            <div class="flex items-center gap-3 flex-1 min-w-0">
              <div class="size-10 rounded-md flex items-center justify-center bg-primary/10 group-hover:bg-primary/20 transition-colors shrink-0">
                <UIcon name="i-tabler-tag" class="size-5 text-primary" />
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="text-sm font-semibold text-highlighted truncate mb-0.5">{{ tag.name }}</h3>
                <p class="text-xs text-muted">{{ $t('site.tags.post_count', { n: tag.count }) }}</p>
              </div>
            </div>
            <div class="flex items-center gap-3">
              <UBadge v-if="tag.count > 3" color="primary" variant="soft" class="hidden sm:inline-flex">{{ $t('site.tags.hot_tag') }}</UBadge>
              <UIcon name="i-tabler-chevron-right" class="size-5 text-muted group-hover:text-primary group-hover:translate-x-1 transition-all" />
            </div>
          </NuxtLink>
        </div>

        <!-- Pagination -->
        <div v-if="totalPages > 1" class="flex justify-center mt-6">
          <UPagination
            v-model:page="currentPage"
            :total="filteredTags.length"
            :items-per-page="pageSize" />
        </div>
      </template>

    </div>
  </div>
</template>

<script setup lang="ts">
const { containerClass } = useContainerWidth()
const { t } = useI18n()
import type { TermDetailResponse } from '~/types/api/term'

const COLORS = [
  '#3B82F6', '#10B981', '#F59E0B', '#8B5CF6', '#EF4444',
  '#EC4899', '#14B8A6', '#F97316', '#6366F1', '#84CC16',
]
const getColor = (index: number) => COLORS[index % COLORS.length]

const sortOptions = computed(() => [
  { label: t('site.tags.sort_name'), value: 'name' },
  { label: t('site.tags.sort_count'), value: 'count' },
])
const viewModeOptions = computed(() => [
  { label: t('site.tags.view_cloud'), value: 'cloud' },
  { label: t('site.tags.view_card'), value: 'card' },
  { label: t('site.tags.view_list'), value: 'list' },
])

const termApi = useTermApi()

const isLoading = ref(true)
const searchQuery = ref('')
const sortBy = ref<'name' | 'count'>('count')
const viewMode = ref<'cloud' | 'card' | 'list'>('cloud')
const currentPage = ref(1)
const pageSize = 12

const tags = ref<TermDetailResponse[]>([])

const filteredTags = computed(() => {
  let result = tags.value
  if (searchQuery.value) {
    const q = searchQuery.value.trim().toLowerCase()
    result = result.filter(t => t.name.toLowerCase().includes(q))
  }
  return [...result].sort((a, b) =>
    sortBy.value === 'name' ? a.name.localeCompare(b.name, 'zh-CN') : b.count - a.count
  )
})

const pagedTags = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredTags.value.slice(start, start + pageSize)
})

const totalPages = computed(() => Math.ceil(filteredTags.value.length / pageSize))
const totalArticles = computed(() => tags.value.reduce((s, t) => s + t.count, 0))
const averageArticles = computed(() => tags.value.length ? totalArticles.value / tags.value.length : 0)
const popularTagCount = computed(() => tags.value.filter(t => t.count > 3).length)

const getTagSize = (count: number) => {
  const minSize = 0.875
  const maxSize = 2
  const counts = tags.value.map(t => t.count)
  const maxCount = Math.max(...counts, 1)
  const minCount = Math.min(...counts, 0)
  if (maxCount === minCount) return (minSize + maxSize) / 2
  return minSize + ((count - minCount) / (maxCount - minCount)) * (maxSize - minSize)
}

onMounted(async () => {
  try {
    tags.value = await termApi.getTerms({ taxonomy: 'tag' })
  } catch (e) {
    console.error('Failed to load tags:', e)
  } finally {
    isLoading.value = false
  }
})
</script>
