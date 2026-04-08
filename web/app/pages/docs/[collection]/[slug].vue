<script setup lang="ts">
import type { DocDetailItem, DocItem, DocCollectionItem } from '~/types/api/doc'

const route = useRoute()
const docApi = useDocApi()
const { containerClass } = useContainerWidth()

const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)

const doc = ref<DocDetailItem | null>(null)
const collection = ref<DocCollectionItem | null>(null)
const collectionDocs = ref<DocItem[]>([])

const collectionSlug = computed(() => route.params.collection as string)
const docSlug = computed(() => route.params.slug as string)


useHead(computed(() => ({
  title: doc.value ? `${doc.value.title} - ${collection.value?.title ?? '文档'}` : '文档',
  meta: [
    { name: 'description', content: doc.value?.seo?.meta_desc || doc.value?.excerpt || '' },
    { property: 'og:title', content: doc.value?.seo?.og_title || doc.value?.title || '' },
    { property: 'og:image', content: doc.value?.seo?.og_image || '' },
  ],
  bodyAttrs: {
    'data-page-id': doc.value?.id ? String(doc.value.id) : '',
    'data-page-slug': doc.value?.slug || '',
    'data-page-collection': collectionSlug.value,
  },
})))

// Top-level docs for sidebar nav
const topLevelDocs = computed(() => collectionDocs.value.filter(d => !d.parent_id))
const childDocsByParent = computed(() => {
  const map: Record<number, DocItem[]> = {}
  collectionDocs.value.forEach(d => {
    if (d.parent_id) {
      if (!map[d.parent_id]) map[d.parent_id] = []
      map[d.parent_id]!.push(d)
    }
  })
  return map
})

// Prev/next navigation among top-level docs
const currentDocIndex = computed(() => topLevelDocs.value.findIndex(d => d.slug === docSlug.value))
const prevDoc = computed(() => currentDocIndex.value > 0 ? topLevelDocs.value[currentDocIndex.value - 1] : null)
const nextDoc = computed(() => currentDocIndex.value < topLevelDocs.value.length - 1 ? topLevelDocs.value[currentDocIndex.value + 1] : null)

async function fetchDoc() {
  rawLoading.value = true
  try {
    const d = await docApi.getDocBySlug(docSlug.value)
    doc.value = d

    // Fetch collection info
    const colRes = await docApi.getCollections({ page_size: 100 })
    const found = (colRes.data ?? []).find(c => c.slug === collectionSlug.value)
    if (found) {
      collection.value = found
      // Fetch all docs in collection for sidebar
      const docsRes = await docApi.getDocs({ collection_id: found.id, status: 2, page_size: 100 })
      collectionDocs.value = docsRes.data ?? []
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

onMounted(async () => {
  await fetchDoc()
  // Increment view count
  if (doc.value) {
    docApi.incrementView(doc.value.id).catch(() => {})
  }
})

// Re-fetch when slug changes
watch(docSlug, fetchDoc)
</script>

<template>
  <div class="min-h-screen pb-16">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-8']">
      <!-- Skeleton -->
      <div v-if="loading" class="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <div class="lg:col-span-1 space-y-2">
          <USkeleton v-for="i in 8" :key="i" class="h-8 w-full" />
        </div>
        <div class="lg:col-span-3 space-y-4">
          <USkeleton class="h-8 w-3/4" />
          <USkeleton class="h-4 w-1/2" />
          <div class="space-y-3 mt-6">
            <USkeleton v-for="i in 10" :key="i" class="h-4 w-full" />
          </div>
        </div>
      </div>

      <div v-else-if="!doc" class="flex flex-col items-center justify-center py-20">
        <div class="size-16 rounded-full bg-muted flex items-center justify-center mb-4">
          <UIcon name="i-tabler-file-off" class="size-8 text-muted" />
        </div>
        <h3 class="text-lg font-medium text-highlighted mb-1">文档不存在</h3>
        <NuxtLink :to="`/docs/${collectionSlug}`" class="text-sm text-primary hover:underline">返回合集</NuxtLink>
      </div>

      <div v-else class="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <!-- Left sidebar nav -->
        <aside class="hidden lg:block lg:col-span-1">
          <div class="sticky top-6">
            <NuxtLink :to="`/docs/${collectionSlug}`" class="flex items-center gap-1.5 text-sm text-muted hover:text-primary transition-colors mb-3">
              <UIcon name="i-tabler-arrow-left" class="size-4" />
              {{ collection?.title ?? '返回合集' }}
            </NuxtLink>

            <nav class="rounded-lg bg-default ring-1 ring-default shadow-sm p-3">
              <ul class="space-y-0.5">
                <template v-for="item in topLevelDocs" :key="item.id">
                  <li>
                    <NuxtLink
                      :to="`/docs/${collectionSlug}/${item.slug}`"
                      class="flex items-center gap-1.5 px-2 py-1.5 rounded-md text-sm transition-colors hover:bg-elevated"
                      :class="item.slug === docSlug ? 'bg-primary/10 text-primary font-medium' : 'text-muted'">
                      <UIcon name="i-tabler-file-text" class="size-3.5 shrink-0" />
                      <span class="truncate">{{ item.title }}</span>
                    </NuxtLink>
                    <ul v-if="childDocsByParent[item.id]?.length" class="ml-4 mt-0.5 space-y-0.5 border-l border-default pl-2">
                      <li v-for="child in childDocsByParent[item.id]" :key="child.id">
                        <NuxtLink
                          :to="`/docs/${collectionSlug}/${child.slug}`"
                          class="flex items-center gap-1.5 px-2 py-1.5 rounded-md text-sm transition-colors hover:bg-elevated"
                          :class="child.slug === docSlug ? 'text-primary font-medium' : 'text-muted'">
                          <span class="truncate">{{ child.title }}</span>
                        </NuxtLink>
                      </li>
                    </ul>
                  </li>
                </template>
              </ul>
            </nav>
          </div>
        </aside>

        <!-- Main doc content -->
        <main class="lg:col-span-3">
          <!-- Breadcrumb (mobile) -->
          <div class="flex items-center gap-1.5 text-sm text-muted mb-4 lg:hidden">
            <NuxtLink to="/docs" class="hover:text-primary">文档</NuxtLink>
            <UIcon name="i-tabler-chevron-right" class="size-3.5" />
            <NuxtLink :to="`/docs/${collectionSlug}`" class="hover:text-primary">{{ collection?.title }}</NuxtLink>
            <UIcon name="i-tabler-chevron-right" class="size-3.5" />
            <span class="text-highlighted truncate">{{ doc.title }}</span>
          </div>

          <!-- Doc header -->
          <div class="mb-6">
            <h1 class="text-2xl font-bold text-highlighted mb-2">{{ doc.title }}</h1>
            <div class="flex items-center gap-4 text-sm text-muted">
              <span v-if="doc.updated_at" class="flex items-center gap-1">
                <UIcon name="i-tabler-calendar" class="size-4" />
                {{ formatDate(doc.updated_at) }}
              </span>
              <span v-if="doc.stats?.view_count" class="flex items-center gap-1">
                <UIcon name="i-tabler-eye" class="size-4" />
                {{ doc.stats.view_count }}
              </span>
            </div>
          </div>

          <!-- Doc content -->
          <MarkdownContent v-if="doc.content" :content="doc.content" />
          <div v-else class="text-muted text-sm italic">此文档暂无内容。</div>

          <!-- Prev/Next navigation -->
          <div class="mt-10 pt-6 border-t border-default grid grid-cols-2 gap-4">
            <NuxtLink
              v-if="prevDoc"
              :to="`/docs/${collectionSlug}/${prevDoc.slug}`"
              class="flex items-center gap-2 p-3 rounded-lg ring-1 ring-default hover:ring-primary/50 hover:shadow-sm transition-all group">
              <UIcon name="i-tabler-arrow-left" class="size-4 text-muted group-hover:text-primary shrink-0" />
              <div class="min-w-0">
                <p class="text-xs text-muted">上一篇</p>
                <p class="text-sm font-medium text-highlighted group-hover:text-primary truncate">{{ prevDoc.title }}</p>
              </div>
            </NuxtLink>
            <div v-else />

            <NuxtLink
              v-if="nextDoc"
              :to="`/docs/${collectionSlug}/${nextDoc.slug}`"
              class="flex items-center justify-end gap-2 p-3 rounded-lg ring-1 ring-default hover:ring-primary/50 hover:shadow-sm transition-all group text-right">
              <div class="min-w-0">
                <p class="text-xs text-muted">下一篇</p>
                <p class="text-sm font-medium text-highlighted group-hover:text-primary truncate">{{ nextDoc.title }}</p>
              </div>
              <UIcon name="i-tabler-arrow-right" class="size-4 text-muted group-hover:text-primary shrink-0" />
            </NuxtLink>
          </div>
        </main>
      </div>
    </div>
  </div>
</template>
