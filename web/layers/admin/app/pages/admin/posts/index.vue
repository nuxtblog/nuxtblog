<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.posts.title')" :subtitle="$t('admin.posts.subtitle')">
      <template #actions>
        <UButton
          v-if="filterStatus !== 'trash'"
          as="NuxtLink"
          to="/admin/posts/new"
          icon="i-tabler-plus"
          color="primary">
          {{ $t('admin.posts.new_post') }}
        </UButton>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- Status tabs -->
      <div class="flex items-center gap-1 border-b border-default pb-0 mb-4 overflow-x-auto">
        <button
          v-for="s in statusTabs"
          :key="s.value"
          class="px-3 py-2 text-sm font-medium rounded-t transition-colors whitespace-nowrap"
          :class="filterStatus === s.value
            ? 'text-primary border-b-2 border-primary'
            : 'text-muted hover:text-highlighted'"
          @click="onStatusTab(s.value)">
          {{ s.label }}
          <span v-if="s.count != null" class="ml-1 text-xs text-muted">({{ s.count }})</span>
        </button>
      </div>

      <!-- Filters toolbar -->
      <PostFiltersToolbar
        v-model:search-keyword="searchKeyword"
        v-model:filter-category="filterCategory"
        v-model:filter-author="filterAuthor"
        v-model:sort-by="sortBy"
        v-model:page-size="pageSize"
        v-model:view-mode="viewMode"
        :categories="categories"
        :authors="authors"
        :filter-status="filterStatus"
        :total-posts="totalPosts"
        @clear-trash="showClearTrashModal = true"
      />

      <!-- Post list content -->
      <div class="flex-1">
        <!-- Loading -->
        <div v-if="displayLoading" class="space-y-3">
          <div
            v-for="i in 8"
            :key="i"
            class="flex items-center gap-4 p-4 border border-default rounded-md">
            <USkeleton class="size-4 rounded shrink-0" />
            <USkeleton class="h-16 w-24 rounded shrink-0" />
            <div class="flex-1 space-y-2 min-w-0">
              <USkeleton class="h-4 w-3/4" />
              <USkeleton class="h-3 w-1/2" />
              <div class="flex gap-4">
                <USkeleton class="h-3 w-16" />
                <USkeleton class="h-3 w-12" />
                <USkeleton class="h-3 w-10" />
              </div>
            </div>
            <USkeleton class="h-5 w-14 rounded-full shrink-0" />
          </div>
        </div>

        <!-- List view -->
        <div v-else-if="viewMode === 'list' && posts.length > 0" class="space-y-3">
          <PostListRow
            v-for="post in posts"
            :key="post.id"
            :post="post"
            :is-selected="selectedPosts.includes(post.id)"
            :actions="getPostActions(post)"
            :filter-status="filterStatus"
            @toggle-select="toggleSelect(post.id)" />
        </div>

        <!-- Grid view -->
        <div v-else-if="viewMode === 'grid' && posts.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          <PostGridCard
            v-for="post in posts"
            :key="post.id"
            :post="post"
            :is-selected="selectedPosts.includes(post.id)"
            :actions="getPostActions(post)"
            :filter-status="filterStatus"
            @toggle-select="toggleSelect(post.id)" />
        </div>

        <!-- Empty state -->
        <div
          v-else-if="!displayLoading && posts.length === 0"
          class="flex flex-col items-center justify-center py-16">
          <UIcon
            :name="hasActiveFilters ? 'i-tabler-search-off' : filterStatus === 'trash' ? 'i-tabler-trash-off' : 'i-tabler-file-off'"
            class="size-16 text-muted mb-4" />
          <h3 class="text-lg font-medium text-highlighted mb-1">
            {{ hasActiveFilters ? $t('admin.posts.no_results') : filterStatus === 'trash' ? $t('admin.posts.trash_empty') : $t('admin.posts.no_posts') }}
          </h3>
          <p class="text-sm text-muted mb-4">
            {{ hasActiveFilters ? $t('admin.posts.try_modify_filters') : filterStatus === 'trash' ? $t('admin.posts.no_deleted') : $t('admin.posts.start_writing') }}
          </p>
          <UButton
            v-if="hasActiveFilters"
            color="neutral" variant="outline" icon="i-tabler-filter-off"
            @click="clearFilters">
            {{ $t('admin.posts.clear_filters') }}
          </UButton>
          <UButton
            v-else-if="filterStatus !== 'trash'"
            as="NuxtLink" to="/admin/posts/new" icon="i-tabler-plus" color="primary">
            {{ $t('admin.posts.new_post') }}
          </UButton>
        </div>
      </div>
    </AdminPageContent>

    <AdminPageFooter>
      <template #left>
        <template v-if="posts.length > 0">
          <UCheckbox
            :model-value="isAllSelected"
            :indeterminate="isIndeterminate"
            @update:model-value="toggleSelectAll" />
          <template v-if="selectedPosts.length > 0">
            <span>{{ $t('common.selected_n', { n: selectedPosts.length }) }}</span>
            <USeparator orientation="vertical" class="h-4" />
            <UButton
              v-if="filterStatus !== 'trash'"
              color="neutral" variant="outline" size="sm" icon="i-tabler-pencil"
              @click="showBatchEditModal = true">
              {{ $t('admin.posts.batch_edit_btn') }}
            </UButton>
            <USelect
              v-model="batchAction"
              :items="batchItems"
              :placeholder="$t('admin.posts.batch_action')"
              class="w-36" size="sm" />
            <UButton color="primary" variant="soft" size="sm" :disabled="!batchAction" @click="applyBatch">
              {{ $t('common.apply') }}
            </UButton>
            <UButton color="neutral" variant="ghost" size="sm" @click="selectedPosts = []">
              {{ $t('common.cancel') }}
            </UButton>
          </template>
          <span v-else class="text-xs">{{ $t('common.selectAll') }}</span>
        </template>
        <span v-else-if="!displayLoading && totalPosts > 0">
          {{ $t('common.total', { count: totalPosts }) }} · {{ currentPage }}/{{ totalPages }}
        </span>
      </template>
      <template #right>
        <UPagination
          v-if="!displayLoading && totalPosts > pageSize"
          v-model:page="currentPage"
          :total="totalPosts"
          :items-per-page="pageSize"
          size="sm" />
      </template>
    </AdminPageFooter>

    <!-- Delete confirm modal -->
    <PostDeleteModal
      v-model:open="showDeleteModal"
      mode="single"
      @confirm="confirmDelete" />

    <!-- Clear trash confirm modal -->
    <PostDeleteModal
      v-model:open="showClearTrashModal"
      mode="clearTrash"
      :trash-count="totalPosts"
      @confirm="confirmClearTrash" />

    <!-- Batch edit modal -->
    <PostBatchEditModal
      v-model:open="showBatchEditModal"
      :post-ids="selectedPosts"
      :categories="categories"
      :tags="tags"
      :authors="authors"
      :is-admin="isAdmin"
      @applied="onBatchEditApplied" />
  </AdminPageContainer>
</template>

<script setup lang="ts">
import { dispatchCommand } from '~/composables/useNuxtblogAdmin'
import { CONTRIBUTION_SLOTS } from '~/config/contribution-slots'
import { usePluginContributionsStore } from '~/stores/plugin-contributions'
import { usePluginContextStore } from '~/stores/plugin-context'
import { usePluginLocale } from '~/composables/usePluginLocale'

interface PostListItem {
  id: number
  post_type: string
  title: string
  slug: string
  excerpt: string
  featured_img?: { id: number; url: string }
  status: string
  published_at?: string
  comment_count: number
  view_count: number
  author?: { id: number; username: string; nickname: string; avatar?: string }
  metas?: Record<string, string>
  created_at: string
  updated_at: string
}

interface PostListData {
  data: PostListItem[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

interface PostStats {
  total_posts: number
  published_posts: number
  draft_posts: number
  private_posts: number
  archived_posts: number
}

interface TaxonomyItem {
  id: number
  term: { id: number; name: string; slug: string }
  taxonomy: string
  post_count: number
  parent_id?: number
}

interface UserItem {
  id: number
  username: string
  display_name: string
}

const { apiFetch } = useApiFetch()
const toast = useToast()
const route = useRoute()
const router = useRouter()
const { t } = useI18n()

// ── State ──────────────────────────────────────────────────────────────────
const posts = ref<PostListItem[]>([])
const loading = ref(true)
const displayLoading = useMinLoading(loading)
const filterStatus = ref((route.query.status as string) || 'all')
const filterCategory = ref<string | undefined>((route.query.cat as string) || undefined)
const filterAuthor = ref<string | undefined>((route.query.author as string) || undefined)
const sortBy = ref((route.query.sort as string) || 'created_desc')
const viewMode = ref<'list' | 'grid'>('list')
const selectedPosts = ref<number[]>([])
const batchAction = ref<string | undefined>(undefined)
const searchKeyword = ref((route.query.q as string) || '')
const currentPage = ref(Number(route.query.page) || 1)
const pageSize = ref(Number(route.query.size) || 20)
const totalPosts = ref(0)
const totalPages = ref(0)

// ── Filter options ──────────────────────────────────────────────────────────
const categories = ref<TaxonomyItem[]>([])
const tags = ref<TaxonomyItem[]>([])
const authors = ref<UserItem[]>([])

// ── Status tabs ─────────────────────────────────────────────────────────────
const tabCounts = ref({ all: null as number | null, published: null as number | null, draft: null as number | null, private: null as number | null, trash: null as number | null })

const statusTabs = computed(() => [
  { label: t('admin.posts.status_all'), value: 'all', count: tabCounts.value.all },
  { label: t('admin.posts.status_published'), value: 'published', count: tabCounts.value.published },
  { label: t('admin.posts.status_draft'), value: 'draft', count: tabCounts.value.draft },
  { label: t('admin.posts.status_private'), value: 'private', count: tabCounts.value.private },
  { label: t('admin.posts.status_trash'), value: 'trash', count: tabCounts.value.trash },
])

const batchItems = computed(() =>
  filterStatus.value === 'trash'
    ? [
        { label: t('admin.posts.restore_action'), value: 'restore' },
        { label: t('admin.posts.hard_delete'), value: 'delete' },
      ]
    : [
        { label: t('common.publish'), value: 'publish' },
        { label: t('admin.posts.status_draft'), value: 'draft' },
        { label: t('admin.posts.trash_action'), value: 'trash' },
      ],
)

// ── Active filters ──────────────────────────────────────────────────────────
const hasActiveFilters = computed(() =>
  searchKeyword.value.trim() !== '' || !!filterCategory.value || !!filterAuthor.value,
)
const clearFilters = () => {
  searchKeyword.value = ''
  filterCategory.value = undefined
  filterAuthor.value = undefined
}

// ── Fetch ──────────────────────────────────────────────────────────────────
const fetchPosts = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: pageSize.value,
      post_type: 'post',
    }
    if (filterStatus.value !== 'all') params.status = filterStatus.value
    if (filterCategory.value) params.include_category_ids = filterCategory.value
    if (filterAuthor.value) params.author_id = Number(filterAuthor.value)
    if (searchKeyword.value.trim()) params.keyword = searchKeyword.value.trim()
    const sortMap: Record<string, string> = {
      created_desc: 'created_at',
      created_asc: 'created_at',
      updated_desc: 'updated_at',
      views_desc: 'view_count',
    }
    if (sortMap[sortBy.value]) params.sort_by = sortMap[sortBy.value]
    const data = await apiFetch<PostListData>('/posts', { params })
    posts.value = data.data || []
    totalPosts.value = data.total || 0
    totalPages.value = data.total_pages || 0
  } catch {
    posts.value = []
    totalPosts.value = 0
  } finally {
    loading.value = false
  }
}

const fetchStats = async () => {
  try {
    const [statsRes, trashRes] = await Promise.allSettled([
      apiFetch<PostStats>('/posts/stats'),
      apiFetch<PostListData>('/posts', { params: { status: 'trash', post_type: 'post', page_size: 1 } }),
    ])
    if (statsRes.status === 'fulfilled') {
      const s = statsRes.value
      tabCounts.value.all = s.total_posts
      tabCounts.value.published = s.published_posts
      tabCounts.value.draft = s.draft_posts
      tabCounts.value.private = s.private_posts
    }
    if (trashRes.status === 'fulfilled') {
      tabCounts.value.trash = trashRes.value.total
    }
  } catch {}
}

const fetchFilters = async () => {
  const [catRes, tagRes, userRes] = await Promise.allSettled([
    apiFetch<{ list: TaxonomyItem[] }>('/taxonomies', { params: { taxonomy: 'category', page_size: 100 } }),
    apiFetch<{ list: TaxonomyItem[] }>('/taxonomies', { params: { taxonomy: 'tag', page_size: 200 } }),
    apiFetch<{ list: UserItem[] }>('/users', { params: { size: 100 } }),
  ])
  if (catRes.status === 'fulfilled') categories.value = catRes.value.list || []
  if (tagRes.status === 'fulfilled') tags.value = tagRes.value.list || []
  if (userRes.status === 'fulfilled') authors.value = userRes.value.list || []
}

const debouncedKeyword = refDebounced(searchKeyword, 350)

watch([filterStatus, filterCategory, filterAuthor, sortBy, debouncedKeyword, pageSize], () => {
  currentPage.value = 1
  selectedPosts.value = []
  batchAction.value = undefined
  fetchPosts()
  fetchStats()
})

watch(currentPage, fetchPosts)

watch([filterStatus, filterCategory, filterAuthor, sortBy, debouncedKeyword, currentPage, pageSize], () => {
  const query: Record<string, string> = {}
  if (filterStatus.value !== 'all') query.status = filterStatus.value
  if (debouncedKeyword.value) query.q = debouncedKeyword.value
  if (currentPage.value > 1) query.page = String(currentPage.value)
  if (filterCategory.value) query.cat = filterCategory.value
  if (filterAuthor.value) query.author = filterAuthor.value
  if (sortBy.value !== 'created_desc') query.sort = sortBy.value
  if (pageSize.value !== 20) query.size = String(pageSize.value)
  router.replace({ query })
})

onMounted(() => {
  fetchPosts()
  fetchFilters()
  fetchStats()
})

// ── Tab switch ─────────────────────────────────────────────────────────────
const onStatusTab = (val: string) => {
  filterStatus.value = val
  currentPage.value = 1
  selectedPosts.value = []
  batchAction.value = undefined
}

// ── Selection ──────────────────────────────────────────────────────────────
const isAllSelected = computed(() => posts.value.length > 0 && posts.value.every(p => selectedPosts.value.includes(p.id)))
const isIndeterminate = computed(() => selectedPosts.value.length > 0 && !isAllSelected.value)
const toggleSelectAll = () => {
  selectedPosts.value = isAllSelected.value ? [] : posts.value.map(p => p.id)
}
const toggleSelect = (id: number) => {
  const idx = selectedPosts.value.indexOf(id)
  if (idx > -1) selectedPosts.value.splice(idx, 1)
  else selectedPosts.value.push(id)
}

// ── Batch ──────────────────────────────────────────────────────────────────
const applyBatch = async () => {
  if (!batchAction.value || selectedPosts.value.length === 0) return
  try {
    const res = await apiFetch<{ affected: number }>('/posts/batch', {
      method: 'POST',
      body: { ids: selectedPosts.value, action: batchAction.value },
    })
    toast.add({ title: t('admin.posts.op_success', { n: res.affected }), icon: 'i-tabler-circle-check', color: 'success' })
    selectedPosts.value = []
    batchAction.value = undefined
    fetchPosts()
    fetchStats()
  } catch {
    toast.add({ title: t('admin.posts.batch_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  }
}

// ── Batch edit modal ────────────────────────────────────────────────────────
const showBatchEditModal = ref(false)
const authStore = useAuthStore()
const isAdmin = computed(() => (authStore.user?.role ?? 0) >= 3)

const onBatchEditApplied = () => {
  selectedPosts.value = []
  fetchPosts()
  fetchStats()
}

// ── Clear trash ────────────────────────────────────────────────────────────
const showClearTrashModal = ref(false)

const confirmClearTrash = async () => {
  try {
    const data = await apiFetch<PostListData>('/posts', {
      params: { status: 'trash', post_type: 'post', page_size: 1000 },
    })
    const ids = (data.data || []).map(p => p.id)
    if (ids.length > 0) {
      await apiFetch('/posts/batch', { method: 'POST', body: { ids, action: 'delete' } })
    }
    toast.add({ title: t('admin.posts.cleared_trash', { n: ids.length }), icon: 'i-tabler-trash-x', color: 'success' })
    showClearTrashModal.value = false
    fetchPosts()
    fetchStats()
  } catch {
    toast.add({ title: t('admin.posts.clear_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  }
}

// ── Post operations ─────────────────────────────────────────────────────────
const previewPost = (post: PostListItem) => window.open(`/posts/${post.slug}`, '_blank')

const copyPost = async (post: PostListItem) => {
  try {
    const detail = await apiFetch<{ id: number; title: string; content: string; excerpt: string }>(`/posts/${post.id}`)
    await apiFetch('/posts', {
      method: 'POST',
      body: { title: `${detail.title} (副本)`, slug: `${post.slug}-copy-${Date.now()}`, content: detail.content, excerpt: detail.excerpt, status: 1, post_type: 1 },
    })
    toast.add({ title: t('admin.posts.copy_created'), icon: 'i-tabler-copy', color: 'success' })
    fetchPosts()
    fetchStats()
  } catch {
    toast.add({ title: t('admin.posts.copy_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  }
}

const trashPost = async (post: PostListItem) => {
  try {
    await apiFetch(`/posts/${post.id}/trash`, { method: 'POST' })
    toast.add({ title: t('admin.posts.moved_to_trash'), icon: 'i-tabler-trash', color: 'neutral' })
    fetchPosts()
    fetchStats()
  } catch {
    toast.add({ title: t('common.operation_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  }
}

const restorePost = async (post: PostListItem) => {
  try {
    await apiFetch(`/posts/${post.id}/restore`, { method: 'POST' })
    toast.add({ title: t('admin.posts.restored'), icon: 'i-tabler-restore', color: 'success' })
    fetchPosts()
    fetchStats()
  } catch {
    toast.add({ title: t('admin.posts.restore_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  }
}

const showDeleteModal = ref(false)
const pendingDeleteId = ref<number | null>(null)

const confirmDeletePost = (post: PostListItem) => {
  pendingDeleteId.value = post.id
  showDeleteModal.value = true
}

const confirmDelete = async () => {
  if (!pendingDeleteId.value) return
  try {
    await apiFetch(`/posts/${pendingDeleteId.value}`, { method: 'DELETE' })
    toast.add({ title: t('admin.posts.permanently_deleted'), icon: 'i-tabler-trash-x', color: 'error' })
    showDeleteModal.value = false
    fetchPosts()
    fetchStats()
  } catch {
    toast.add({ title: t('common.delete_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  }
}

// ── Plugin contribution commands ──────────────────────────────────────────
const contributionsStore = usePluginContributionsStore()
const contextStore = usePluginContextStore()
const { t: pluginT } = usePluginLocale()
const pluginMenuItems = contributionsStore.getMenuItems(CONTRIBUTION_SLOTS.POST_LIST_ROW_ACTION)

const getPluginActions = (post: PostListItem) => {
  const items = pluginMenuItems.value.filter(item => !item.when || contextStore.evaluateWhen(item.when))
  if (items.length === 0) return []
  return [items.map(item => ({
    label: pluginT(item) || item.command,
    icon: item.icon,
    onClick: () => dispatchCommand(item.command, { source: 'post-list', postId: post.id, postTitle: post.title }),
  }))]
}

// ── Actions menu ───────────────────────────────────────────────────────────
const getPostActions = (post: PostListItem) => {
  if (filterStatus.value === 'trash') {
    return [
      [{ label: t('admin.posts.restore_action'), icon: 'i-tabler-restore', onClick: () => restorePost(post) }],
      [{ label: t('admin.posts.hard_delete'), icon: 'i-tabler-trash-x', color: 'error' as const, onClick: () => confirmDeletePost(post) }],
    ]
  }
  return [
    [
      { label: t('common.edit'), icon: 'i-tabler-pencil', to: `/admin/posts/edit/${post.id}` },
      { label: t('common.preview'), icon: 'i-tabler-eye', onClick: () => previewPost(post) },
      { label: t('common.copy'), icon: 'i-tabler-copy', onClick: () => copyPost(post) },
    ],
    ...getPluginActions(post),
    [{ label: t('admin.posts.trash_action'), icon: 'i-tabler-trash', color: 'error' as const, onClick: () => trashPost(post) }],
  ]
}
</script>
