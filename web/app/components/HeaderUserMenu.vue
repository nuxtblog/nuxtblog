<script setup lang="ts">
import type { NavMenuItem } from '~/types/api/navMenu'
import { executePublicCommand } from '~/composables/useNuxtblogPublic'
import { usePluginLocale } from '~/composables/usePluginLocale'

const { t } = useI18n();
const authStore = useAuthStore();
const followApi = useFollowApi();
const contributionsStore = usePluginContributionsStore();
const contextStore = usePluginContextStore();
const { t: pluginT } = usePluginLocale();

const displayName = computed(
  () => authStore.user?.display_name || authStore.user?.username || "",
);

const followStats = ref({ follower_count: 0, following_count: 0 });

onMounted(async () => {
  if (!authStore.isLoggedIn || !authStore.user?.id) return;
  const res = await followApi.getStatus(authStore.user.id).catch(() => null);
  if (res) followStats.value = res;
});

// Keep isLoggedIn synced in plugin context for `when` expressions
watchEffect(() => {
  contextStore.set('isLoggedIn', authStore.isLoggedIn)
})

const { getOption } = useOption()
const disabledViews = computed(() => getOption('disabled_plugin_views'))

const pluginMenuItems = contributionsStore.getMenuItems('public:user-menu')
const visiblePluginItems = computed(() =>
  pluginMenuItems.value.filter(item =>
    (!item.when || contextStore.evaluateWhen(item.when)) &&
    !disabledViews.value.includes(item.command),
  ),
)

const handleLogout = async () => {
  await authStore.logout();
  await navigateTo("/");
};

// Read menu config (array order = display order)
const userMenuItems = computed(() => getOption('user_menu') as NavMenuItem[])

const dropdownItems = computed(() => {
  const uid = authStore.user?.id;
  const items: any[][] = [];

  // Profile header (always first, not configurable)
  items.push([{ slot: "profile", disabled: true }]);

  let currentGroup: any[] = [];

  // Action handlers
  const actionHandlers: Record<string, () => any> = {
    'action:admin_dashboard': () => ({ label: t("site.user_menu.dashboard"), icon: "i-tabler-layout-dashboard", onSelect: () => navigateTo("/admin/dashboard") }),
    'action:user_profile': () => ({ label: t("site.user_menu.profile"), icon: "i-tabler-user", onSelect: () => navigateTo(`/user/${uid}`) }),
    'action:user_posts': () => ({ label: t("site.user_menu.my_posts"), icon: "i-tabler-file-text", onSelect: () => navigateTo(`/user/${uid}?tab=posts`) }),
    'action:user_favorites': () => ({ label: t("site.user_menu.favorites"), icon: "i-tabler-bookmark", onSelect: () => navigateTo(`/user/${uid}?tab=favorites`) }),
    'action:user_settings': () => ({ label: t("site.user_menu.settings"), icon: "i-tabler-settings", onSelect: () => navigateTo("/user/settings") }),
    'action:user_logout': () => ({ label: t("site.user_menu.logout"), icon: "i-tabler-logout", color: "error" as const, onSelect: handleLogout }),
  }

  for (const menuItem of userMenuItems.value) {
    if (menuItem.object_type === 'separator') {
      if (currentGroup.length) {
        items.push(currentGroup)
        currentGroup = []
      }
      continue
    }

    // Skip admin_dashboard for non-admin users
    if (menuItem.local_id === 'action:admin_dashboard' && authStore.user?.role !== 3) continue

    const handler = actionHandlers[menuItem.local_id]
    if (handler) {
      currentGroup.push(handler())
    } else if (menuItem.local_id.startsWith('plugin:')) {
      const pluginId = menuItem.local_id.replace('plugin:', '')
      currentGroup.push({
        label: menuItem.label,
        icon: menuItem.css_classes || 'i-tabler-puzzle',
        onSelect: () => executePublicCommand(pluginId),
      })
    } else if (menuItem.object_type === 'custom') {
      currentGroup.push({
        label: menuItem.label,
        icon: menuItem.css_classes || 'i-tabler-link',
        onSelect: () => navigateTo(menuItem.url, { external: menuItem.target === '_blank', open: menuItem.target === '_blank' ? { target: '_blank' } : undefined }),
      })
    }
  }

  // Plugin-injected menu items
  if (visiblePluginItems.value.length > 0) {
    if (currentGroup.length) {
      items.push(currentGroup)
      currentGroup = []
    }
    items.push(
      visiblePluginItems.value.map(item => ({
        label: pluginT(item) || item.command,
        icon: item.icon,
        onSelect: () => executePublicCommand(item.command),
      })),
    );
  }

  if (currentGroup.length) {
    items.push(currentGroup)
  }

  return items;
});
</script>

<template>
  <template v-if="authStore.isLoggedIn">
    <UDropdownMenu
      :items="dropdownItems"
      :content="{ align: 'end' }"
      :ui="{ content: 'w-64' }">
      <UButton
        variant="ghost"
        color="neutral"
        class="flex items-center gap-2 px-2">
        <BaseAvatar
          :src="authStore.user?.avatar"
          :alt="displayName"
          size="sm" />
        <span class="hidden sm:block text-sm font-medium truncate max-w-28">{{
          displayName
        }}</span>
        <UIcon
          name="i-tabler-chevron-down"
          class="hidden sm:block size-3.5 text-muted shrink-0" />
      </UButton>

      <template #profile>
        <div class="px-3 py-3 space-y-3">
          <div class="flex items-center gap-3">
            <BaseAvatar
              :src="authStore.user?.avatar"
              :alt="displayName"
              size="md" />
            <div class="flex-1 min-w-0">
              <p class="font-semibold text-highlighted text-sm truncate">
                {{ displayName }}
              </p>
              <p class="text-xs text-muted truncate text-left">
                @{{ authStore.user?.username }}
              </p>
            </div>
          </div>

          <div class="flex gap-4 text-xs">
            <NuxtLink
              :to="`/user/${authStore.user?.id}?tab=following`"
              class="text-muted hover:text-highlighted transition-colors">
              <span class="font-semibold text-highlighted">{{
                followStats.following_count
              }}</span>
              {{ $t("site.user_menu.following") }}
            </NuxtLink>
            <NuxtLink
              :to="`/user/${authStore.user?.id}?tab=followers`"
              class="text-muted hover:text-highlighted transition-colors">
              <span class="font-semibold text-highlighted">{{
                followStats.follower_count
              }}</span>
              {{ $t("site.user_menu.followers") }}
            </NuxtLink>
          </div>
          <ContributionSlot name="public:user-menu-cards" />
        </div>
      </template>
    </UDropdownMenu>
  </template>

  <template v-else>
    <div class="flex items-center gap-1">
      <UButton to="/auth/login" variant="ghost" color="neutral" size="sm">
        {{ $t("site.user_menu.login") }}
      </UButton>
      <UButton to="/auth/register" color="primary" size="sm">
        {{ $t("site.user_menu.register") }}
      </UButton>
    </div>
  </template>
</template>
