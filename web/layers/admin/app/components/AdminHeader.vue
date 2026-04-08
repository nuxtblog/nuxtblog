<template>
  <header class="w-full border-b border-default bg-default/95 backdrop-blur supports-backdrop-filter:bg-default/60">
    <div class="flex h-full items-center justify-between px-3 md:px-5 gap-3">

      <!-- 左侧：汉堡菜单（移动端）+ Logo -->
      <div class="flex items-center gap-2 md:gap-3 min-w-0">
        <UButton
          class="lg:hidden shrink-0"
          color="neutral"
          variant="ghost"
          :icon="mobileOpen ? 'i-tabler-x' : 'i-tabler-menu-2'"
          square
          size="sm"
          @click="toggleMobile"
        />
        <NuxtLink to="/admin/dashboard" class="flex items-center gap-2 min-w-0">
          <div class="h-7 w-7 md:h-8 md:w-8 rounded-md bg-primary flex items-center justify-center shrink-0">
            <UIcon name="i-tabler-file-text" class="text-white size-4 md:size-5" />
          </div>
          <h1 class="text-base md:text-lg font-semibold text-highlighted hidden sm:block truncate">
            {{ $t('admin.title') }}
          </h1>
        </NuxtLink>
      </div>

      <!-- 右侧操作区 -->
      <div class="flex items-center gap-1 md:gap-2 shrink-0">
        <ClientOnly>
          <ThemeToggle />
        </ClientOnly>

        <!-- 通知 -->
        <AdminNotificationBox />

        <!-- 用户菜单 -->
        <UDropdownMenu :items="userMenuItems" :popper="{ placement: 'bottom-end' }">
          <UButton color="neutral" variant="ghost" class="flex items-center gap-1.5 px-1.5 md:px-2 h-9">
            <BaseAvatar
              :src="authStore.user?.avatar"
              :alt="displayName"
              size="xs"
              class="shrink-0"
            />
            <div class="hidden md:block text-left leading-tight">
              <p class="text-sm font-medium text-highlighted">{{ displayName }}</p>
              <p class="text-xs text-muted">{{ roleLabel }}</p>
            </div>
            <UIcon name="i-tabler-chevron-down" class="size-3.5 text-muted hidden md:block" />
          </UButton>
        </UDropdownMenu>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
const { mobileOpen, toggleMobile } = useAdminSidebar()
const authStore = useAuthStore()
const router = useRouter()
const { t } = useI18n()

const displayName = computed(() =>
  authStore.user?.display_name || authStore.user?.username || t('admin.roles.admin')
)

const roleLabel = computed(() => {
  switch (authStore.user?.role) {
    case 3: return t('admin.roles.super_admin')
    case 2: return t('admin.roles.editor')
    default: return t('admin.roles.user')
  }
})

const handleLogout = async () => {
  await authStore.logout()
  router.push('/auth/login')
}


const userMenuItems = computed(() => [
  [
    {
      label: displayName.value,
      description: authStore.user?.email ?? '',
      disabled: true,
    },
  ],
  [
    { label: t('admin.header.back_to_site'), icon: 'i-tabler-home', to: '/' },
    { label: t('admin.header.user_management'), icon: 'i-tabler-users', to: '/admin/users' },
    { label: t('admin.header.post_management'), icon: 'i-tabler-file-text', to: '/admin/posts' },
  ],
  [
    {
      label: t('admin.header.logout'),
      icon: 'i-tabler-logout',
      color: 'error' as const,
      onSelect: handleLogout,
    },
  ],
])
</script>
