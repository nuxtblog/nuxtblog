<template>
  <div class="min-h-screen pb-16">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-8']">

      <!-- Plugin slot: categories top -->
      <ClientOnly><ContributionSlot name="public:categories-top" /></ClientOnly>

      <!-- Title -->
      <div class="flex items-center gap-2 mb-6">
        <UIcon name="i-tabler-stack-2" class="size-6 text-primary shrink-0" />
        <h1 class="text-xl font-bold text-highlighted">{{ $t('site.categories.title') }}</h1>
        <span class="text-sm text-muted ml-1">{{ $t('site.categories.subtitle', { n: filteredCategories.length }) }}</span>
      </div>

      <!-- Search and Filter -->
      <div class="flex flex-col sm:flex-row gap-3 mb-6">
        <UInput
          v-model="searchQuery"
          :placeholder="$t('site.categories.search_placeholder')"
          leading-icon="i-tabler-search"
          class="flex-1" />
        <USelect
          v-model="sortBy"
          :items="sortOptions"
          value-key="value"
          label-key="label" />
      </div>

      <!-- Loading skeleton -->
      <div v-if="isLoading" class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3 sm:gap-4">
        <div v-for="i in 8" :key="i" class="rounded-md bg-default ring-1 ring-default shadow-sm p-5">
          <div class="flex items-start gap-3 mb-2">
            <USkeleton class="size-12 rounded-md shrink-0" />
            <div class="flex-1 space-y-2">
              <USkeleton class="h-4 w-24" />
              <USkeleton class="h-3 w-16" />
            </div>
          </div>
          <USkeleton class="h-3 w-full mt-2" />
          <USkeleton class="h-3 w-3/4 mt-1" />
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredCategories.length === 0" class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">
        <div class="py-20 flex flex-col items-center gap-3 text-center px-4">
          <div class="size-16 rounded-full bg-muted flex items-center justify-center">
            <UIcon name="i-tabler-category" class="size-8 text-muted" />
          </div>
          <p class="font-semibold text-highlighted">{{ $t('site.categories.no_categories') }}</p>
          <p class="text-sm text-muted">
            {{ searchQuery ? $t('site.categories.no_results') : $t('site.categories.no_any') }}
          </p>
        </div>
      </div>

      <div v-else class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3 sm:gap-4">
        <TransitionGroup name="list" tag="div" class="contents">
          <NuxtLink
            v-for="(category, index) in paginatedCategories"
            :key="category.term_taxonomy_id"
            :to="`/categories/${category.slug}`"
            class="rounded-md ring-1 ring-default bg-default shadow-sm hover:shadow-md transition-shadow group">
            <div class="p-5">
              <div class="flex items-start gap-3 mb-2">
                <div
                  class="w-12 h-12 rounded-md flex items-center justify-center shrink-0"
                  :style="{ backgroundColor: getColor(index) + '20', color: getColor(index) }">
                  <UIcon name="i-tabler-category" class="size-5" />
                </div>
                <div class="flex-1 min-w-0">
                  <h3 class="font-semibold text-highlighted text-base truncate group-hover:text-primary transition-colors">
                    {{ category.name }}
                  </h3>
                  <p class="text-sm text-muted">{{ $t('site.categories.post_count', { n: category.count }) }}</p>
                </div>
              </div>

              <p v-if="category.description" class="text-sm text-muted mb-3 line-clamp-2">
                {{ category.description }}
              </p>

              <div class="flex items-center justify-between pt-3 border-t border-default mt-auto">
                <span class="text-xs text-muted">{{ category.slug }}</span>
              </div>
            </div>
          </NuxtLink>
        </TransitionGroup>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center mt-6">
        <UPagination
          v-model:page="currentPage"
          :total="filteredCategories.length"
          :items-per-page="pageSize" />
      </div>

      <!-- Plugin slot: categories bottom -->
      <ClientOnly><ContributionSlot name="public:categories-bottom" /></ClientOnly>
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

const termApi = useTermApi()

const sortOptions = computed(() => [
  { label: t('site.categories.sort_name'), value: 'name' },
  { label: t('site.categories.sort_count'), value: 'count' },
])

const isLoading = ref(true)
const searchQuery = ref('')
const sortBy = ref<'name' | 'count'>('name')
const categories = ref<TermDetailResponse[]>([])
const currentPage = ref(1)
const pageSize = 12

const filteredCategories = computed(() => {
  let result = categories.value
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    result = result.filter(c =>
      c.name.toLowerCase().includes(q) || c.description.toLowerCase().includes(q)
    )
  }
  return [...result].sort((a, b) =>
    sortBy.value === 'name' ? a.name.localeCompare(b.name, 'zh-CN') : b.count - a.count
  )
})

const totalPages = computed(() => Math.ceil(filteredCategories.value.length / pageSize))

const paginatedCategories = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredCategories.value.slice(start, start + pageSize)
})

watch([searchQuery, sortBy], () => {
  currentPage.value = 1
})

onMounted(async () => {
  try {
    categories.value = await termApi.getTerms({ taxonomy: 'category' })
  } catch (e) {
    console.error('Failed to load categories:', e)
  } finally {
    isLoading.value = false
  }
})
</script>

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: scale(0.9);
}
</style>
