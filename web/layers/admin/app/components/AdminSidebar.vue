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
            <span class="flex-1 truncate">{{ tl(item.label) }}</span>
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
            :text="tl(item.label)"
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
            <span class="flex-1 truncate">{{ tl(item.label) }}</span>
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
              {{ tl(child.label) }}
            </NuxtLink>
          </div>
        </details>

        <!-- 有子菜单 - 折叠状态 -->
        <UDropdownMenu
          v-if="item.children && desktopCollapsed"
          :items="item.children.map((c) => ({ label: tl(c.label), to: c.to }))"
          :popper="{ placement: 'right-start' }">
          <UTooltip :text="tl(item.label)" :popper="{ placement: 'right' }">
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
import { ADMIN_MENU, type MenuItem } from "../config/admin-menu";

const route = useRoute();
const { t, te } = useI18n();
/** Translate label if i18n key exists, otherwise return as-is (for plugin items). */
const tl = (label: string) => te(label) ? t(label) : label;
const { desktopCollapsed, closeMobile, toggleDesktop } = useAdminSidebar();
const { can } = usePermissions();

const mediaCount = ref(0);
const pendingComments = ref(0);

// Badge is runtime state — kept in component, keyed by menu item name
const badgeMap = computed<Record<string, () => number>>(() => ({
  media: () => Math.min(mediaCount.value, 99),
  comments: () => pendingComments.value,
}));

// Attach badges to static menu definitions at runtime
const ALL_MENU = computed<MenuItem[]>(() =>
  ADMIN_MENU.map((item) => {
    const badgeFn = badgeMap.value[item.name];
    return badgeFn ? { ...item, badge: badgeFn } : item;
  }),
);

// Filter menu items by capability
const filterItems = (items: MenuItem[]): MenuItem[] =>
  items
    .filter((item) => !item.cap || can(item.cap))
    .map((item) => ({
      ...item,
      children: item.children ? filterItems(item.children) : undefined,
    }))
    .filter((item) => !item.children || item.children.length > 0);

const contributionsStore = usePluginContributionsStore();

const navToMenuItem = (nav: { pluginId: string; route: string; title: string; icon?: string; order: number }): MenuItem => ({
  name: `plugin-${nav.pluginId}-${nav.route}`,
  label: nav.title,
  to: nav.route,
  icon: nav.icon,
  order: nav.order,
});

const menu = computed(() => {
  const base = filterItems(ALL_MENU.value);
  const menuNames = new Set(base.map((i) => i.name));

  // 1. Plugin top-level items: no parent, or parent doesn't match any menu
  //    Orphaned items (parent set but unmatched) get pushed to the bottom (order 9999)
  const allPluginNav = contributionsStore.navigation;
  const topLevelPluginItems = allPluginNav
    .filter((n) => !n.parent || !menuNames.has(n.parent))
    .map((n) => {
      const item = navToMenuItem(n);
      if (n.parent && !menuNames.has(n.parent)) {
        item.order = Math.max(n.order, 9999);
      }
      return item;
    });

  const allTopLevel = [...base, ...topLevelPluginItems]
    .sort((a, b) => (a.order ?? 999) - (b.order ?? 999));

  // 2. For each menu item, merge plugin-contributed children
  return allTopLevel.map((item) => {
    const pluginChildren = allPluginNav
      .filter((n) => n.parent === item.name)
      .sort((a, b) => a.order - b.order);
    if (pluginChildren.length === 0) return item;

    const extraChildren = pluginChildren.map(navToMenuItem);

    // Flat menu item auto-upgrades to parent: original `to` becomes first child
    const existingChildren = item.children || [
      { name: `${item.name}-index`, label: item.label, to: item.to!, icon: item.icon, order: 0 },
    ];

    return {
      ...item,
      to: undefined,
      children: [...existingChildren, ...extraChildren].sort(
        (a, b) => (a.order ?? 999) - (b.order ?? 999),
      ),
    };
  });
});

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
