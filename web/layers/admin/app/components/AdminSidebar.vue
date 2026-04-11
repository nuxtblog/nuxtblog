<template>
  <aside
    class="flex flex-col h-full bg-default border-r border-default transition-all duration-300"
    :class="desktopCollapsed ? 'w-14' : 'w-52'">
    <!-- 移动端关闭按钮 -->
    <div
      class="flex items-center justify-between px-3 py-3 lg:hidden border-b border-default">
      <span class="text-sm font-semibold text-highlighted">{{
        t("admin.menu")
      }}</span>
      <UButton
        color="neutral"
        variant="ghost"
        icon="i-tabler-x"
        square
        size="xs"
        @click="closeMobile" />
    </div>

    <!-- 菜单列表 -->
    <nav class="flex-1 overflow-y-auto py-3 px-2 space-y-0.5">
      <template v-for="item in menu" :key="item.name">
        <!-- 无子菜单 -->
        <template v-if="!item.children">
          <!-- 展开状态 -->
          <NuxtLink
            v-if="!desktopCollapsed"
            :to="item.to"
            class="flex items-center gap-2.5 px-2.5 py-2 rounded-md text-sm transition-colors hover:bg-elevated"
            :class="
              route.path.startsWith(item.to!)
                ? 'bg-primary/10 text-primary font-medium'
                : 'text-muted'
            "
            @click="closeMobile">
            <UIcon :name="item.icon!" class="size-4.5 shrink-0" />
            <span class="flex-1 truncate">{{ t(item.label) }}</span>
            <UBadge
              v-if="item.badge && item.badge() > 0"
              :label="String(item.badge())"
              size="xs"
              color="neutral"
              variant="soft" />
          </NuxtLink>

          <!-- 折叠状态（图标模式） -->
          <UTooltip
            v-else
            :text="t(item.label)"
            :popper="{ placement: 'right' }">
            <NuxtLink
              :to="item.to"
              class="flex items-center justify-center w-10 h-10 rounded-md mx-auto transition-colors relative hover:bg-elevated"
              :class="
                route.path.startsWith(item.to!)
                  ? 'bg-primary/10 text-primary'
                  : 'text-muted'
              "
              @click="closeMobile">
              <UIcon :name="item.icon!" class="size-4.5" />
              <span
                v-if="item.badge && item.badge() > 0"
                class="absolute -top-0.5 -right-0.5 h-3.5 w-3.5 rounded-full bg-error text-white text-[9px] flex items-center justify-center">
                {{ item.badge() }}
              </span>
            </NuxtLink>
          </UTooltip>
        </template>

        <!-- 有子菜单 - 展开状态 -->
        <details
          v-if="item.children && !desktopCollapsed"
          :open="isParentActive(item)"
          class="group/detail">
          <summary
            class="flex items-center gap-2.5 px-2.5 py-2 rounded-md cursor-pointer text-sm transition-colors hover:bg-elevated list-none"
            :class="
              isParentActive(item)
                ? 'text-highlighted font-medium'
                : 'text-muted'
            ">
            <UIcon :name="item.icon!" class="size-4.5 shrink-0" />
            <span class="flex-1 truncate">{{ t(item.label) }}</span>
            <UIcon
              name="i-tabler-chevron-right"
              class="size-3.5 transition-transform duration-200 group-open/detail:rotate-90" />
          </summary>

          <div class="mt-0.5 ml-3 pl-3 border-l border-default space-y-0.5">
            <NuxtLink
              v-for="child in item.children"
              :key="child.name"
              :to="child.to"
              class="flex items-center gap-2 px-2 py-1.5 rounded-md text-sm transition-colors hover:bg-elevated"
              :class="
                route.path === child.to
                  ? 'text-primary font-medium'
                  : 'text-muted'
              "
              @click="closeMobile">
              {{ t(child.label) }}
            </NuxtLink>
          </div>
        </details>

        <!-- 有子菜单 - 折叠状态 -->
        <UDropdownMenu
          v-if="item.children && desktopCollapsed"
          :items="item.children.map((c) => ({ label: t(c.label), to: c.to }))"
          :popper="{ placement: 'right-start' }">
          <UTooltip :text="t(item.label)" :popper="{ placement: 'right' }">
            <button
              class="flex items-center justify-center w-10 h-10 rounded-md mx-auto transition-colors hover:bg-elevated"
              :class="
                isParentActive(item)
                  ? 'bg-primary/10 text-primary'
                  : 'text-muted'
              ">
              <UIcon :name="item.icon!" class="size-4.5" />
            </button>
          </UTooltip>
        </UDropdownMenu>
      </template>

      <!-- Phase 2.3: Plugin navigation items -->
      <div v-if="!desktopCollapsed" class="mt-1 pt-1 border-t border-default">
        <ContributionSlot name="admin:sidebar-nav" />
      </div>
    </nav>

    <!-- 底部：折叠切换按钮 -->
    <div class="shrink-0 border-t border-default px-2 py-3 hidden lg:block">
      <UButton
        color="neutral"
        variant="ghost"
        :class="
          desktopCollapsed ? 'w-10 mx-auto' : 'w-full justify-start gap-2.5'
        "
        :square="desktopCollapsed"
        size="sm"
        @click="toggleDesktop">
        <UIcon
          :name="
            desktopCollapsed
              ? 'i-tabler-layout-sidebar-right-expand'
              : 'i-tabler-layout-sidebar-left-collapse'
          "
          class="size-4.5 shrink-0" />
        <span v-if="!desktopCollapsed" class="text-sm">{{
          t("admin.collapse")
        }}</span>
      </UButton>

      <div v-if="!desktopCollapsed" class="mt-3 px-2.5 text-xs text-muted">
        <p>{{ t("admin.version") }} v1.0.0</p>
        <p>{{ t("admin.copyright") }}</p>
      </div>
    </div>
  </aside>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router";

interface MenuItem {
  name: string;
  label: string;
  icon?: string;
  to?: string;
  children?: MenuItem[];
  badge?: () => number;
  cap?: import("~/config/permissions").Capability; // required capability to show
}

const route = useRoute();
const { t } = useI18n();
const { desktopCollapsed, closeMobile, toggleDesktop } = useAdminSidebar();
const { can } = usePermissions();

const mediaCount = ref(0);
const pendingComments = ref(0);

const ALL_MENU: MenuItem[] = [
  {
    name: "dashboard",
    label: "admin.nav.dashboard",
    icon: "i-tabler-layout-grid",
    to: "/admin/dashboard",
  },
  {
    name: "posts",
    label: "admin.nav.posts",
    icon: "i-tabler-file-text",
    cap: "edit_posts",
    children: [
      {
        name: "posts-all",
        label: "admin.nav.posts_all",
        to: "/admin/posts",
        icon: "i-tabler-list",
      },
      {
        name: "posts-new",
        label: "admin.nav.posts_new",
        to: "/admin/posts/new",
        icon: "i-tabler-plus",
        cap: "publish_posts",
      },
      {
        name: "posts-categories",
        label: "admin.nav.posts_categories",
        to: "/admin/posts/categories",
        icon: "i-tabler-folder",
        cap: "manage_categories",
      },
      {
        name: "posts-tags",
        label: "admin.nav.posts_tags",
        to: "/admin/posts/tags",
        icon: "i-tabler-tag",
        cap: "manage_categories",
      },
    ],
  },
  {
    name: "docs",
    label: "admin.nav.docs",
    icon: "i-tabler-books",
    cap: "edit_posts",
    children: [
      { name: "docs-all", label: "admin.nav.docs_all", to: "/admin/docs", icon: "i-tabler-list" },
      { name: "docs-new", label: "admin.nav.docs_new", to: "/admin/docs/new", icon: "i-tabler-plus", cap: "publish_posts" },
      { name: "docs-collections", label: "admin.nav.docs_collections", to: "/admin/docs/collections", icon: "i-tabler-folders", cap: "manage_categories" },
    ],
  },
  {
    name: "moments",
    label: "admin.nav.moments",
    icon: "i-tabler-camera",
    to: "/admin/moments",
    cap: "moderate_comments",
  },
  {
    name: "media",
    label: "admin.nav.media",
    icon: "i-tabler-photo",
    to: "/admin/media",
    cap: "upload_files",
    badge: () => Math.min(mediaCount.value, 99),
  },
  {
    name: "pages",
    label: "admin.nav.pages",
    icon: "i-tabler-file",
    cap: "edit_pages",
    children: [
      {
        name: "pages-all",
        label: "admin.nav.pages_all",
        to: "/admin/pages",
        icon: "i-tabler-list",
      },
      {
        name: "pages-new",
        label: "admin.nav.pages_new",
        to: "/admin/pages/new",
        icon: "i-tabler-plus",
        cap: "publish_pages",
      },
    ],
  },
  {
    name: "comments",
    label: "admin.nav.comments",
    icon: "i-tabler-message",
    to: "/admin/comments",
    cap: "moderate_comments",
    badge: () => pendingComments.value,
  },
  {
    name: "reports",
    label: "admin.nav.reports",
    icon: "i-tabler-flag",
    to: "/admin/reports",
    cap: "moderate_comments",
  },
  {
    name: "announcements",
    label: "admin.nav.announcements",
    icon: "i-tabler-speakerphone",
    to: "/admin/announcements",
    cap: "moderate_comments",
  },
  {
    name: "friendlinks",
    label: "admin.nav.friendlinks",
    icon: "i-tabler-link",
    to: "/admin/friendlinks",
    cap: "manage_options",
  },
  {
    name: "users",
    label: "admin.nav.users",
    icon: "i-tabler-users-group",
    cap: "list_users",
    children: [
      {
        name: "users-all",
        label: "admin.nav.users_all",
        to: "/admin/users",
        icon: "i-tabler-list",
        cap: "list_users",
      },
      {
        name: "users-new",
        label: "admin.nav.users_new",
        to: "/admin/users/new",
        icon: "i-tabler-user-plus",
        cap: "create_users",
      },
      {
        name: "users-roles",
        label: "admin.nav.users_roles",
        to: "/admin/users/roles",
        icon: "i-tabler-shield-check",
        cap: "manage_roles",
      },
    ],
  },
  {
    name: "appearance",
    label: "admin.nav.appearance",
    icon: "i-tabler-palette",
    cap: "manage_appearance",
    children: [
      {
        name: "appearance-themes",
        label: "admin.nav.appearance_themes",
        to: "/admin/appearance/themes",
        icon: "i-tabler-paint",
      },
      {
        name: "appearance-customize",
        label: "admin.nav.appearance_customize",
        to: "/admin/appearance/customize",
        icon: "i-tabler-settings",
      },
      {
        name: "appearance-menus",
        label: "admin.nav.appearance_menus",
        to: "/admin/appearance/menus",
        icon: "i-tabler-menu",
      },
    ],
  },
  {
    name: "ai",
    label: "admin.nav.ai",
    icon: "i-tabler-sparkles",
    to: "/admin/ai",
    cap: "manage_options",
  },
  {
    name: "plugins",
    label: "admin.nav.plugins",
    icon: "i-tabler-plug",
    cap: "manage_options",
    children: [
      {
        name: "plugins-installed",
        label: "admin.nav.plugins_installed",
        to: "/admin/plugins",
        icon: "i-tabler-puzzle",
      },
      {
        name: "plugins-marketplace",
        label: "admin.nav.plugins_marketplace",
        to: "/admin/plugins/marketplace",
        icon: "i-tabler-building-store",
      },
      {
        name: "plugins-monitor",
        label: "admin.nav.plugins_monitor",
        to: "/admin/plugins/monitor",
        icon: "i-tabler-activity",
      },
    ],
  },
  {
    name: "system",
    label: "admin.nav.system",
    icon: "i-tabler-server",
    to: "/admin/system",
    cap: "manage_options",
  },
  {
    name: "settings",
    label: "admin.nav.settings",
    icon: "i-tabler-settings",
    cap: "manage_options",
    children: [
      {
        name: "settings-homepage",
        label: "admin.nav.settings_homepage",
        to: "/admin/settings/homepage",
        icon: "i-tabler-home",
      },
      {
        name: "settings-general",
        label: "admin.nav.settings_general",
        to: "/admin/settings/general",
        icon: "i-tabler-adjustments",
      },
      {
        name: "settings-writing",
        label: "admin.nav.settings_writing",
        to: "/admin/settings/writing",
        icon: "i-tabler-edit",
      },
      {
        name: "settings-reading",
        label: "admin.nav.settings_reading",
        to: "/admin/settings/reading",
        icon: "i-tabler-book",
      },
      {
        name: "settings-discussion",
        label: "admin.nav.settings_discussion",
        to: "/admin/settings/discussion",
        icon: "i-tabler-message-circle",
      },
      {
        name: "settings-media",
        label: "admin.nav.settings_media",
        to: "/admin/settings/media",
        icon: "i-tabler-photo",
      },
      {
        name: "settings-notifications",
        label: "admin.nav.settings_notifications",
        to: "/admin/settings/notifications",
        icon: "i-tabler-bell-ringing",
      },
      {
        name: "settings-oauth",
        label: "admin.nav.settings_oauth",
        to: "/admin/settings/oauth",
        icon: "i-tabler-brand-oauth",
      },
      {
        name: "settings-payment",
        label: "admin.nav.settings_payment",
        to: "/admin/settings/payment",
        icon: "i-tabler-credit-card",
      },
    ],
  },
];

// Filter menu items by capability
const filterItems = (items: MenuItem[]): MenuItem[] =>
  items
    .filter((item) => !item.cap || can(item.cap))
    .map((item) => ({
      ...item,
      children: item.children ? filterItems(item.children) : undefined,
    }))
    .filter((item) => !item.children || item.children.length > 0);

const menu = computed(() => filterItems(ALL_MENU));

const isParentActive = (item: MenuItem) => {
  if (!item.children) return false;
  return item.children.some((child) => route.path.startsWith(child.to || ""));
};
</script>

<style scoped>
aside {
  scrollbar-width: thin;
  scrollbar-color: var(--ui-border) transparent;
}
</style>
