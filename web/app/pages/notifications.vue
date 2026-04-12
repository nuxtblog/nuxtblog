<template>
  <div class="min-h-screen pb-16">
    <PageHeader
      icon="i-tabler-bell"
      :title="$t('site.notifications.title')">
      <template #actions>
        <div class="flex gap-1">
          <UButton
            v-if="unreadCount > 0"
            variant="ghost"
            color="neutral"
            size="sm"
            icon="i-tabler-checks"
            @click="handleMarkAllRead">
            {{ $t("site.notifications.mark_all_read") }}
          </UButton>
          <UButton
            v-if="list.length > 0 && currentFilter !== 'announcement'"
            variant="ghost"
            color="neutral"
            size="sm"
            icon="i-tabler-trash"
            class="text-error"
            @click="handleClearAll">
            {{ $t("site.notifications.clear") }}
          </UButton>
        </div>
      </template>
      <template #toolbar>
        <div class="flex gap-1 p-1 bg-default rounded-md ring-1 ring-default w-fit">
          <UButton
            v-for="tab in filterItems"
            :key="tab.value"
            :variant="currentFilter === tab.value ? 'solid' : 'ghost'"
            :color="currentFilter === tab.value ? 'primary' : 'neutral'"
            size="sm"
            :class="currentFilter === tab.value ? 'shadow' : 'text-muted'"
            @click="currentFilter = tab.value as any">
            {{ tab.label }}
            <span
              v-if="tab.value === 'unread' && notifUnreadOnly > 0"
              class="ml-1 inline-flex items-center justify-center min-w-[16px] h-4 px-0.5 rounded-full bg-error text-white text-[10px] font-bold">
              {{ notifUnreadOnly > 99 ? "99+" : notifUnreadOnly }}
            </span>
            <span
              v-if="tab.value === 'announcement' && announcementUnread > 0"
              class="ml-1 inline-flex items-center justify-center min-w-[16px] h-4 px-0.5 rounded-full bg-warning text-white text-[10px] font-bold">
              {{ announcementUnread > 99 ? "99+" : announcementUnread }}
            </span>
          </UButton>
        </div>
      </template>
    </PageHeader>

    <PageContent>
      <!-- 统计行 -->
      <div class="grid grid-cols-3 rounded-md overflow-hidden ring-1 ring-default bg-default shadow-sm mb-6">
        <div class="flex flex-col items-center py-4">
          <span class="text-2xl font-bold text-highlighted">{{ totalAll }}</span>
          <span class="text-xs text-muted mt-0.5">{{ $t("site.notifications.stat_all") }}</span>
        </div>
        <div class="flex flex-col items-center py-4 border-x border-default">
          <span class="text-2xl font-bold" :class="unreadCount > 0 ? 'text-primary' : 'text-highlighted'">{{ unreadCount }}</span>
          <span class="text-xs text-muted mt-0.5">{{ $t("site.notifications.stat_unread") }}</span>
        </div>
        <div class="flex flex-col items-center py-4">
          <span class="text-2xl font-bold text-highlighted">{{ Math.max(0, totalAll - unreadCount) }}</span>
          <span class="text-xs text-muted mt-0.5">{{ $t("site.notifications.stat_read") }}</span>
        </div>
      </div>

      <!-- 通知列表 -->
      <BaseCardList :loading="displayLoading" :empty="displayList.length === 0" :skeleton-count="5">

        <template #skeleton>
          <div class="flex items-center gap-3">
            <USkeleton class="size-10 rounded-full shrink-0" />
            <div class="flex-1 space-y-2">
              <div class="flex items-center justify-between">
                <USkeleton class="h-3.5 w-32" />
                <USkeleton class="h-3 w-12" />
              </div>
              <USkeleton class="h-3 w-full" />
              <USkeleton class="h-3 w-2/3" />
            </div>
          </div>
        </template>

        <template #empty>
          <div class="text-center py-12">
            <UIcon name="i-tabler-bell-off" class="size-12 mx-auto mb-3 text-muted" />
            <p class="font-semibold text-highlighted mb-1">{{ $t("site.notifications.no_notifications") }}</p>
            <p class="text-sm text-muted">
              {{ currentFilter === "unread" ? $t("site.notifications.no_unread") : $t("site.notifications.no_any") }}
            </p>
          </div>
        </template>

        <TransitionGroup tag="div" name="notif" class="space-y-3">
          <UCard
            v-for="item in displayList"
            :key="item._isAnnouncement ? 'a-' + item.id : String(item.id)"
            class="group hover:shadow-md transition-shadow"
            :class="item.related_link && !item._isAnnouncement && 'cursor-pointer'"
            @click="handleItemClick(item)">
            <div class="flex items-start gap-3">

              <!-- 类型图标 -->
              <div
                class="shrink-0 size-10 rounded-full flex items-center justify-center"
                :class="getItemIconBg(item)">
                <UIcon :name="getItemIcon(item)" class="size-4" :class="getItemIconColor(item)" />
              </div>

              <!-- 内容 -->
              <div class="flex-1 min-w-0">
                <div class="flex items-start justify-between gap-2">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-1.5 flex-wrap mb-0.5">
                      <p class="text-sm font-semibold text-highlighted leading-snug">
                        {{ item.title || getTypeLabel(item.type) }}
                      </p>
                      <UBadge
                        v-if="item._isAnnouncement"
                        :label="$t('site.notifications.type_announcement')"
                        :color="annBadgeColor(item._annType) as any"
                        variant="subtle"
                        size="xs" />
                      <span v-if="!item.read" class="size-1.5 rounded-full bg-primary shrink-0" />
                    </div>

                    <p v-if="item.user_name && !item._isAnnouncement" class="text-xs text-muted mb-1">
                      {{ item.user_name }}
                    </p>
                    <p v-if="item.content" class="text-sm text-muted leading-relaxed line-clamp-2">
                      {{ item.content }}
                    </p>
                    <div
                      v-if="item.related_link && item.related_title"
                      class="inline-flex items-center gap-1 text-xs text-primary mt-1.5">
                      <UIcon name="i-tabler-link" class="size-3 shrink-0" />
                      <span class="truncate">{{ item.related_title }}</span>
                    </div>
                    <p class="text-xs text-muted mt-1.5">{{ formatTime(item.created_at) }}</p>
                  </div>

                  <!-- Hover 操作（公告不显示删除） -->
                  <div
                    v-if="!item._isAnnouncement"
                    class="flex items-center gap-0.5 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity">
                    <UButton
                      v-if="!item.read"
                      variant="ghost"
                      color="primary"
                      size="xs"
                      square
                      icon="i-tabler-check"
                      :title="$t('site.notifications.mark_read_title')"
                      @click.stop="handleMarkRead(item)" />
                    <UButton
                      variant="ghost"
                      color="neutral"
                      size="xs"
                      square
                      icon="i-tabler-trash"
                      :title="$t('site.notifications.delete_title')"
                      class="hover:text-error"
                      @click.stop="handleDelete(item)" />
                  </div>
                </div>
              </div>

            </div>
          </UCard>
        </TransitionGroup>

      </BaseCardList>
    </PageContent>

    <PageFooter v-if="totalPages > 1">
      <div class="flex justify-center">
        <UPagination
          v-model:page="currentPage"
          :total="currentFilter === 'announcement' ? announcementTotal : currentTotal"
          :items-per-page="pageSize" />
      </div>
    </PageFooter>
  </div>
</template>

<script setup lang="ts">
import type { NotificationItem } from "~/composables/useNotificationApi";
import type { AnnouncementItem } from "~/composables/useAnnouncementApi";

const { t } = useI18n();

definePageMeta({ middleware: "auth" });

const authStore = useAuthStore();
const notifApi = useNotificationApi();
const announcementApi = useAnnouncementApi();
const toast = useToast();

// ── State ─────────────────────────────────────────────────────────────────────

type AnnType = 'info' | 'warning' | 'success' | 'danger'
type DisplayItem = NotificationItem & {
  _isAnnouncement: boolean
  _annType?: AnnType
}

const rawLoading = ref(true);
const displayLoading = useMinLoading(rawLoading);
const currentFilter = ref<"all" | "unread" | "read" | "announcement">("all");
const currentPage = ref(1);
const pageSize = 20;

const list = ref<NotificationItem[]>([]);
const currentTotal = ref(0);
const totalAll = ref(0);
const unreadCount = ref(0);

const announcementList = ref<AnnouncementItem[]>([]);
const announcementTotal = ref(0);
const announcementUnread = ref(0);

// ── Computed ──────────────────────────────────────────────────────────────────

// unreadCount from backend includes announcement unreads; strip them out for the notification-only tab
const notifUnreadOnly = computed(() => Math.max(0, unreadCount.value - announcementUnread.value));

const totalPages = computed(() =>
  Math.ceil(
    (currentFilter.value === "announcement" ? announcementTotal.value : currentTotal.value) / pageSize,
  ),
);

const filterItems = computed(() => [
  { label: t("site.notifications.filter_all"), value: "all" },
  { label: t("site.notifications.filter_unread"), value: "unread" },
  { label: t("site.notifications.filter_read"), value: "read" },
  { label: t("site.notifications.filter_announcement"), value: "announcement" },
]);

const displayList = computed((): DisplayItem[] => {
  const toDisplay = (n: NotificationItem): DisplayItem => ({ ...n, _isAnnouncement: false });
  const annToDisplay = (a: AnnouncementItem): DisplayItem => ({
    id: a.id,
    type: "announcement" as any,
    content: a.content,
    title: a.title,
    read: !a.unread,
    created_at: a.created_at,
    _isAnnouncement: true,
    _annType: a.type as AnnType,
  });

  if (currentFilter.value === "announcement") return announcementList.value.map(annToDisplay);
  if (currentFilter.value === "all")
    return [...announcementList.value.map(annToDisplay), ...list.value.map(toDisplay)];
  return list.value.map(toDisplay);
});

// ── Fetch ─────────────────────────────────────────────────────────────────────

async function fetchNotifications() {
  if (!authStore.user?.id) return;
  rawLoading.value = true;
  try {
    if (currentFilter.value === "announcement") {
      const res = await announcementApi.getAnnouncements({ page: currentPage.value, size: pageSize });
      announcementList.value = res.list ?? [];
      announcementTotal.value = res.total;
      announcementUnread.value = res.unread_count;
    } else {
      const [notifRes, announcRes] = await Promise.all([
        notifApi.getNotifications({
          user_id: authStore.user!.id,
          filter: currentFilter.value,
          page: currentPage.value,
          size: pageSize,
        }),
        currentFilter.value === "all"
          ? announcementApi.getAnnouncements({ page: 1, size: 5 })
          : Promise.resolve(null),
      ]);
      list.value = notifRes.list ?? [];
      currentTotal.value = notifRes.total;
      // unread from backend already includes announcement unreads
      unreadCount.value = notifRes.unread;
      if (announcRes) {
        announcementList.value = announcRes.list ?? [];
        announcementTotal.value = announcRes.total;
        announcementUnread.value = announcRes.unread_count;
        // totalAll = notifications + announcements
        totalAll.value = notifRes.total + announcRes.total;
      } else if (currentFilter.value === "unread" || currentFilter.value === "read") {
        totalAll.value = currentTotal.value;
      }
    }
  } catch {
    toast.add({ title: t("site.notifications.load_failed"), color: "error" });
  } finally {
    rawLoading.value = false;
  }
}

watch(currentFilter, async () => {
  currentPage.value = 1;
  await fetchNotifications();
  // Switching to announcement tab = user has seen the announcements → mark all as read
  if (currentFilter.value === "announcement" && announcementUnread.value > 0) {
    const n = announcementUnread.value;
    announcementList.value.forEach((a) => (a.unread = false));
    announcementUnread.value = 0;
    unreadCount.value = Math.max(0, unreadCount.value - n);
    announcementApi.markAnnouncementsRead().catch(() => {});
  }
});
watch(currentPage, fetchNotifications);
onMounted(fetchNotifications);

// ── Actions ───────────────────────────────────────────────────────────────────

async function handleItemClick(item: DisplayItem) {
  if (item._isAnnouncement) return;
  if (!item.related_link) return;
  if (!item.read) await handleMarkRead(item);
  await navigateTo(item.related_link);
}

async function handleMarkRead(item: DisplayItem) {
  if (item._isAnnouncement) return;
  const prev = item.read;
  item.read = true;
  unreadCount.value = Math.max(0, unreadCount.value - 1);
  try {
    await notifApi.markRead(item.id);
  } catch {
    item.read = prev;
    unreadCount.value++;
  }
}

async function handleMarkAllRead() {
  if (!authStore.user?.id) return;
  const prev = list.value.map((n) => n.read);
  const prevUnread = unreadCount.value;
  const prevAnnUnread = announcementUnread.value;
  list.value.forEach((n) => (n.read = true));
  announcementList.value.forEach((a) => (a.unread = false));
  unreadCount.value = 0;
  announcementUnread.value = 0;
  try {
    await Promise.all([
      notifApi.markAllRead(authStore.user.id),
      announcementApi.markAnnouncementsRead(),
    ]);
  } catch {
    list.value.forEach((n, i) => (n.read = prev[i]!));
    unreadCount.value = prevUnread;
    announcementUnread.value = prevAnnUnread;
    toast.add({ title: t("site.notifications.op_failed"), color: "error" });
  }
}

async function handleDelete(item: DisplayItem) {
  const idx = list.value.findIndex((n) => n.id === item.id);
  if (idx === -1) return;
  const [removed] = list.value.splice(idx, 1);
  currentTotal.value--;
  totalAll.value = Math.max(0, totalAll.value - 1);
  if (!removed!.read) unreadCount.value = Math.max(0, unreadCount.value - 1);
  try {
    await notifApi.deleteNotification(item.id);
  } catch {
    list.value.splice(idx, 0, removed!);
    currentTotal.value++;
    totalAll.value++;
    if (!removed!.read) unreadCount.value++;
    toast.add({ title: t("site.notifications.delete_failed"), color: "error" });
  }
}

async function handleClearAll() {
  if (!authStore.user?.id || currentFilter.value === "announcement") return;
  if (!confirm(t("site.notifications.confirm_clear"))) return;
  try {
    await notifApi.clearNotifications(authStore.user.id, currentFilter.value);
    await fetchNotifications();
  } catch {
    toast.add({ title: t("site.notifications.clear_failed"), color: "error" });
  }
}

// ── Helpers ───────────────────────────────────────────────────────────────────

const TYPE_ICONS: Record<string, string> = {
  follow: "i-tabler-user-plus",
  like: "i-tabler-heart",
  comment: "i-tabler-message-circle",
  reply: "i-tabler-message-reply",
  mention: "i-tabler-at",
  system: "i-tabler-bell",
};

const TYPE_BG: Record<string, string> = {
  follow: "bg-primary/10",
  like: "bg-rose-500/10",
  comment: "bg-info/10",
  reply: "bg-info/10",
  mention: "bg-warning/10",
  system: "bg-elevated",
};

const TYPE_COLOR: Record<string, string> = {
  follow: "text-primary",
  like: "text-rose-500",
  comment: "text-info",
  reply: "text-info",
  mention: "text-warning",
  system: "text-muted",
};

const ANN_ICON: Record<string, string> = {
  info: "i-tabler-info-circle",
  warning: "i-tabler-alert-triangle",
  success: "i-tabler-circle-check",
  danger: "i-tabler-alert-circle",
};
const ANN_BG: Record<string, string> = {
  info: "bg-info/10",
  warning: "bg-warning/10",
  success: "bg-success/10",
  danger: "bg-error/10",
};
const ANN_COLOR: Record<string, string> = {
  info: "text-info",
  warning: "text-warning",
  success: "text-success",
  danger: "text-error",
};

function getItemIcon(item: DisplayItem): string {
  if (item._isAnnouncement) return ANN_ICON[item._annType ?? "info"] ?? ANN_ICON.info!;
  return TYPE_ICONS[item.type] ?? TYPE_ICONS.system!;
}
function getItemIconBg(item: DisplayItem): string {
  if (item._isAnnouncement) return ANN_BG[item._annType ?? "info"] ?? ANN_BG.info!;
  return TYPE_BG[item.type] ?? TYPE_BG.system!;
}
function getItemIconColor(item: DisplayItem): string {
  if (item._isAnnouncement) return ANN_COLOR[item._annType ?? "info"] ?? ANN_COLOR.info!;
  return TYPE_COLOR[item.type] ?? TYPE_COLOR.system!;
}
function annBadgeColor(annType?: string): string {
  const map: Record<string, string> = { info: "info", warning: "warning", success: "success", danger: "error" };
  return map[annType ?? "info"] ?? "neutral";
}

const TYPE_LABELS = computed<Record<string, string>>(() => ({
  follow: t("site.notifications.type_follow"),
  like: t("site.notifications.type_like"),
  comment: t("site.notifications.type_comment"),
  reply: t("site.notifications.type_reply"),
  mention: t("site.notifications.type_mention"),
  system: t("site.notifications.type_system"),
  announcement: t("site.notifications.type_announcement"),
}));

const getTypeLabel = (type: string) => TYPE_LABELS.value[type] ?? t("site.notifications.type_default");

const formatTime = (str: string) => {
  const date = new Date(str);
  const diff = Date.now() - date.getTime();
  const m = Math.floor(diff / 60000);
  const h = Math.floor(diff / 3600000);
  const d = Math.floor(diff / 86400000);
  if (m < 1) return t("site.activity.time_just_now");
  if (m < 60) return t("site.activity.time_minutes_ago", { n: m });
  if (h < 24) return t("site.activity.time_hours_ago", { n: h });
  if (d < 7) return t("site.activity.time_days_ago", { n: d });
  return date.toLocaleDateString("zh-CN");
};
</script>

<style scoped>
.notif-enter-active,
.notif-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}
.notif-enter-from { opacity: 0; transform: translateY(-4px); }
.notif-leave-to   { opacity: 0; transform: translateX(8px); }
</style>
