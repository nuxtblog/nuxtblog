<script setup lang="ts">
import type { DocCollectionItem, DocItem } from '~/types/api/doc'
import { useDebounceFn } from '@vueuse/core'

const route = useRoute()
const docApi = useDocApi()

const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)

const collection = ref<DocCollectionItem | null>(null)
const docs = ref<DocItem[]>([])

const searchQuery = ref('')
const searchResults = ref<DocItem[] | null>(null)
const searching = ref(false)

const collectionSlug = computed(() => route.params.collection as string)

useHead(computed(() => ({
  title: collection.value?.title ?? '文档合集',
  meta: [{ name: 'description', content: collection.value?.description ?? '' }],
})))

// Group docs by parent_id
const topLevelDocs = computed(() => docs.value.filter(d => !d.parent_id))
const childDocsByParent = computed(() => {
  const map: Record<number, DocItem[]> = {}
  docs.value.forEach(d => {
    if (d.parent_id) {
      if (!map[d.parent_id]) map[d.parent_id] = []
      map[d.parent_id]!.push(d)
    }
  })
  return map
})

async function fetchCollection() {
  rawLoading.value = true
  try {
    // Get collections and find by slug
    const res = await docApi.getCollections({ status: 2, page_size: 100 })
    const found = (res.data ?? []).find(c => c.slug === collectionSlug.value)
    if (found) {
      collection.value = found
      // Fetch docs in this collection
      const docsRes = await docApi.getDocs({ collection_id: found.id, status: 2, page_size: 100 })
      docs.value = docsRes.data ?? []
    }
  } catch (e: any) {
    // silently fail
  } finally {
    rawLoading.value = false
  }
}

function formatDate(d: string) {
  return d ? new Date(d).toLocaleDateString('zh-CN') : ''
}

const doSearch = useDebounceFn(async () => {
  const kw = searchQuery.value.trim()
  if (!kw) {
    searchResults.value = null
    return
  }
  if (!collection.value) return
  searching.value = true
  try {
    const res = await docApi.getDocs({ collection_id: collection.value.id, keyword: kw, status: 2, page_size: 50 })
    searchResults.value = res.data ?? []
  } catch {
    searchResults.value = []
  } finally {
    searching.value = false
  }
}, 350)

watch(searchQuery, doSearch)

onMounted(fetchCollection)
</script>

<template>
  <div class="min-h-screen pb-16">
    <!-- Skeleton -->
    <PageContent v-if="loading">
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <div class="lg:col-span-1 space-y-3">
          <USkeleton class="h-6 w-3/4" />
          <USkeleton class="h-4 w-full" />
          <div class="space-y-2 mt-4">
            <USkeleton v-for="i in 6" :key="i" class="h-8 w-full" />
          </div>
        </div>
        <div class="lg:col-span-3 space-y-3">
          <div v-for="i in 5" :key="i" class="rounded-lg bg-default ring-1 ring-default p-4 space-y-2">
            <USkeleton class="h-5 w-1/2" />
            <USkeleton class="h-4 w-full" />
          </div>
        </div>
      </div>
    </PageContent>

    <PageContent v-else-if="!collection">
      <div class="flex flex-col items-center justify-center py-20">
        <div class="size-16 rounded-full bg-muted flex items-center justify-center mb-4">
          <UIcon name="i-tabler-books-off" class="size-8 text-muted" />
        </div>
        <h3 class="text-lg font-medium text-highlighted mb-1">合集不存在</h3>
        <NuxtLink to="/docs" class="text-sm text-primary hover:underline">返回文档列表</NuxtLink>
      </div>
    </PageContent>

    <template v-else>
      <PageHeader
        icon="i-tabler-books"
        :title="collection.title"
        :subtitle="collection.description">
      </PageHeader>

      <PageContent>
        <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
          <!-- Left sidebar -->
          <aside class="lg:col-span-1">
            <!-- Collection info -->
            <div class="rounded-lg bg-default ring-1 ring-default shadow-sm p-4 mb-4">
              <div class="flex items-center gap-2 mb-2">
                <UIcon name="i-tabler-books" class="size-5 text-primary shrink-0" />
                <h2 class="font-semibold text-highlighted truncate">{{ collection.title }}</h2>
              </div>
              <p v-if="collection.description" class="text-sm text-muted line-clamp-3">{{ collection.description }}</p>
              <div v-if="collection.doc_count != null" class="mt-2">
                <UBadge :label="`${collection.doc_count} 篇文档`" color="neutral" variant="soft" size="xs" />
              </div>
            </div>

            <!-- Doc tree nav -->
            <nav class="rounded-lg bg-default ring-1 ring-default shadow-sm p-3">
              <h3 class="text-xs font-semibold text-muted uppercase tracking-wider mb-2 px-1">目录</h3>
              <DocNavTree
                :docs="topLevelDocs"
                :child-docs-by-parent="childDocsByParent"
                :collection-slug="collectionSlug"
                current-slug=""
              />
            </nav>
          </aside>

          <!-- Main content: doc list -->
          <main class="lg:col-span-3">
            <!-- Search -->
            <div class="mb-4">
              <UInput
                v-model="searchQuery"
                placeholder="搜索文档..."
                leading-icon="i-tabler-search"
                :loading="searching"
                class="w-full" />
            </div>

            <!-- Search results -->
            <template v-if="searchResults !== null">
              <div v-if="searchResults.length === 0" class="flex flex-col items-center justify-center py-16">
                <div class="size-16 rounded-full bg-muted flex items-center justify-center mb-4">
                  <UIcon name="i-tabler-search-off" class="size-8 text-muted" />
                </div>
                <p class="text-sm text-muted">未找到相关文档</p>
              </div>
              <div v-else class="rounded-lg bg-default ring-1 ring-default shadow-sm overflow-hidden">
                <div class="divide-y divide-default">
                  <NuxtLink
                    v-for="item in searchResults" :key="item.id"
                    :to="`/docs/${collectionSlug}/${item.slug}`"
                    class="flex items-start gap-4 p-4 hover:bg-elevated transition-colors group">
                    <div class="flex-1 min-w-0">
                      <h3 class="font-medium text-highlighted group-hover:text-primary transition-colors truncate">
                        {{ item.title }}
                      </h3>
                      <p v-if="item.excerpt" class="text-sm text-muted line-clamp-2 mt-0.5">{{ item.excerpt }}</p>
                    </div>
                    <UIcon name="i-tabler-chevron-right" class="size-4 text-muted shrink-0 mt-1 group-hover:text-primary transition-colors" />
                  </NuxtLink>
                </div>
              </div>
            </template>

            <!-- Normal tree list -->
            <template v-else>
              <div v-if="topLevelDocs.length === 0" class="flex flex-col items-center justify-center py-16">
                <div class="size-16 rounded-full bg-muted flex items-center justify-center mb-4">
                  <UIcon name="i-tabler-file-text-off" class="size-8 text-muted" />
                </div>
                <p class="text-sm text-muted">此合集暂无文档</p>
              </div>

              <div v-else class="rounded-lg bg-default ring-1 ring-default shadow-sm overflow-hidden">
                <div class="divide-y divide-default">
                  <NuxtLink
                    v-for="doc in topLevelDocs" :key="doc.id"
                    :to="`/docs/${collectionSlug}/${doc.slug}`"
                    class="flex items-start gap-4 p-4 hover:bg-elevated transition-colors group">
                    <div class="flex-1 min-w-0">
                      <h3 class="font-medium text-highlighted group-hover:text-primary transition-colors truncate">
                        {{ doc.title }}
                      </h3>
                      <p v-if="doc.excerpt" class="text-sm text-muted line-clamp-2 mt-0.5">{{ doc.excerpt }}</p>
                      <div class="flex items-center gap-3 mt-2">
                        <span v-if="doc.stats?.view_count" class="flex items-center gap-1 text-xs text-muted">
                          <UIcon name="i-tabler-eye" class="size-3.5" />
                          {{ doc.stats.view_count }}
                        </span>
                        <span v-if="doc.updated_at" class="text-xs text-muted">{{ formatDate(doc.updated_at) }}</span>
                        <UBadge v-if="childDocsByParent[doc.id]?.length" :label="`${childDocsByParent[doc.id]?.length ?? 0} 子文档`" color="neutral" variant="soft" size="xs" />
                      </div>
                    </div>
                    <UIcon name="i-tabler-chevron-right" class="size-4 text-muted shrink-0 mt-1 group-hover:text-primary transition-colors" />
                  </NuxtLink>
                </div>
              </div>
            </template>
          </main>
        </div>
      </PageContent>
    </template>
  </div>
</template>
