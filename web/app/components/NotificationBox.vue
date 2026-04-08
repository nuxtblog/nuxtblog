<template>
  <UPopover>
    <div class="relative">
      <UButton
        icon="i-tabler-bell"
        variant="ghost"
        color="neutral"
        square
        size="sm"
        :aria-label="$t('site.notifications.title')" />
      <span
        v-if="unreadCount > 0"
        class="absolute -top-0.5 -right-0.5 min-w-4 h-4 px-0.5 rounded-full bg-error text-white text-[10px] font-bold flex items-center justify-center pointer-events-none">
        {{ unreadCount > 99 ? "99+" : unreadCount }}
      </span>
    </div>

    <template #content>
      <div class="w-80">
        <!-- 头部 -->
        <div class="px-4 py-3 border-b border-default flex items-center justify-between">
          <h3 class="font-semibold text-highlighted">{{ $t("site.notifications.title") }}</h3>
          <UButton
            v-if="unreadCount > 0"
            variant="ghost"
            color="primary"
            size="xs"
            :loading="markingAll"
            @click="markAllRead">
            {{ $t("site.notifications.mark_all_read") }}
          </UButton>
        </div>

        <!-- 列表 -->
        <ul class="max-h-96 overflow-y-auto divide-y divide-default">
          <li v-if="pending" class="px-4 py-6 text-center">
            <UIcon name="i-tabler-loader-2" class="size-5 text-muted animate-spin mx-auto" />
          </li>

          <li v-else-if="!mergedList.length" class="px-4 py-8 text-center">
            <UIcon name="i-tabler-bell-off" class="size-8 text-muted mx-auto mb-2" />
            <p class="text-sm text-muted">{{ $t("site.notifications.empty") }}</p>
          </li>

          <li
            v-for="item in mergedList"
            :key="(item._isAnnouncement ? 'a-' : 'n-') + item.id"
            class="px-4 py-3 hover:bg-elevated/50 transition-colors cursor-pointer"
            @click="handleClick(item)">
            <div class="flex gap-3">
              <!-- 图标 -->
              <div
                class="shrink-0 size-7 rounded-full flex items-center justify-center mt-0.5"
                :class="getItemIconBg(item)">
                <UIcon :name="getItemIcon(item)" class="size-3.5" :class="getItemIconColor(item)" />
              </div>

              <div class="flex-1 min-w-0">
                <div class="flex items-start gap-1.5">
                  <p class="text-sm font-medium text-highlighted leading-snug flex-1">
                    {{ item.title || item.content }}
                  </p>
                  <span v-if="!item.read" class="size-1.5 rounded-full bg-primary shrink-0 mt-1.5" />
                </div>
                <p v-if="item.title && item.content" class="text-xs text-muted mt-0.5 line-clamp-1">
                  {{ item.content }}
                </p>
                <p v-if="item.user_name" class="text-xs text-muted mt-0.5">{{ item.user_name }}</p>
                <p class="text-xs text-muted mt-1">{{ formatTime(item.created_at) }}</p>
              </div>
            </div>
          </li>
        </ul>

        <!-- 底部 -->
        <div class="p-2 border-t border-default">
          <UButton to="/notifications" variant="ghost" color="primary" block size="sm">
            {{ $t("site.notifications.view_all") }}
          </UButton>
        </div>
      </div>
    </template>
  </UPopover>
</template>

<script setup lang="ts">
import type { NotificationItem } from "~/composables/useNotificationApi";
import type { AnnouncementItem } from "~/composables/useAnnouncementApi";

const { t } = useI18n();
const authStore = useAuthStore();
const notificationApi = useNotificationApi();
const announcementApi = useAnnouncementApi();

const markingAll = ref(false);

type BoxItem = {
  id: number
  type: string
  title?: string
  content: string
  user_name?: string
  related_link?: string
  read: boolean
  created_at: string
  _isAnnouncement: boolean
  _annType?: string
}

const { data, pending, refresh } = await useAsyncData(
  "notifications-box",
  async () => {
    if (!authStore.user?.id) return null;
    const [notifRes, annRes] = await Promise.all([
      notificationApi.getNotifications({ user_id: authStore.user!.id, page: 1, size: 8 }).catch(() => null),
      announcementApi.getAnnouncements({ page: 1, size: 3 }).catch(() => null),
    ]);
    return { notifRes, annRes };
  },
);

const mergedList = computed((): BoxItem[] => {
  const notifications: BoxItem[] = (data.value?.notifRes?.list ?? []).map((n: NotificationItem) => ({
    id: n.id,
    type: n.type,
    title: n.title,
    content: n.content,
    user_name: n.user_name,
    related_link: n.related_link,
    read: n.read,
    created_at: n.created_at,
    _isAnnouncement: false,
  }));
  const announcements: BoxItem[] = (data.value?.annRes?.list ?? []).map((a: AnnouncementItem) => ({
    id: a.id,
    type: "announcement",
    title: a.title,
    content: a.content,
    read: !a.unread,
    created_at: a.created_at,
    _isAnnouncement: true,
    _annType: a.type,
  }));
  // announcements first, then notifications; total max 10
  return [...announcements, ...notifications].slice(0, 10);
});

// unread = notification unreads + announcement unreads
const unreadCount = computed(() =>
  (data.value?.notifRes?.unread ?? 0)
);

const handleClick = async (item: BoxItem) => {
  if (item._isAnnouncement) {
    if (!item.read) {
      await announcementApi.markAnnouncementsRead().catch(() => {});
      refresh();
    }
    return;
  }
  if (!item.read) {
    await notificationApi.markRead(item.id).catch(() => {});
    refresh();
  }
  if (item.related_link) {
    await navigateTo(item.related_link);
  }
};

const markAllRead = async () => {
  if (!authStore.user?.id) return;
  markingAll.value = true;
  await Promise.all([
    notificationApi.markAllRead(authStore.user.id).catch(() => {}),
    announcementApi.markAnnouncementsRead().catch(() => {}),
  ]);
  await refresh();
  markingAll.value = false;
};

// ── Icon helpers ──────────────────────────────────────────────────────────────

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
  info: "bg-info/10", warning: "bg-warning/10", success: "bg-success/10", danger: "bg-error/10",
};
const ANN_COLOR: Record<string, string> = {
  info: "text-info", warning: "text-warning", success: "text-success", danger: "text-error",
};

const getItemIcon = (item: BoxItem) =>
  item._isAnnouncement ? (ANN_ICON[item._annType ?? "info"] ?? ANN_ICON.info!) : (TYPE_ICONS[item.type] ?? TYPE_ICONS.system!);
const getItemIconBg = (item: BoxItem) =>
  item._isAnnouncement ? (ANN_BG[item._annType ?? "info"] ?? ANN_BG.info!) : (TYPE_BG[item.type] ?? TYPE_BG.system!);
const getItemIconColor = (item: BoxItem) =>
  item._isAnnouncement ? (ANN_COLOR[item._annType ?? "info"] ?? ANN_COLOR.info!) : (TYPE_COLOR[item.type] ?? TYPE_COLOR.system!);

const formatTime = (iso: string) => {
  const diff = Date.now() - new Date(iso).getTime();
  const m = Math.floor(diff / 60000);
  if (m < 1) return t("common.just_now");
  if (m < 60) return t("common.minutes_ago", { n: m });
  const h = Math.floor(m / 60);
  if (h < 24) return t("common.hours_ago", { n: h });
  return t("common.days_ago", { n: Math.floor(h / 24) });
};
</script>
