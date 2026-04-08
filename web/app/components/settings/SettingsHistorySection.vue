<template>
  <UCard>
    <template #header>
      <div class="flex items-center justify-between">
        <h2 class="font-semibold text-highlighted flex items-center gap-2">
          <UIcon name="i-tabler-history" class="size-5 text-primary" />
          浏览历史
        </h2>
        <UButton
          v-if="historyItems.length > 0"
          color="error" variant="soft" size="sm" icon="i-tabler-trash"
          @click="clearHistory">
          清空历史
        </UButton>
      </div>
    </template>

    <div v-if="historyLoading" class="space-y-3">
      <div v-for="i in 5" :key="i" class="flex gap-3 py-2">
        <USkeleton class="w-16 h-12 rounded shrink-0" />
        <div class="flex-1 space-y-1.5">
          <USkeleton class="h-4 w-3/4" />
          <USkeleton class="h-3 w-1/3" />
        </div>
      </div>
    </div>

    <div v-else-if="historyItems.length === 0" class="py-10 text-center">
      <UIcon name="i-tabler-history-off" class="size-12 text-muted mx-auto mb-3" />
      <p class="text-sm text-muted">暂无浏览记录</p>
    </div>

    <div v-else class="divide-y divide-default -my-1">
      <NuxtLink
        v-for="item in historyItems"
        :key="item.post_id"
        :to="`/posts/${item.post_slug}`"
        class="flex gap-3 py-3 hover:bg-muted/50 rounded-md px-2 -mx-2 transition-colors">
        <img
          v-if="item.post_cover"
          :src="item.post_cover"
          class="w-16 h-12 rounded object-cover shrink-0 bg-muted" />
        <div v-else class="w-16 h-12 rounded bg-muted shrink-0 flex items-center justify-center">
          <UIcon name="i-tabler-file-text" class="size-5 text-muted" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm font-medium text-highlighted truncate">{{ item.post_title }}</p>
          <p class="text-xs text-muted mt-0.5">{{ new Date(item.viewed_at).toLocaleString('zh-CN') }}</p>
        </div>
      </NuxtLink>
    </div>

    <template v-if="historyTotal > 20" #footer>
      <div class="flex justify-center">
        <UPagination v-model:page="historyPage" :total="historyTotal" :items-per-page="20" />
      </div>
    </template>
  </UCard>
</template>

<script setup lang="ts">
import type { HistoryItem } from '~/composables/useHistoryApi'

const historyApi = useHistoryApi()
const toast = useToast()

const historyItems = ref<HistoryItem[]>([])
const historyTotal = ref(0)
const historyPage = ref(1)
const historyLoading = ref(false)

async function loadHistory() {
  if (historyLoading.value) return
  historyLoading.value = true
  try {
    const res = await historyApi.list(historyPage.value)
    historyItems.value = res.items
    historyTotal.value = res.total
  } catch {} finally {
    historyLoading.value = false
  }
}

async function clearHistory() {
  await historyApi.clear()
  historyItems.value = []
  historyTotal.value = 0
  toast.add({ title: '浏览历史已清空', color: 'success' })
}

onMounted(loadHistory)
watch(historyPage, loadHistory)
</script>
