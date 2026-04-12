<template>
  <div class="min-h-screen pb-16">
    <PageHeader
      icon="i-tabler-stack-2"
      :title="$t('site.categories.title')"
      :subtitle="$t('site.categories.subtitle', { n: filteredCategories.length })">
      <template #actions>
        <div class="flex items-center gap-1 p-1 bg-default rounded-md ring-1 ring-default">
          <UButton
            v-for="opt in sortOptions"
            :key="opt.value"
            :variant="sortBy === opt.value ? 'solid' : 'ghost'"
            :color="sortBy === opt.value ? 'primary' : 'neutral'"
            size="xs"
            :class="sortBy === opt.value ? 'shadow-sm' : 'text-muted'"
            @click="sortBy = opt.value as 'name' | 'count'">
            {{ opt.label }}
          </UButton>
        </div>
      </template>
      <template #toolbar>
        <UInput
          v-model="searchQuery"
          :placeholder="$t('site.categories.search_placeholder')"
          leading-icon="i-tabler-search"
          class="w-full"
          size="md">
          <template v-if="searchQuery" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="searchQuery = ''" />
          </template>
        </UInput>
      </template>
    </PageHeader>

    <!-- Plugin slot: categories top -->
    <ClientOnly><ContributionSlot name="public:categories-top" /></ClientOnly>

    <PageContent>
      <!-- Loading skeleton -->
      <div v-if="loading" class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3 sm:gap-4">
        <div v-for="i in 8" :key="i" class="rounded-md bg-default ring-1 ring-default shadow-sm">
          <div class="p-5">
            <div class="flex items-start gap-3 mb-3">
              <USkeleton class="size-11 rounded-md shrink-0" />
              <div class="flex-1 space-y-2">
                <USkeleton class="h-4 w-24" />
                <USkeleton class="h-3 w-16" />
              </div>
            </div>
            <USkeleton class="h-3 w-full" />
            <USkeleton class="h-3 w-3/4 mt-1.5" />
            <div class="flex items-center justify-between pt-3 mt-3 border-t border-default">
              <USkeleton class="h-3 w-20" />
              <USkeleton class="size-4 rounded-full" />
            </div>
          </div>
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

      <!-- Category Cards Grid -->
      <div v-else class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-4 gap-3 sm:gap-4">
        <TransitionGroup name="list" tag="div" class="contents">
          <NuxtLink
            v-for="(category, index) in paginatedCategories"
            :key="category.term_taxonomy_id"
            :to="`/categories/${category.slug}`"
            class="rounded-md ring-1 ring-default bg-default shadow-sm hover:shadow-md hover:-translate-y-0.5 transition-all duration-200 group">
            <div class="p-5">
              <div class="flex items-start gap-3 mb-3">
                <div
                  class="size-11 rounded-md flex items-center justify-center shrink-0 transition-transform duration-200 group-hover:scale-110"
                  :style="{ backgroundColor: getColor(index) + '18', color: getColor(index) }">
                  <UIcon name="i-tabler-category" class="size-5" />
                </div>
                <div class="flex-1 min-w-0">
                  <h3 class="font-semibold text-highlighted text-base truncate group-hover:text-primary transition-colors">
                    {{ category.name }}
                  </h3>
                  <UBadge
                    :label="$t('site.categories.post_count', { n: category.count })"
                    variant="subtle"
                    color="neutral"
                    size="xs"
                    class="mt-1" />
                </div>
              </div>

              <p v-if="category.description" class="text-sm text-muted line-clamp-2 leading-relaxed">
                {{ category.description }}
              </p>

              <div class="flex items-center justify-between pt-3 mt-3 border-t border-default">
                <span class="text-xs text-muted font-mono">{{ category.slug }}</span>
                <UIcon
                  name="i-tabler-arrow-right"
                  class="size-4 text-muted group-hover:text-primary group-hover:translate-x-0.5 transition-all duration-200" />
              </div>
            </div>
          </NuxtLink>
        </TransitionGroup>
      </div>
    </PageContent>

    <!-- Pagination -->
    <PageFooter v-if="totalPages > 1">
      <div class="flex justify-center">
        <UPagination
          v-model:page="currentPage"
          :total="filteredCategories.length"
          :items-per-page="pageSize" />
      </div>
    </PageFooter>

    <!-- Plugin slot: categories bottom -->
    <ClientOnly><ContributionSlot name="public:categories-bottom" /></ClientOnly>
  </div>
</template>

<script setup lang="ts">
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

const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)
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
    rawLoading.value = false
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
