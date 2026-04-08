<script setup lang="ts">
definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const toast = useToast()
const docApi = useDocApi()

// loading
const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)

// state
const filterStatus = ref('all')
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const docs = ref<any[]>([])
const collections = ref<any[]>([])
const filterCollectionId = ref<number | undefined>(undefined)

// delete modal
const showDeleteModal = ref(false)
const deleteTarget = ref<any>(null)
const deleteLoading = ref(false)

// batch selection
const selectedDocs = ref<number[]>([])
const batchAction = ref<string | undefined>(undefined)

const isAllDocsSelected = computed(() =>
  docs.value.length > 0 && docs.value.every(d => selectedDocs.value.includes(d.id))
)
const isDocsIndeterminate = computed(() =>
  selectedDocs.value.length > 0 && !isAllDocsSelected.value
)

function toggleSelectDoc(id: number) {
  const idx = selectedDocs.value.indexOf(id)
  if (idx > -1) selectedDocs.value.splice(idx, 1)
  else selectedDocs.value.push(id)
}

function toggleSelectAllDocs() {
  if (isAllDocsSelected.value) selectedDocs.value = []
  else selectedDocs.value = docs.value.map(d => d.id)
}

async function applyBatchDocs() {
  if (!batchAction.value || !selectedDocs.value.length) return
  try {
    if (batchAction.value === 'delete') {
      await Promise.all(selectedDocs.value.map(id => docApi.deleteDoc(id)))
      toast.add({ title: t('common.deleted'), color: 'success' })
    } else {
      const status = Number(batchAction.value)
      await Promise.all(selectedDocs.value.map(id => docApi.updateDoc(id, { status })))
      toast.add({ title: t('common.updated'), color: 'success' })
    }
    selectedDocs.value = []
    batchAction.value = undefined
    fetchDocs()
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.operation_failed'), color: 'error' })
  }
}

// computed
const statusTabs = computed(() => [
  { value: 'all', label: t('admin.docs.tab_all'), count: total.value },
  { value: '2', label: t('admin.docs.tab_published') },
  { value: '1', label: t('admin.docs.tab_draft') },
  { value: '3', label: t('admin.docs.tab_archived') },
])

const collectionOptions = computed(() => [
  { label: t('admin.docs.all_collections'), value: undefined },
  ...collections.value.map(c => ({ label: c.title, value: c.id })),
])

const statusBadge: Record<number, { label: string; color: 'success' | 'neutral' | 'warning' }> = {
  1: { label: t('admin.docs.status_draft'), color: 'neutral' },
  2: { label: t('admin.docs.status_published'), color: 'success' },
  3: { label: t('admin.docs.status_archived'), color: 'warning' },
}

// methods
async function fetchDocs() {
  rawLoading.value = true
  try {
    const params: any = { page: currentPage.value, page_size: pageSize.value }
    if (filterStatus.value !== 'all') params.status = Number(filterStatus.value)
    if (filterCollectionId.value) params.collection_id = filterCollectionId.value
    if (searchKeyword.value) params.keyword = searchKeyword.value
    const res = await docApi.getDocs(params)
    docs.value = res.data ?? []
    total.value = res.total ?? 0
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.load_failed'), color: 'error' })
  } finally {
    rawLoading.value = false
  }
}

async function fetchCollections() {
  try {
    const res = await docApi.getCollections({ page_size: 100 })
    collections.value = res.data ?? []
  } catch {}
}

function onStatusTab(val: string) {
  filterStatus.value = val
  currentPage.value = 1
}

function openDeleteModal(item: any) {
  deleteTarget.value = item
  showDeleteModal.value = true
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  deleteLoading.value = true
  try {
    await docApi.deleteDoc(deleteTarget.value.id)
    toast.add({ title: t('common.deleted'), color: 'success' })
    showDeleteModal.value = false
    fetchDocs()
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.delete_failed'), color: 'error' })
  } finally {
    deleteLoading.value = false
  }
}

function formatDate(d: string) {
  return d ? new Date(d).toLocaleDateString('zh-CN') : '-'
}

function openDocPreview(doc: any) {
  const slug = getCollectionSlug(doc.collection_id)
  if (typeof window !== "undefined") window.open(`/docs/${slug}/${doc.slug}`, "_blank")
}

function getCollectionSlug(collectionId: number) {
  return collections.value.find(c => c.id === collectionId)?.slug ?? ''
}

function getItemActions(item: any) {
  return [
    [
      { label: t('common.edit'), icon: 'i-tabler-pencil', to: `/admin/docs/${item.id}/edit` },
      { label: t('admin.docs.view_doc'), icon: 'i-tabler-eye', onClick: () => openDocPreview(item) },
    ],
    [
      { label: t('common.delete'), icon: 'i-tabler-trash', color: 'error' as const, onClick: () => openDeleteModal(item) },
    ],
  ]
}

watch([filterStatus, filterCollectionId], () => {
  currentPage.value = 1
  selectedDocs.value = []
  batchAction.value = undefined
  fetchDocs()
})
watch(currentPage, fetchDocs)
watch(searchKeyword, useDebounceFn(() => {
  currentPage.value = 1
  selectedDocs.value = []
  batchAction.value = undefined
  fetchDocs()
}, 400))

onMounted(() => { fetchCollections(); fetchDocs() })
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.docs.title')" :subtitle="$t('admin.docs.subtitle')">
      <template #actions>
        <div class="flex items-center gap-2">
          <UButton as="NuxtLink" to="/admin/docs/collections" variant="outline" color="neutral" icon="i-tabler-folders" size="sm">
            {{ $t('admin.docs.collections') }}
          </UButton>
          <UButton as="NuxtLink" to="/admin/docs/new" icon="i-tabler-plus" color="primary">
            {{ $t('admin.docs.new_doc') }}
          </UButton>
        </div>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- Status Tabs -->
      <div class="flex items-center gap-1 border-b border-default pb-0 mb-4 overflow-x-auto">
        <button
          v-for="s in statusTabs" :key="s.value"
          class="px-3 py-2 text-sm font-medium rounded-t transition-colors whitespace-nowrap"
          :class="filterStatus === s.value ? 'text-primary border-b-2 border-primary' : 'text-muted hover:text-highlighted'"
          @click="onStatusTab(s.value)">
          {{ s.label }}
          <span v-if="s.count != null" class="ml-1 text-xs text-muted">({{ s.count }})</span>
        </button>
      </div>

      <!-- Toolbar -->
      <div class="flex flex-wrap items-center gap-3 mb-4">
        <UInput v-model="searchKeyword" :placeholder="$t('admin.docs.search_placeholder')"
          leading-icon="i-tabler-search" class="w-56" size="sm">
          <template v-if="searchKeyword" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="searchKeyword = ''" />
          </template>
        </UInput>
        <AdminSearchableSelect
          v-model="filterCollectionId"
          :items="collectionOptions"
          :placeholder="$t('admin.docs.all_collections')"
          :search-placeholder="$t('common.search')" />
      </div>

      <!-- Skeleton -->
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 8" :key="i" class="flex items-center gap-4 p-4 border border-default rounded-lg">
          <div class="flex-1 space-y-2">
            <USkeleton class="h-4 w-3/4" />
            <div class="flex gap-3">
              <USkeleton class="h-3 w-16 rounded-full" />
              <USkeleton class="h-3 w-12 rounded-full" />
              <USkeleton class="h-3 w-20" />
            </div>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="docs.length === 0" class="flex flex-col items-center justify-center py-16">
        <UIcon name="i-tabler-file-text-off" class="size-16 text-muted mb-4" />
        <h3 class="text-lg font-medium text-highlighted mb-1">{{ $t('admin.docs.no_docs') }}</h3>
        <p class="text-sm text-muted mb-4">{{ $t('admin.docs.no_docs_desc') }}</p>
        <UButton as="NuxtLink" to="/admin/docs/new" color="primary" icon="i-tabler-plus">
          {{ $t('admin.docs.new_doc') }}
        </UButton>
      </div>

      <!-- List -->
      <div v-else class="space-y-3">
        <div
          v-for="doc in docs" :key="doc.id"
          class="flex items-center gap-4 p-4 border border-default rounded-lg group hover:shadow-sm transition-all bg-default">
          <UCheckbox :model-value="selectedDocs.includes(doc.id)" @update:model-value="toggleSelectDoc(doc.id)" />
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-4">
              <div class="flex-1 min-w-0">
                <h3 class="text-sm font-medium text-highlighted truncate group-hover:text-primary transition-colors">
                  {{ doc.title }}
                </h3>
                <div class="flex items-center gap-3 mt-1.5 flex-wrap">
                  <span class="text-xs text-muted">{{ formatDate(doc.updated_at) }}</span>
                  <span v-if="doc.excerpt" class="text-xs text-muted truncate max-w-xs">{{ doc.excerpt }}</span>
                </div>
              </div>
              <div class="flex items-center gap-3 shrink-0">
                <UBadge
                  v-if="statusBadge[doc.status]"
                  :label="statusBadge[doc.status]?.label ?? ''"
                  :color="(statusBadge[doc.status]?.color ?? 'neutral') as any"
                  variant="soft" size="sm" />
                <UDropdownMenu :items="getItemActions(doc)" :popper="{ placement: 'bottom-end' }">
                  <UButton
                    icon="i-tabler-dots-vertical"
                    color="neutral"
                    variant="ghost"
                    size="xs"
                    square
                    class="opacity-0 group-hover:opacity-100 transition-opacity" />
                </UDropdownMenu>
              </div>
            </div>
          </div>
        </div>
      </div>
    </AdminPageContent>

    <AdminPageFooter>
      <template #left>
        <template v-if="docs.length > 0">
          <UCheckbox
            :model-value="isAllDocsSelected"
            :indeterminate="isDocsIndeterminate"
            @update:model-value="toggleSelectAllDocs" />
          <template v-if="selectedDocs.length > 0">
            <span>{{ $t('common.selected_n', { n: selectedDocs.length }) }}</span>
            <USeparator orientation="vertical" class="h-4" />
            <USelect
              v-model="batchAction"
              :items="[
                { label: $t('admin.docs.status_published'), value: '2' },
                { label: $t('admin.docs.status_draft'), value: '1' },
                { label: $t('admin.docs.status_archived'), value: '3' },
                { label: $t('common.delete'), value: 'delete' },
              ]"
              :placeholder="$t('admin.posts.batch_action')"
              class="w-36"
              size="sm" />
            <UButton color="primary" variant="soft" size="sm" :disabled="!batchAction" @click="applyBatchDocs">{{ $t('common.apply') }}</UButton>
            <UButton color="neutral" variant="ghost" size="sm" @click="selectedDocs = []; batchAction = undefined">{{ $t('common.cancel') }}</UButton>
          </template>
          <span v-else class="text-xs">{{ $t('common.selectAll') }}</span>
        </template>
      </template>
      <template #right>
        <UPagination v-if="total > pageSize" v-model:page="currentPage" :total="total" :items-per-page="pageSize" size="sm" />
      </template>
    </AdminPageFooter>

    <!-- Delete Modal -->
    <UModal v-model:open="showDeleteModal">
      <template #content>
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
            </div>
            <div>
              <h3 class="font-semibold text-highlighted">{{ $t('admin.docs.delete_title') }}</h3>
              <p class="text-sm text-muted mt-0.5">{{ $t('admin.docs.delete_desc') }}</p>
            </div>
          </div>
          <div class="flex justify-end gap-2 mt-6">
            <UButton color="neutral" variant="outline" @click="showDeleteModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="error" :loading="deleteLoading" @click="confirmDelete">{{ $t('common.delete') }}</UButton>
          </div>
        </div>
      </template>
    </UModal>
  </AdminPageContainer>
</template>
