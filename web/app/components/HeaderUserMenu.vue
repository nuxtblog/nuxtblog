<script setup lang="ts">
const { t } = useI18n();
const authStore = useAuthStore();
const followApi = useFollowApi();

const displayName = computed(
  () => authStore.user?.display_name || authStore.user?.username || "",
);

const followStats = ref({ follower_count: 0, following_count: 0 });

onMounted(async () => {
  if (!authStore.isLoggedIn || !authStore.user?.id) return;
  const res = await followApi.getStatus(authStore.user.id).catch(() => null);
  if (res) followStats.value = res;
});

const handleLogout = async () => {
  await authStore.logout();
  await navigateTo("/");
};

const dropdownItems = computed(() => {
  const uid = authStore.user?.id;
  const items: any[][] = [];

  items.push([{ slot: "profile", disabled: true }]);

  if (authStore.user?.role === 3) {
    items.push([
      {
        label: t("site.user_menu.dashboard"),
        icon: "i-tabler-layout-dashboard",
        onSelect: () => navigateTo("/admin/dashboard"),
      },
    ]);
  }

  items.push([
    {
      label: t("site.user_menu.profile"),
      icon: "i-tabler-user",
      onSelect: () => navigateTo(`/user/${uid}`),
    },
    {
      label: t("site.user_menu.my_posts"),
      icon: "i-tabler-file-text",
      onSelect: () => navigateTo(`/user/${uid}?tab=posts`),
    },
    {
      label: t("site.user_menu.favorites"),
      icon: "i-tabler-bookmark",
      onSelect: () => navigateTo(`/user/${uid}?tab=favorites`),
    },
  ]);

  items.push([
    {
      label: t("site.user_menu.settings"),
      icon: "i-tabler-settings",
      onSelect: () => navigateTo("/user/settings"),
    },
  ]);

  items.push([
    {
      label: t("site.user_menu.logout"),
      icon: "i-tabler-logout",
      color: "error" as const,
      onSelect: handleLogout,
    },
  ]);

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
