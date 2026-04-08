<template>
  <div class="min-h-screen pb-16">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-8']">

      <!-- Title -->
      <div class="flex items-center gap-2 mb-6">
        <UIcon name="i-tabler-activity" class="size-6 text-primary shrink-0" />
        <h1 class="text-xl font-bold text-highlighted">{{ $t('site.activity.title') }}</h1>
      </div>

      <!-- Pill tabs + actions -->
      <div class="flex items-center justify-between mb-4 flex-wrap gap-2">
        <div class="flex gap-1 p-1 bg-default rounded-md ring-1 ring-default w-fit">
          <UButton
            v-for="filter in filters"
            :key="filter.value"
            :variant="currentFilter === filter.value ? 'solid' : 'ghost'"
            :color="currentFilter === filter.value ? 'primary' : 'neutral'"
            size="sm"
            :class="currentFilter === filter.value ? 'shadow' : 'text-muted'"
            @click="onFilterChange(filter.value)">
            {{ filter.label }}
            <span
              v-if="filter.value === 'unread' && unreadCount > 0"
              class="ml-1 inline-flex items-center justify-center min-w-[16px] h-4 px-0.5 rounded-full bg-error text-white text-[10px] font-bold">
              {{ unreadCount > 99 ? '99+' : unreadCount }}
            </span>
          </UButton>
        </div>

        <div class="flex gap-1">
          <UButton
            v-if="unreadCount > 0"
            variant="ghost"
            color="neutral"
            size="sm"
            icon="i-tabler-checks"
            @click="markAllAsRead">
            {{ $t('site.activity.mark_all_read') }}
          </UButton>
          <UButton
            v-if="activities.length > 0"
            variant="ghost"
            color="neutral"
            size="sm"
            icon="i-tabler-trash"
            class="text-error"
            @click="clearAll">
            {{ clearLabel }}
          </UButton>
        </div>
      </div>

      <!-- Main card -->
      <div class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">

        <!-- Loading skeleton -->
        <div v-if="displayLoading" class="divide-y divide-default">
          <div v-for="i in 5" :key="i" class="flex items-start gap-3 px-4 py-4">
            <USkeleton class="size-9 rounded-full shrink-0" />
            <div class="flex-1 space-y-2 pt-0.5">
              <div class="flex items-center justify-between">
                <USkeleton class="h-3.5 w-36" />
                <USkeleton class="h-3 w-12" />
              </div>
              <USkeleton class="h-3 w-full" />
              <USkeleton class="h-3 w-3/4" />
            </div>
          </div>
        </div>

        <!-- Empty state -->
        <div
          v-else-if="activities.length === 0"
          class="py-20 flex flex-col items-center gap-3 text-center px-4">
          <div class="size-16 rounded-full bg-muted flex items-center justify-center">
            <UIcon name="i-tabler-bell-off" class="size-8 text-muted" />
          </div>
          <p class="font-semibold text-highlighted">{{ $t('site.activity.no_activity') }}</p>
          <p class="text-sm text-muted">
            {{ currentFilter === 'unread' ? $t('site.activity.no_unread') : $t('site.activity.no_any') }}
          </p>
        </div>

        <!-- Activity rows -->
        <div v-else class="divide-y divide-default">
          <div
            v-for="activity in activities"
            :key="activity.id"
            class="group flex items-start gap-3 px-4 py-3.5 transition-colors"
            :class="activity.read ? 'hover:bg-elevated/50' : 'bg-primary/5 hover:bg-primary/8'">

            <!-- Icon / Avatar -->
            <div class="shrink-0 mt-0.5">
              <div
                v-if="activity.type === 'system'"
                class="size-9 rounded-full flex items-center justify-center"
                :class="getSystemIconClass(activity.sub_type!)">
                <UIcon :name="getSystemIcon(activity.sub_type!)" class="size-4" />
              </div>
              <BaseAvatar
                v-else
                :src="activity.avatar"
                :alt="activity.user_name"
                size="sm"
              />
            </div>

            <!-- Content -->
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between gap-2">
                <div class="flex-1 min-w-0">
                  <!-- User action -->
                  <p v-if="activity.type !== 'system'" class="text-sm text-default leading-snug mb-0.5">
                    <span class="font-semibold text-highlighted">{{ activity.user_name }}</span>
                    <span class="text-muted ml-1">{{ activity.action }}</span>
                    <span v-if="!activity.read" class="inline-block size-1.5 rounded-full bg-primary shrink-0 ml-1.5 align-middle" />
                  </p>
                  <!-- System notification -->
                  <div v-else class="flex items-center gap-1.5 mb-0.5">
                    <p class="text-sm font-semibold text-highlighted leading-snug">{{ activity.title }}</p>
                    <span v-if="!activity.read" class="inline-block size-1.5 rounded-full bg-primary shrink-0" />
                  </div>

                  <p v-if="activity.content" class="text-sm text-default leading-relaxed">{{ activity.content }}</p>

                  <!-- Related content link -->
                  <NuxtLink
                    v-if="activity.related_title && activity.related_link"
                    :to="activity.related_link"
                    class="inline-flex items-center gap-1 text-xs text-primary hover:underline mt-1">
                    <UIcon name="i-tabler-link" class="size-3" />
                    {{ activity.related_title }}
                  </NuxtLink>

                  <p class="text-xs text-muted mt-1.5">{{ formatTime(activity.created_at) }}</p>
                </div>

                <!-- Hover actions -->
                <div class="flex items-center gap-0.5 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
                  <UButton
                    v-if="!activity.read"
                    variant="ghost"
                    color="primary"
                    size="xs"
                    square
                    icon="i-tabler-check"
                    :title="$t('site.activity.mark_read')"
                    @click="markAsRead(activity.id)" />
                  <UButton
                    variant="ghost"
                    color="neutral"
                    size="xs"
                    square
                    icon="i-tabler-trash"
                    class="hover:text-error"
                    :title="$t('site.activity.delete')"
                    @click="deleteActivity(activity.id)" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex justify-center mt-6">
        <UPagination
          v-model:page="currentPage"
          :total="total"
          :items-per-page="pageSize" />
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
const { containerClass } = useContainerWidth()
const { t } = useI18n()
import type { NotificationItem } from '~/composables/useNotificationApi'

// Hardcoded to user 1 (single-user blog owner)
const OWNER_USER_ID = 1

type FilterType = 'all' | 'unread' | 'interaction' | 'system'

const notificationApi = useNotificationApi()

const loading = ref(false)
const displayLoading = useMinLoading(loading)
const currentFilter = ref<FilterType>('all')
const currentPage = ref(1)
const pageSize = 10

const activities = ref<NotificationItem[]>([])
const total = ref(0)
const totalPages = ref(0)
const unreadCount = ref(0)

const filters = computed(() => [
  { label: t('site.activity.filter_all'), value: 'all' },
  { label: t('site.activity.filter_unread'), value: 'unread' },
  { label: t('site.activity.filter_interaction'), value: 'interaction' },
  { label: t('site.activity.filter_system'), value: 'system' },
])

const fetchNotifications = async () => {
  loading.value = true
  try {
    const res = await notificationApi.getNotifications({
      user_id: OWNER_USER_ID,
      filter: currentFilter.value,
      page: currentPage.value,
      size: pageSize,
    })
    activities.value = res.list ?? []
    total.value = res.total
    totalPages.value = res.total_pages
    unreadCount.value = res.unread
  } catch {
    activities.value = []
    total.value = 0
    totalPages.value = 0
  } finally {
    loading.value = false
  }
}

onMounted(fetchNotifications)
watch(currentPage, fetchNotifications)

const onFilterChange = async (filter: FilterType) => {
  currentFilter.value = filter
  currentPage.value = 1
  await fetchNotifications()
}

const formatTime = (dateStr: string) => {
  const date = new Date(dateStr)
  const diff = Date.now() - date.getTime()
  const m = Math.floor(diff / 60000)
  const h = Math.floor(diff / 3600000)
  const d = Math.floor(diff / 86400000)
  if (m < 1) return t('site.activity.time_just_now')
  if (m < 60) return t('site.activity.time_minutes_ago', { n: m })
  if (h < 24) return t('site.activity.time_hours_ago', { n: h })
  if (d < 7) return t('site.activity.time_days_ago', { n: d })
  return date.toLocaleDateString('zh-CN')
}

const getSystemIcon = (subType: string) => {
  const icons: Record<string, string> = {
    approved: 'i-tabler-circle-check',
    rejected: 'i-tabler-circle-x',
    update: 'i-tabler-sparkles',
    security: 'i-tabler-shield-x',
    newsletter: 'i-tabler-mail',
  }
  return icons[subType] || 'i-tabler-bell'
}

const getSystemIconClass = (subType: string) => {
  const classes: Record<string, string> = {
    approved: 'bg-success/10 text-success',
    rejected: 'bg-error/10 text-error',
    update: 'bg-info/10 text-info',
    security: 'bg-warning/10 text-warning',
    newsletter: 'bg-secondary/10 text-secondary',
  }
  return classes[subType] || 'bg-muted text-muted'
}

const markAsRead = async (id: number) => {
  await notificationApi.markRead(id)
  const item = activities.value.find(a => a.id === id)
  if (item) {
    item.read = true
    unreadCount.value = Math.max(0, unreadCount.value - 1)
  }
}

const markAllAsRead = async () => {
  await notificationApi.markAllRead(OWNER_USER_ID)
  activities.value.forEach(a => (a.read = true))
  unreadCount.value = 0
}

const deleteActivity = async (id: number) => {
  if (confirm(t('site.activity.confirm_delete'))) {
    await notificationApi.deleteNotification(id)
    const idx = activities.value.findIndex(a => a.id === id)
    if (idx > -1) {
      const wasUnread = !activities.value[idx].read
      activities.value.splice(idx, 1)
      total.value = Math.max(0, total.value - 1)
      if (wasUnread) unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
  }
}

const clearLabelMap = computed<Record<string, string>>(() => ({
  all: t('site.activity.clear_all'),
  unread: t('site.activity.clear_unread'),
  interaction: t('site.activity.clear_interaction'),
  system: t('site.activity.clear_system'),
}))
const clearLabel = computed(() => clearLabelMap.value[currentFilter.value] ?? t('site.activity.clear_default'))

const clearAll = async () => {
  if (confirm(t('site.activity.confirm_clear', { label: clearLabel.value }))) {
    await notificationApi.clearNotifications(OWNER_USER_ID, currentFilter.value)
    await fetchNotifications()
  }
}
</script>
