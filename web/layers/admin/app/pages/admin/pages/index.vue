<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.pages.title')" :subtitle="$t('admin.pages.subtitle')">
      <template #actions>
        <UButton as="NuxtLink" to="/admin/pages/new" icon="i-tabler-plus" color="primary">
          {{ $t('admin.pages.new_page') }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- 筛选工具栏 -->
      <div class="flex flex-wrap items-center gap-3 mb-4">
        <UInput
          v-model="searchQuery"
          :placeholder="$t('admin.pages.search_placeholder')"
          leading-icon="i-tabler-search"
          class="w-56"
          size="sm">
          <template v-if="searchQuery" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="searchQuery = ''" />
          </template>
        </UInput>
        <USelect
          v-model="filterStatus"
          :items="statusFilters.map(s => ({ label: s.label, value: s.value }))"
          class="w-36"
          size="sm"
        />
      </div>

      <!-- 加载骨架 -->
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 6" :key="i" class="flex items-center gap-4 p-4 border border-default rounded-md">
          <USkeleton class="size-4 rounded shrink-0" />
          <div class="flex-1 space-y-1.5">
            <USkeleton class="h-4 w-2/3" />
            <USkeleton class="h-3 w-1/3" />
          </div>
          <USkeleton class="h-5 w-14 rounded-full shrink-0" />
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else-if="!pages.length" class="flex flex-col items-center justify-center py-16">
        <UIcon name="i-tabler-file-off" class="size-16 text-muted mb-4" />
        <h3 class="text-lg font-medium text-highlighted mb-1">{{ searchQuery ? $t('admin.pages.no_results') : $t('admin.pages.no_pages') }}</h3>
        <p class="text-sm text-muted mb-4">{{ searchQuery ? $t('admin.pages.try_modify_filters') : $t('admin.pages.start_creating') }}</p>
        <UButton v-if="!searchQuery" as="NuxtLink" to="/admin/pages/new" icon="i-tabler-plus" color="primary">
          {{ $t('admin.pages.new_page') }}
        </UButton>
      </div>

      <!-- 列表 -->
      <div v-else class="space-y-3">
        <div
          v-for="page in pages"
          :key="page.id"
          class="flex items-center group gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all">
          <UCheckbox
            :model-value="selected.includes(page.id)"
            @update:model-value="toggleSelect(page.id)" />

          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-4">
              <div class="flex-1 min-w-0">
                <h3 class="text-sm font-medium text-highlighted group-hover:text-primary transition-colors truncate">
                  {{ page.title }}
                </h3>
                <div class="flex items-center gap-4 mt-1.5 text-xs text-muted">
                  <span>{{ page.author?.display_name || page.author?.username || '—' }}</span>
                  <span>{{ formatDate(page.updated_at) }}</span>
                </div>
              </div>

              <!-- 状态 + 操作 -->
              <div class="flex items-center gap-3 shrink-0">
                <UBadge :label="statusLabel(page.status)" :color="statusColor(page.status)" variant="soft" size="sm" />
                <UDropdownMenu :items="getPageActions(page)" :popper="{ placement: 'bottom-end' }">
                  <UButton
                    color="neutral"
                    variant="ghost"
                    icon="i-tabler-dots-vertical"
                    square
                    size="xs"
                    class="opacity-0 group-hover:opacity-100 transition-opacity" />
                </UDropdownMenu>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 删除确认 -->
      <UModal v-model:open="deleteModal">
        <template #content>
          <div class="p-6">
            <div class="flex items-center gap-3 mb-4">
              <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
                <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
              </div>
              <div>
                <h3 class="font-semibold text-highlighted">{{ $t('common.confirm_delete') }}</h3>
                <p class="text-sm text-muted mt-0.5">
                  {{ $t('admin.pages.delete_confirm', { title: deleteTarget?.title }) }}
                </p>
              </div>
            </div>
            <div class="flex justify-end gap-2 mt-6">
              <UButton color="neutral" variant="outline" @click="deleteModal = false">{{ $t('common.cancel') }}</UButton>
              <UButton color="error" :loading="deleteLoading" @click="doDelete">{{ $t('common.delete') }}</UButton>
            </div>
          </div>
        </template>
      </UModal>
    </AdminPageContent>

    <AdminPageFooter>
      <template #left>
        <template v-if="pages.length > 0">
          <UCheckbox
            :model-value="allSelected"
            :indeterminate="selected.length > 0 && !allSelected"
            @update:model-value="toggleAll" />
          <template v-if="selected.length > 0">
            <span>{{ $t('common.selected_n', { n: selected.length }) }}</span>
            <USeparator orientation="vertical" class="h-4" />
            <USelect
              v-model="batchAction"
              :items="[
                { label: $t('admin.pages.batch_publish'), value: 'publish' },
                { label: $t('admin.pages.batch_draft'), value: 'draft' },
                { label: $t('admin.pages.batch_delete'), value: 'delete' },
              ]"
              :placeholder="$t('admin.pages.batch_ops')"
              class="w-36"
              size="sm" />
            <UButton color="primary" variant="soft" size="sm" :disabled="!batchAction" :loading="batchLoading" @click="applyBatch">{{ $t('common.apply') }}</UButton>
            <UButton color="neutral" variant="ghost" size="sm" @click="selected = []; batchAction = undefined">{{ $t('common.cancel') }}</UButton>
          </template>
          <span v-else class="text-xs">{{ $t('common.selectAll') }}</span>
        </template>
      </template>
      <template #right>
        <UPagination v-if="total > pageSize" v-model:page="currentPage" :total="total" :items-per-page="pageSize" size="sm" />
      </template>
    </AdminPageFooter>
  </AdminPageContainer>
</template>

<script setup lang="ts">
const { apiFetch } = useApiFetch()
const authStore = useAuthStore()
const toast = useToast()
const { t } = useI18n()

const authHeaders = computed(() =>
  authStore.token ? { Authorization: `Bearer ${authStore.token}` } : {}
)

interface PageItem {
  id: number
  title: string
  slug: string
  status: number
  updated_at: string
  author?: { id: number; username: string; display_name: string }
}

interface PageListRes {
  data: PageItem[]
  total: number
}

const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)
const pages = ref<PageItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = 20
const filterStatus = ref('all')
const searchQuery = ref('')
const selected = ref<number[]>([])
const batchAction = ref<string | undefined>(undefined)
const batchLoading = ref(false)
const deleteModal = ref(false)
const deleteTarget = ref<PageItem | null>(null)
const deleteLoading = ref(false)

const statusFilters = computed(() => [
  { label: t('admin.pages.status_all'), value: 'all' },
  { label: t('admin.pages.status_published'), value: 'published' },
  { label: t('admin.pages.status_draft'), value: 'draft' },
])

const statusMap: Record<string, number> = { draft: 1, published: 2, private: 3, archived: 4 }

const fetchPages = async () => {
  try {
    const params: Record<string, any> = { page: currentPage.value, page_size: pageSize, post_type: 2 }
    if (filterStatus.value !== 'all') params.status = statusMap[filterStatus.value]
    if (searchQuery.value) params.search = searchQuery.value
    const res = await apiFetch<PageListRes>('/posts', { params })
    pages.value = res.data ?? []
    total.value = res.total ?? 0
  } catch {
    pages.value = []
    total.value = 0
  } finally {
    rawLoading.value = false
  }
}

watch(filterStatus, () => {
  currentPage.value = 1
  selected.value = []
  batchAction.value = undefined
  fetchPages()
})
watch(currentPage, fetchPages)
watchDebounced(searchQuery, () => {
  currentPage.value = 1
  selected.value = []
  batchAction.value = undefined
  fetchPages()
}, { debounce: 300 })
onMounted(fetchPages)

// Selection
const allSelected = computed(() => pages.value.length > 0 && pages.value.every(p => selected.value.includes(p.id)))
const toggleAll = (v: boolean) => { selected.value = v ? pages.value.map(p => p.id) : [] }
const toggleSelect = (id: number) => {
  const i = selected.value.indexOf(id)
  i > -1 ? selected.value.splice(i, 1) : selected.value.push(id)
}

// Batch
const applyBatch = async () => {
  if (!batchAction.value || !selected.value.length) return
  batchLoading.value = true
  try {
    if (batchAction.value === 'delete') {
      await Promise.all(selected.value.map(id =>
        apiFetch(`/posts/${id}`, { method: 'DELETE', headers: authHeaders.value })
      ))
      toast.add({ title: t('admin.pages.batch_deleted', { n: selected.value.length }), color: 'success' })
    } else {
      const status = batchAction.value === 'publish' ? 2 : 1
      await Promise.all(selected.value.map(id =>
        apiFetch(`/posts/${id}`, { method: 'PUT', body: { status }, headers: authHeaders.value })
      ))
      toast.add({ title: t('admin.pages.batch_success'), color: 'success' })
    }
    selected.value = []
    batchAction.value = undefined
    fetchPages()
  } catch {
    toast.add({ title: t('common.operation_failed'), color: 'error' })
  } finally {
    batchLoading.value = false
  }
}

// Actions menu
const getPageActions = (page: PageItem) => [
  [
    { label: t('common.edit'), icon: 'i-tabler-pencil', to: `/admin/pages/edit/${page.id}` },
    { label: t('common.view'), icon: 'i-tabler-eye', onClick: () => window.open(`/pages/${page.slug}`, '_blank') },
  ],
  [{ label: t('common.delete'), icon: 'i-tabler-trash', color: 'error' as const, onClick: () => confirmDelete(page) }],
]

// Delete
const confirmDelete = (page: PageItem) => { deleteTarget.value = page; deleteModal.value = true }
const doDelete = async () => {
  if (!deleteTarget.value) return
  deleteLoading.value = true
  try {
    await apiFetch(`/posts/${deleteTarget.value.id}`, { method: 'DELETE', headers: authHeaders.value })
    toast.add({ title: t('admin.pages.deleted'), color: 'success' })
    deleteModal.value = false
    fetchPages()
  } catch {
    toast.add({ title: t('admin.pages.delete_failed'), color: 'error' })
  } finally {
    deleteLoading.value = false
  }
}

// Helpers
const statusLabel = (s: number) => ({ 1: t('admin.pages.status_draft'), 2: t('admin.pages.status_published'), 3: t('admin.pages.status_private'), 4: t('admin.pages.status_archived') }[s] ?? t('common.unknown'))
type BadgeColor = 'success' | 'neutral' | 'warning'
const statusColor = (s: number): BadgeColor => ({ 1: 'neutral', 2: 'success', 3: 'warning' }[s] as BadgeColor ?? 'neutral')

const formatDate = (s: string) => {
  const d = Math.floor((Date.now() - new Date(s).getTime()) / 86400000)
  if (d === 0) return t('common.today')
  if (d === 1) return t('common.yesterday')
  if (d < 30) return t('common.days_ago', { n: d })
  return new Date(s).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })
}
</script>
