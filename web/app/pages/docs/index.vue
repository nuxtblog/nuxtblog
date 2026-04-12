<script setup lang="ts">
import type { DocCollectionItem } from '~/types/api/doc'

const docApi = useDocApi()

const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)
const collections = ref<DocCollectionItem[]>([])
const total = ref(0)

useHead({
  title: '文档',
  meta: [{ name: 'description', content: '浏览所有文档合集' }],
})

async function fetchCollections() {
  rawLoading.value = true
  try {
    const res = await docApi.getCollections({ status: 2, page_size: 50 })
    collections.value = res.data ?? []
    total.value = res.total ?? 0
  } catch (e: any) {
    // silently fail on public page
  } finally {
    rawLoading.value = false
  }
}

onMounted(fetchCollections)
</script>

<template>
  <div class="min-h-screen pb-16">
    <PageHeader
      icon="i-tabler-books"
      title="文档"
      :subtitle="`共 ${total} 个合集`" />

    <!-- Plugin slot: docs top -->
    <ClientOnly><ContributionSlot name="public:docs-top" /></ClientOnly>

    <PageContent>
      <!-- Skeleton -->
      <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="i in 6" :key="i" class="rounded-md bg-default ring-1 ring-default shadow-sm p-5 space-y-3">
          <USkeleton class="h-5 w-3/4" />
          <USkeleton class="h-4 w-full" />
          <USkeleton class="h-4 w-2/3" />
          <div class="flex gap-2">
            <USkeleton class="h-5 w-12 rounded-full" />
            <USkeleton class="h-5 w-10 rounded-full" />
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="collections.length === 0" class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">
        <div class="py-20 flex flex-col items-center gap-3 text-center px-4">
          <div class="size-16 rounded-full bg-muted flex items-center justify-center">
            <UIcon name="i-tabler-books-off" class="size-8 text-muted" />
          </div>
          <p class="font-semibold text-highlighted">暂无文档</p>
          <p class="text-sm text-muted">还没有发布任何文档合集</p>
        </div>
      </div>

      <!-- Collection grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <NuxtLink
          v-for="col in collections" :key="col.id"
          :to="`/docs/${col.slug}`"
          class="rounded-md bg-default ring-1 ring-default shadow-sm p-5 hover:shadow-md transition-all cursor-pointer group">
          <h2 class="font-semibold text-highlighted group-hover:text-primary transition-colors mb-2">
            {{ col.title }}
          </h2>
          <p v-if="col.description" class="text-sm text-muted line-clamp-2 mb-3">{{ col.description }}</p>
          <div class="flex items-center gap-2 flex-wrap">
            <UBadge v-if="col.doc_count != null" :label="`${col.doc_count} 篇`" color="neutral" variant="soft" size="xs" />
            <UBadge v-if="col.locale" :label="col.locale.toUpperCase()" color="neutral" variant="outline" size="xs" />
          </div>
        </NuxtLink>
      </div>
    </PageContent>

    <!-- Plugin slot: docs bottom -->
    <ClientOnly><ContributionSlot name="public:docs-bottom" /></ClientOnly>
  </div>
</template>
