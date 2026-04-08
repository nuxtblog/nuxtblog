<template>
  <AdminPageContainer>
    <AdminPageHeader :title="$t('admin.users.title')" :subtitle="$t('admin.users.subtitle')">
      <template #actions>
        <div class="flex items-center gap-3">
          <div class="flex items-center gap-3">
            <div class="text-center px-4 py-2 bg-default rounded-md border border-default">
              <div class="text-lg font-semibold text-highlighted">{{ stats.total }}</div>
              <div class="text-xs text-muted">{{ $t('admin.users.total') }}</div>
            </div>
            <div class="text-center px-4 py-2 bg-success/10 rounded-md border border-success/20">
              <div class="text-lg font-semibold text-success">{{ stats.active }}</div>
              <div class="text-xs text-success">{{ $t('admin.users.active') }}</div>
            </div>
          </div>
          <UButton color="primary" icon="i-tabler-plus" @click="openCreateModal">
            {{ $t('admin.users.add_user') }}
          </UButton>
        </div>
      </template>
    </AdminPageHeader>

    <AdminPageContent>
      <!-- 筛选工具栏 -->
      <div class="flex flex-wrap items-center gap-3 mb-4">
        <div class="flex items-center gap-1">
          <UButton
            v-for="s in statusFilters"
            :key="String(s.value)"
            :color="filterStatus === s.value ? 'primary' : 'neutral'"
            :variant="filterStatus === s.value ? 'soft' : 'ghost'"
            size="sm"
            @click="filterStatus = s.value">
            {{ s.label }}
            <span v-if="s.count !== undefined" class="ml-1 text-xs opacity-60">({{ s.count }})</span>
          </UButton>
        </div>

        <AdminSearchableSelect
          v-model="filterRole"
          :items="[{ label: $t('admin.users.all_roles'), value: 'all' }, ...roleList.map(r => ({ label: r.name, value: r.id }))]"
          :placeholder="$t('admin.users.all_roles')"
          :search-placeholder="$t('common.search')" />

        <UInput
          v-model="searchQuery"
          :placeholder="$t('admin.users.search_placeholder')"
          leading-icon="i-tabler-search"
          class="w-56"
          size="sm">
          <template v-if="searchQuery" #trailing>
            <UButton icon="i-tabler-x" color="neutral" variant="ghost" size="xs" @click="searchQuery = ''" />
          </template>
        </UInput>

        <USelect
          v-model="sortBy"
          :items="[
            { label: $t('admin.users.sort_newest'), value: 'created_desc' },
            { label: $t('admin.users.sort_oldest'), value: 'created_asc' },
            { label: $t('admin.users.sort_username'), value: 'username_asc' },
          ]"
          class="w-36"
          size="sm" />
      </div>

      <!-- 用户列表 -->
      <div>
        <!-- 加载状态 -->
        <div v-if="loading" class="space-y-3">
          <div
            v-for="i in 8" :key="i"
            class="flex items-center gap-4 p-4 bg-default border border-default rounded-md">
            <USkeleton class="size-4 rounded shrink-0" />
            <USkeleton class="size-10 rounded-full shrink-0" />
            <div class="flex-1 space-y-1.5 min-w-0">
              <USkeleton class="h-4 w-32" />
              <USkeleton class="h-3 w-48" />
            </div>
            <USkeleton class="h-5 w-16 rounded-full hidden sm:block" />
            <USkeleton class="h-5 w-14 rounded-full hidden md:block" />
            <USkeleton class="h-3 w-20 hidden lg:block" />
            <div class="flex gap-1 shrink-0">
              <USkeleton class="size-7 rounded" />
              <USkeleton class="size-7 rounded" />
            </div>
          </div>
        </div>

        <!-- 用户列表 -->
        <div v-else-if="filteredUsers.length > 0" class="space-y-3">
          <div
            v-for="user in filteredUsers"
            :key="user.id"
            class="flex items-center group gap-4 p-4 bg-default border border-default rounded-md hover:shadow-sm transition-all">
            <UCheckbox
              :model-value="selectedUsers.includes(user.id)"
              @update:model-value="toggleSelectUser(user.id)" />
            <BaseAvatar :src="user.avatar" :alt="user.display_name || user.username" size="md" class="shrink-0" />
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 flex-wrap">
                <span class="text-sm font-medium text-highlighted">{{ user.display_name || user.username }}</span>
                <span class="text-xs text-muted">@{{ user.username }}</span>
              </div>
              <div class="text-xs text-muted mt-0.5">{{ user.email }}</div>
              <div class="text-xs text-muted mt-0.5">{{ formatDate(user.created_at) }}</div>
            </div>
            <div class="flex items-center gap-2 shrink-0">
              <UBadge
                :label="roleMap[user.role]?.name || String(user.role)"
                :color="(roleMap[user.role]?.color as any) || 'neutral'"
                variant="soft" />
              <UBadge
                :label="statusConfig[user.status]?.label || String(user.status)"
                :color="(statusConfig[user.status]?.color as any) || 'neutral'"
                variant="soft" />
              <UDropdownMenu :items="getUserActions(user)" :popper="{ placement: 'bottom-end' }">
                <UButton
                  color="neutral" variant="ghost" icon="i-tabler-dots-vertical" square size="xs"
                  class="opacity-0 group-hover:opacity-100 transition-opacity" />
              </UDropdownMenu>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="flex flex-col items-center justify-center py-16">
          <UIcon name="i-tabler-users" class="text-muted mb-4 size-16" />
          <h3 class="text-lg font-medium text-highlighted mb-1">{{ $t('admin.users.no_users') }}</h3>
          <p class="text-sm text-muted">
            {{ searchQuery ? $t('admin.users.no_results') : $t('admin.users.no_registered') }}
          </p>
        </div>
      </div>
    </AdminPageContent>

    <AdminPageFooter>
      <template #left>
        <template v-if="filteredUsers.length > 0">
          <UCheckbox :model-value="isAllSelected" :indeterminate="isIndeterminate" @update:model-value="toggleSelectAll" />
          <template v-if="selectedUsers.length > 0">
            <span>{{ $t('admin.users.selected_n', { n: selectedUsers.length }) }}</span>
            <USeparator orientation="vertical" class="h-4" />
            <USelect
              v-model="batchAction"
              :items="[
                { label: $t('admin.users.batch_set_active'), value: 'active' },
                { label: $t('admin.users.batch_set_pending'), value: 'pending' },
                { label: $t('admin.users.batch_ban'), value: 'banned' },
              ]"
              :placeholder="$t('admin.users.batch_ops')"
              class="w-36"
              size="sm" />
            <UButton color="primary" variant="soft" size="sm" :disabled="!batchAction || batchAction === 'none'" @click="handleBatchAction">
              {{ $t('common.apply') }}
            </UButton>
            <UButton color="neutral" variant="ghost" size="sm" @click="selectedUsers = []; batchAction = 'none'">
              {{ $t('common.cancel') }}
            </UButton>
          </template>
          <span v-else class="text-xs">{{ $t('common.selectAll') }}</span>
        </template>
        <span v-else-if="!loading && totalUsers > 0">
          {{ $t('common.show_range', { from: (currentPage - 1) * pageSize + 1, to: Math.min(currentPage * pageSize, totalUsers), total: totalUsers }) }}
        </span>
      </template>
      <template #right>
        <UPagination
          v-if="!loading && totalUsers > pageSize"
          v-model:page="currentPage"
          :total="totalUsers"
          :items-per-page="pageSize"
          size="sm" />
      </template>
    </AdminPageFooter>

    <!-- Modals -->
    <UserFormModal
      v-model:open="showUserModal"
      :user="editingUser"
      :role-list="roleList"
      @saved="fetchUsers" />

    <UserResetPasswordModal
      v-model:user="resetPasswordUser" />

    <UserDeleteModal
      v-model:user="deleteConfirmUser"
      :current-user-id="authStore.user?.id"
      @deleted="fetchUsers" />

    <UserDetailModal
      v-model:user="viewingUser"
      :role-map="roleMap"
      :status-config="statusConfig"
      @edit="openEditModal"
      @reset-password="u => resetPasswordUser = u" />
  </AdminPageContainer>
</template>

<script setup lang="ts">
import type { UserListResponse, UserStatus, UserRole } from '~/types/api/user'

const toast = useToast()
const { t } = useI18n()
const userApi = useUserApi()
const authStore = useAuthStore()

// ── State ──────────────────────────────────────────────────────────────────
const users = ref<UserListResponse[]>([])
const stats = ref({ total: 0, active: 0, pending: 0, banned: 0 })
const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)

const filterStatus = ref<number | 'all'>('all')
const filterRole = ref<number | 'all'>('all')
const searchQuery = ref('')
const sortBy = ref('created_desc')
const selectedUsers = ref<number[]>([])
const batchAction = ref('none')
const currentPage = ref(1)
const pageSize = ref(20)
const totalUsers = ref(0)

// Modal state
const showUserModal = ref(false)
const editingUser = ref<UserListResponse | null>(null)
const resetPasswordUser = ref<UserListResponse | null>(null)
const deleteConfirmUser = ref<UserListResponse | null>(null)
const viewingUser = ref<UserListResponse | null>(null)

// ── Config ─────────────────────────────────────────────────────────────────
const getRoleMap = (): Record<number, { name: string; slug: string; color: string }> => ({
  1: { name: t('admin.users_new.role_subscriber'), slug: 'subscriber', color: 'neutral' },
  2: { name: t('admin.users_new.role_editor'), slug: 'editor', color: 'primary' },
  3: { name: t('admin.roles.super_admin'), slug: 'super_admin', color: 'warning' },
})
const roleMap = computed(() => getRoleMap())

const getStatusConfig = (): Record<number, { label: string; color: string }> => ({
  1: { label: t('admin.users.status_active'), color: 'success' },
  2: { label: t('admin.users.status_banned'), color: 'error' },
  3: { label: t('admin.users.status_pending'), color: 'warning' },
})
const statusConfig = computed(() => getStatusConfig())

const batchStatusMap: Record<string, UserStatus> = { active: 1, banned: 2, pending: 3 }

const roleList = computed(() =>
  Object.entries(getRoleMap()).map(([id, r]) => ({ id: Number(id) as UserRole, ...r })),
)

const statusFilters = ref([
  { label: '', value: 'all' as const, count: 0 },
  { label: '', value: 1, count: 0 },
  { label: '', value: 3, count: 0 },
  { label: '', value: 2, count: 0 },
])

// ── Computed ───────────────────────────────────────────────────────────────
const filteredUsers = computed(() => users.value)

const isAllSelected = computed(
  () => filteredUsers.value.length > 0 && filteredUsers.value.every(u => selectedUsers.value.includes(u.id)),
)
const isIndeterminate = computed(
  () => selectedUsers.value.length > 0 && !isAllSelected.value,
)

// ── Fetch ──────────────────────────────────────────────────────────────────
const fetchUsers = async () => {
  try {
    const data = await userApi.getUsers({
      page: currentPage.value,
      size: pageSize.value,
      keyword: searchQuery.value || undefined,
      status: filterStatus.value !== 'all' ? Number(filterStatus.value) as UserStatus : undefined,
      role: filterRole.value !== 'all' ? Number(filterRole.value) as UserRole : undefined,
      order_by: sortBy.value === 'username_asc' ? 'username' : 'created_at',
      order: sortBy.value === 'created_asc' ? 'asc' : 'desc',
    })
    users.value = data.list || []
    totalUsers.value = data.total || 0

    const activeCount = users.value.filter(u => u.status === 1).length
    const bannedCount = users.value.filter(u => u.status === 2).length
    const pendingCount = users.value.filter(u => u.status === 3).length
    stats.value = { total: data.total || 0, active: activeCount, pending: pendingCount, banned: bannedCount }

    statusFilters.value = [
      { label: t('admin.users.status_all'), value: 'all', count: stats.value.total },
      { label: t('admin.users.status_active'), value: 1, count: stats.value.active },
      { label: t('admin.users.status_pending'), value: 3, count: stats.value.pending },
      { label: t('admin.users.status_banned'), value: 2, count: stats.value.banned },
    ]
  } catch (error: any) {
    toast.add({ title: t('admin.users.fetch_failed'), description: error?.message, color: 'error' })
  } finally {
    rawLoading.value = false
  }
}

// ── Selection ──────────────────────────────────────────────────────────────
const toggleSelectUser = (userId: number) => {
  const idx = selectedUsers.value.indexOf(userId)
  if (idx > -1) selectedUsers.value.splice(idx, 1)
  else selectedUsers.value.push(userId)
}

const toggleSelectAll = (val?: boolean) => {
  if (val === false || isAllSelected.value) selectedUsers.value = []
  else selectedUsers.value = filteredUsers.value.map(u => u.id)
}

// ── Batch ──────────────────────────────────────────────────────────────────
const handleBatchAction = async () => {
  if (batchAction.value === 'none' || !selectedUsers.value.length) return
  const statusNum = batchStatusMap[batchAction.value]
  if (!statusNum) return
  try {
    await Promise.all(selectedUsers.value.map(id => userApi.updateUser(id, { status: statusNum })))
    await fetchUsers()
    selectedUsers.value = []
    batchAction.value = 'none'
  } catch (error: any) {
    toast.add({ title: t('admin.users.batch_failed'), description: error?.message, color: 'error' })
  }
}

// ── Modal openers ──────────────────────────────────────────────────────────
const openCreateModal = () => {
  editingUser.value = null
  showUserModal.value = true
}

const openEditModal = (user: UserListResponse) => {
  editingUser.value = user
  showUserModal.value = true
}

// ── Actions menu ───────────────────────────────────────────────────────────
const getUserActions = (user: UserListResponse) => [
  [
    { label: t('admin.users.view_detail'), icon: 'i-tabler-eye', onSelect: () => viewingUser.value = user },
    { label: t('common.edit'), icon: 'i-tabler-pencil', onSelect: () => openEditModal(user) },
    { label: t('admin.users.reset_password'), icon: 'i-tabler-lock', onSelect: () => resetPasswordUser.value = user },
  ],
  [
    {
      label: t('common.delete'),
      icon: 'i-tabler-trash',
      color: 'error' as const,
      disabled: user.id === authStore.user?.id,
      onSelect: () => deleteConfirmUser.value = user,
    },
  ],
]

// ── Helpers ────────────────────────────────────────────────────────────────
const parseDate = (s: string) => new Date(s.includes('T') ? s : s.replace(' ', 'T'))
const formatDate = (s: string) =>
  parseDate(s).toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' })

// ── Watchers ───────────────────────────────────────────────────────────────
const debouncedSearch = useDebounceFn(() => {
  currentPage.value = 1
  selectedUsers.value = []
  batchAction.value = 'none'
  fetchUsers()
}, 350)
watch(searchQuery, debouncedSearch)
watch([filterStatus, filterRole, sortBy], () => {
  currentPage.value = 1
  selectedUsers.value = []
  batchAction.value = 'none'
  fetchUsers()
})

onMounted(fetchUsers)
</script>
