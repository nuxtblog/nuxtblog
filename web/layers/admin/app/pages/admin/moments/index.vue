<script setup lang="ts">
import type { MomentItem, MomentVisibility } from '~/types/api/moment'

definePageMeta({ layout: 'admin' })

const { t } = useI18n()
const toast = useToast()
const momentApi = useMomentApi()

// loading
const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)

// state
const moments = ref<MomentItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)
const filterVisibility = ref('all')
const searchKeyword = ref('')
const filterAuthorId = ref<number | undefined>(undefined)
const authors = ref<Array<{ id: number; username: string; nickname: string }>>([])

// delete modal
const showDeleteModal = ref(false)
const deleteTarget = ref<MomentItem | null>(null)
const deleteLoading = ref(false)

// edit modal
const showEditModal = ref(false)
const editTarget = ref<MomentItem | null>(null)
const editContent = ref('')
const editVisibility = ref<MomentVisibility>(1)
const editLoading = ref(false)

// batch selection
const selectedMoments = ref<number[]>([])
const batchAction = ref<string | undefined>(undefined)

const isAllMomentsSelected = computed(() =>
  moments.value.length > 0 && moments.value.every(m => selectedMoments.value.includes(m.id))
)
const isMomentsIndeterminate = computed(() =>
  selectedMoments.value.length > 0 && !isAllMomentsSelected.value
)

function toggleSelectMoment(id: number) {
  const idx = selectedMoments.value.indexOf(id)
  if (idx > -1) selectedMoments.value.splice(idx, 1)
  else selectedMoments.value.push(id)
}

function toggleSelectAllMoments() {
  if (isAllMomentsSelected.value) selectedMoments.value = []
  else selectedMoments.value = moments.value.map(m => m.id)
}

async function applyBatchMoments() {
  if (!batchAction.value || !selectedMoments.value.length) return
  try {
    if (batchAction.value === 'delete') {
      await Promise.all(selectedMoments.value.map(id => momentApi.deleteMoment(id)))
      toast.add({ title: t('common.deleted'), color: 'success' })
    } else {
      const vis = Number(batchAction.value) as MomentVisibility
      await Promise.all(selectedMoments.value.map(id => momentApi.updateMoment(id, { visibility: vis })))
      toast.add({ title: t('common.updated'), color: 'success' })
    }
    selectedMoments.value = []
    batchAction.value = undefined
    fetchMoments()
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.operation_failed'), color: 'error' })
  }
}

const visibilityTabs = computed(() => [
  { value: 'all', label: t('admin.moments.tab_all') },
  { value: '1', label: t('admin.moments.tab_public') },
  { value: '2', label: t('admin.moments.tab_private') },
  { value: '3', label: t('admin.moments.tab_followers') },
])

const visibilityBadge: Record<number, { label: string; color: 'success' | 'neutral' | 'warning' }> = {
  1: { label: t('admin.moments.visibility_public'), color: 'success' },
  2: { label: t('admin.moments.visibility_private'), color: 'neutral' },
  3: { label: t('admin.moments.visibility_followers'), color: 'warning' },
}

const visibilityOptions = computed(() => [
  { label: t('admin.moments.visibility_public'), value: 1 },
  { label: t('admin.moments.visibility_private'), value: 2 },
  { label: t('admin.moments.visibility_followers'), value: 3 },
])

const batchItems = computed(() => [
  { label: t('admin.moments.visibility_public'), value: '1' },
  { label: t('admin.moments.visibility_private'), value: '2' },
  { label: t('admin.moments.visibility_followers'), value: '3' },
  { label: t('common.delete'), value: 'delete' },
])

async function fetchMoments() {
  rawLoading.value = true
  try {
    const params: any = { page: currentPage.value, page_size: pageSize.value }
    if (filterVisibility.value !== 'all') params.visibility = Number(filterVisibility.value)
    if (searchKeyword.value.trim()) params.keyword = searchKeyword.value.trim()
    if (filterAuthorId.value) params.author_id = filterAuthorId.value
    const res = await momentApi.getMoments(params)
    moments.value = res.data ?? []
    total.value = res.total ?? 0
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.load_failed'), color: 'error' })
  } finally {
    rawLoading.value = false
  }
}

function onVisibilityTab(val: string) {
  filterVisibility.value = val
  currentPage.value = 1
}

function openDeleteModal(item: MomentItem) {
  deleteTarget.value = item
  showDeleteModal.value = true
}

function openEditModal(item: MomentItem) {
  editTarget.value = item
  editContent.value = item.content
  editVisibility.value = item.visibility
  showEditModal.value = true
}

async function confirmEdit() {
  if (!editTarget.value) return
  editLoading.value = true
  try {
    await momentApi.updateMoment(editTarget.value.id, {
      content: editContent.value,
      visibility: editVisibility.value,
    })
    toast.add({ title: t('common.saved'), color: 'success' })
    showEditModal.value = false
    fetchMoments()
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.save_failed'), color: 'error' })
  } finally {
    editLoading.value = false
  }
}

function openMomentFrontend(item: MomentItem) {
  if (import.meta.client) window.open(`/moments/${item.id}`, '_blank')
}

function getMomentActions(item: MomentItem) {
  return [
    [
      { label: t('common.edit'), icon: 'i-tabler-pencil', onClick: () => openEditModal(item) },
      { label: t('common.preview'), icon: 'i-tabler-eye', onClick: () => openMomentFrontend(item) },
    ],
    [
      { label: t('common.delete'), icon: 'i-tabler-trash', color: 'error' as const, onClick: () => openDeleteModal(item) },
    ],
  ]
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  deleteLoading.value = true
  try {
    await momentApi.deleteMoment(deleteTarget.value.id)
    toast.add({ title: t('common.deleted'), color: 'success' })
    showDeleteModal.value = false
    fetchMoments()
  } catch (e: any) {
    toast.add({ title: e?.message ?? t('common.delete_failed'), color: 'error' })
  } finally {
    deleteLoading.value = false
  }
}

function formatDate(d: string) {
  return d ? new Date(d).toLocaleString('zh-CN') : '-'
}

function getAuthorInitials(item: MomentItem) {
  const name = item.author?.nickname || item.author?.username || '?'
  return name.charAt(0).toUpperCase()
}

async function fetchAuthors() {
  try {
    const res = await useApiFetch().apiFetch<{ list: Array<{ id: number; username: string; display_name: string }> }>('/users', { params: { size: 100 } })
    authors.value = (res.list || []).map(u => ({ id: u.id, username: u.username, nickname: u.display_name || u.username }))
  } catch {}
}

watch([filterVisibility, filterAuthorId], () => {
  currentPage.value = 1
  selectedMoments.value = []
  batchAction.value = undefined
  fetchMoments()
})
watch(currentPage, fetchMoments)
watch(searchKeyword, useDebounceFn(() => {
  currentPage.value = 1
  selectedMoments.value = []
  batchAction.value = undefined
  fetchMoments()
}, 400))
onMounted(() => { fetchMoments(); fetchAuthors() })
</script>

<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.moments.title')" :subtitle="$t('admin.moments.subtitle')" />

    <AdminPageContent>
      <!-- Visibility Tabs -->
      <div class="flex items-center gap-1 border-b border-default pb-0 mb-4 overflow-x-auto">
        <button
          v-for="tab in visibilityTabs" :key="tab.value"
          class="px-3 py-2 text-sm font-medium rounded-t transition-colors whitespace-nowrap"
          :class="filterVisibility === tab.value ? 'text-primary border-b-2 border-primary' : 'text-muted hover:text-highlighted'"
          @click="onVisibilityTab(tab.value)">
          {{ tab.label }}
        </button>
      </div>

      <!-- Toolbar -->
      <div class="flex flex-wrap items-center gap-3 mb-4">
        <UInput
          v-model="searchKeyword"
          :placeholder="$t('admin.moments.search_placeholder')"
          leading-icon="i-tabler-search"
          class="w-56"
          size="sm">
          <template v-if="searchKeyword" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="searchKeyword = ''" />
          </template>
        </UInput>
        <AdminSearchableSelect
          v-model="filterAuthorId"
          :items="[{ label: $t('admin.moments.all_authors'), value: undefined }, ...authors.map(a => ({ label: a.nickname || a.username, value: a.id }))]"
          :placeholder="$t('admin.moments.all_authors')"
          :search-placeholder="$t('admin.posts.search_authors')" />
      </div>

      <!-- Skeleton -->
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 5" :key="i" class="flex items-start gap-4 p-4 border border-default rounded-lg">
          <USkeleton class="size-4 rounded shrink-0 mt-3" />
          <USkeleton class="size-10 rounded-full shrink-0" />
          <div class="flex-1 space-y-2">
            <USkeleton class="h-4 w-1/4" />
            <USkeleton class="h-4 w-full" />
            <USkeleton class="h-4 w-3/4" />
            <div class="flex gap-4">
              <USkeleton class="h-3 w-12" />
              <USkeleton class="h-3 w-12" />
              <USkeleton class="h-3 w-12" />
            </div>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="moments.length === 0" class="flex flex-col items-center justify-center py-16">
        <UIcon name="i-tabler-camera-off" class="size-16 text-muted mb-4" />
        <h3 class="text-lg font-medium text-highlighted mb-1">{{ $t('admin.moments.no_moments') }}</h3>
        <p class="text-sm text-muted">{{ $t('admin.moments.no_moments_desc') }}</p>
      </div>

      <!-- List -->
      <div v-else class="space-y-3">
        <div
          v-for="moment in moments" :key="moment.id"
          class="flex items-start gap-4 p-4 border border-default rounded-lg group hover:shadow-sm transition-all bg-default">
          <UCheckbox :model-value="selectedMoments.includes(moment.id)" @update:model-value="toggleSelectMoment(moment.id)" class="mt-2.5" />
          <!-- Author avatar -->
          <div class="size-10 rounded-full bg-primary/10 flex items-center justify-center shrink-0 text-sm font-semibold text-primary">
            {{ getAuthorInitials(moment) }}
          </div>

          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-4">
              <div class="flex-1 min-w-0">
                <!-- Author row -->
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-sm font-medium text-highlighted">
                    {{ moment.author?.nickname || moment.author?.username || '未知用户' }}
                  </span>
                  <span class="text-xs text-muted">{{ formatDate(moment.created_at) }}</span>
                </div>

                <!-- Content -->
                <p class="text-sm text-default leading-relaxed line-clamp-2 whitespace-pre-wrap">{{ moment.content }}</p>

                <!-- Media badges -->
                <div v-if="moment.media && moment.media.length > 0" class="mt-2">
                  <UBadge :label="`${moment.media.length} 张图片`" color="neutral" variant="soft" size="sm" icon="i-tabler-photo" />
                </div>

                <!-- Stats -->
                <div v-if="moment.stats" class="flex items-center gap-4 mt-2">
                  <span class="flex items-center gap-1 text-xs text-muted">
                    <UIcon name="i-tabler-eye" class="size-3.5" />
                    {{ moment.stats.view_count }}
                  </span>
                  <span class="flex items-center gap-1 text-xs text-muted">
                    <UIcon name="i-tabler-heart" class="size-3.5" />
                    {{ moment.stats.like_count }}
                  </span>
                  <span class="flex items-center gap-1 text-xs text-muted">
                    <UIcon name="i-tabler-message" class="size-3.5" />
                    {{ moment.stats.comment_count }}
                  </span>
                </div>
              </div>

              <!-- Visibility badge + 3-dots -->
              <div class="flex items-center gap-3 shrink-0">
                <UBadge
                  v-if="visibilityBadge[moment.visibility]"
                  :label="visibilityBadge[moment.visibility]?.label ?? ''"
                  :color="(visibilityBadge[moment.visibility]?.color ?? 'neutral') as any"
                  variant="soft" size="sm" />
                <UDropdownMenu :items="getMomentActions(moment)" :popper="{ placement: 'bottom-end' }">
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
        <template v-if="moments.length > 0">
          <UCheckbox
            :model-value="isAllMomentsSelected"
            :indeterminate="isMomentsIndeterminate"
            @update:model-value="toggleSelectAllMoments" />
          <template v-if="selectedMoments.length > 0">
            <span>{{ $t('common.selected_n', { n: selectedMoments.length }) }}</span>
            <USeparator orientation="vertical" class="h-4" />
            <USelect
              v-model="batchAction"
              :items="batchItems"
              :placeholder="$t('admin.posts.batch_action')"
              class="w-44"
              size="sm" />
            <UButton color="primary" variant="soft" size="sm" :disabled="!batchAction" @click="applyBatchMoments">{{ $t('common.apply') }}</UButton>
            <UButton color="neutral" variant="ghost" size="sm" @click="selectedMoments = []; batchAction = undefined">{{ $t('common.cancel') }}</UButton>
          </template>
          <span v-else class="text-xs">{{ $t('common.selectAll') }}</span>
        </template>
      </template>
      <template #right>
        <UPagination v-if="total > pageSize" v-model:page="currentPage" :total="total" :items-per-page="pageSize" size="sm" />
      </template>
    </AdminPageFooter>

    <!-- Edit Modal -->
    <UModal v-model:open="showEditModal" :ui="{ content: 'max-w-lg' }">
      <template #content>
        <div class="p-6">
          <h3 class="text-base font-semibold text-highlighted mb-4">{{ $t('admin.moments.edit_title') }}</h3>
          <div class="space-y-4">
            <UFormField :label="$t('admin.moments.field_content')">
              <UTextarea v-model="editContent" :rows="5" class="w-full" />
            </UFormField>
            <UFormField :label="$t('admin.moments.field_visibility')">
              <USelect v-model="editVisibility" :items="visibilityOptions" class="w-full" />
            </UFormField>
          </div>
          <div class="flex justify-end gap-2 mt-6">
            <UButton color="neutral" variant="outline" @click="showEditModal = false">{{ $t('common.cancel') }}</UButton>
            <UButton color="primary" :loading="editLoading" @click="confirmEdit">{{ $t('common.save') }}</UButton>
          </div>
        </div>
      </template>
    </UModal>

    <!-- Delete Modal -->
    <UModal v-model:open="showDeleteModal">
      <template #content>
        <div class="p-6">
          <div class="flex items-center gap-3 mb-4">
            <div class="size-10 rounded-full bg-error/10 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-alert-triangle" class="size-5 text-error" />
            </div>
            <div>
              <h3 class="font-semibold text-highlighted">{{ $t('admin.moments.delete_title') }}</h3>
              <p class="text-sm text-muted mt-0.5">{{ $t('admin.moments.delete_desc') }}</p>
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
