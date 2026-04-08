<script setup lang="ts">
const reactionApi = useReactionApi()

const page = ref(1)
const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)
const bookmarks = ref<Array<{ id: number; title: string; slug: string; excerpt: string; created_at: string }>>([])
const total = ref(0)

const formatDate = (iso: string) =>
  iso ? new Date(iso).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' }) : ''

async function load() {
  rawLoading.value = true
  try {
    const res = await reactionApi.getBookmarks(page.value, 10)
    bookmarks.value = res.list
    total.value = res.total
  } finally {
    rawLoading.value = false
  }
}

watch(page, load)
onMounted(load)
</script>

<template>
  <BaseCardList :loading="loading" :empty="bookmarks.length === 0">

    <template #skeleton>
      <div class="flex items-start gap-3">
        <USkeleton class="size-4 rounded shrink-0 mt-0.5" />
        <div class="flex-1 space-y-2">
          <USkeleton class="h-4 w-3/4" />
          <USkeleton class="h-3 w-full" />
          <USkeleton class="h-3 w-1/3" />
        </div>
      </div>
    </template>

    <template #empty>
      <div class="text-center py-12 text-muted">
        <UIcon name="i-tabler-bookmark-off" class="size-12 mx-auto mb-3" />
        <p>{{ $t('site.user.no_favorites') }}</p>
        <UButton to="/" color="primary" variant="soft" size="sm" class="mt-4">
          {{ $t('site.user.discover_articles') }}
        </UButton>
      </div>
    </template>

    <NuxtLink
      v-for="item in bookmarks"
      :key="item.id"
      :to="`/posts/${item.slug}`"
      class="group block">
      <UCard class="hover:shadow-md transition-shadow">
        <div class="flex items-start gap-3">
          <UIcon name="i-tabler-bookmark-filled" class="size-4 text-primary shrink-0 mt-0.5" />
          <div class="flex-1 min-w-0">
            <h3 class="font-medium text-highlighted group-hover:text-primary transition-colors line-clamp-1">
              {{ item.title }}
            </h3>
            <p class="text-muted text-sm line-clamp-2 mt-0.5">{{ item.excerpt }}</p>
            <p class="text-xs text-muted mt-1">{{ formatDate(item.created_at) }}</p>
          </div>
        </div>
      </UCard>
    </NuxtLink>

    <div v-if="total > 10" class="flex justify-center pt-2">
      <UPagination v-model:page="page" :total="total" :items-per-page="10" />
    </div>

  </BaseCardList>
</template>
