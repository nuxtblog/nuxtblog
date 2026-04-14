<template>
  <header class="sticky top-0 z-50 border-b border-default bg-default/80 backdrop-blur">
    <div class="max-w-7xl mx-auto px-4 h-16 flex items-center gap-4">
      <!-- Logo -->
      <NuxtLink to="/" class="text-xl md:text-2xl font-bold text-primary shrink-0">
        {{ siteName }}
      </NuxtLink>

      <!-- 桌面端导航 -->
      <nav class="hidden md:flex items-center gap-1 mx-auto">
        <template v-for="item in rootItems" :key="item.local_id">
          <!-- Has children → dropdown -->
          <UDropdownMenu
            v-if="getChildren(item.local_id).length > 0"
            :items="[getChildren(item.local_id).map(c => ({
              label: c.label,
              to: c.url,
              target: c.target || undefined,
            }))]"
          >
            <UButton
              variant="ghost"
              :color="isActive(item.url) ? 'primary' : 'neutral'"
              size="sm"
              trailing-icon="i-tabler-chevron-down"
            >
              {{ item.label }}
            </UButton>
          </UDropdownMenu>

          <!-- No children → plain link -->
          <UButton
            v-else
            :to="item.url"
            :target="item.target || undefined"
            variant="ghost"
            :color="isActive(item.url) ? 'primary' : 'neutral'"
            size="sm"
          >
            {{ item.label }}
          </UButton>
        </template>
      </nav>

      <!-- 右侧操作区 -->
      <div class="ml-auto flex items-center gap-1">
        <template v-for="action in headerActions" :key="action.local_id">
          <!-- 分隔符 -->
          <div v-if="action.object_type === 'separator'" class="h-5 w-px bg-border-default border border-default mx-0.5" />
          <!-- Built-in actions -->
          <LangSwitcher v-else-if="action.local_id === 'action:lang_switcher'" />
          <ThemeToggle v-else-if="action.local_id === 'action:theme_toggle'" />
          <template v-else-if="action.local_id === 'action:messages' && authStore.isLoggedIn">
            <NuxtLink to="/messages" class="relative">
              <UButton color="neutral" variant="ghost" icon="i-tabler-message" square size="sm" />
              <span v-if="messageUnread > 0"
                class="absolute -top-1 -right-1 h-4 w-4 rounded-full bg-error text-white text-[10px] flex items-center justify-center">
                {{ messageUnread > 9 ? '9+' : messageUnread }}
              </span>
            </NuxtLink>
          </template>
          <NotificationBox v-else-if="action.local_id === 'action:notifications' && authStore.isLoggedIn" />
          <!-- Plugin item: render contribution or execute command -->
          <ClientOnly v-else-if="action.local_id.startsWith('plugin:')">
            <ContributionSlot :name="'public:header-actions'" :filter-id="action.local_id.replace('plugin:', '')" />
          </ClientOnly>
          <!-- Custom action: icon button link -->
          <UTooltip v-else-if="action.object_type === 'custom'" :text="action.label">
            <UButton
              :to="action.url"
              :target="action.target || undefined"
              color="neutral"
              variant="ghost"
              :icon="action.css_classes || 'i-tabler-link'"
              square
              size="sm"
            />
          </UTooltip>
        </template>
        <HeaderUserMenu />

        <!-- 移动端菜单按钮 -->
        <UButton
          :icon="isMobileMenuOpen ? 'i-tabler-x' : 'i-tabler-menu-2'"
          variant="ghost"
          color="neutral"
          class="md:hidden"
          aria-label="Toggle menu"
          @click="isMobileMenuOpen = !isMobileMenuOpen"
        />
      </div>
    </div>

    <!-- 移动端菜单 -->
    <div v-if="isMobileMenuOpen" class="md:hidden border-t border-default">
      <nav class="flex flex-col p-2 gap-1">
        <template v-for="item in allItems" :key="item.local_id">
          <UButton
            :to="item.url"
            :target="item.target || undefined"
            variant="ghost"
            color="neutral"
            class="justify-start"
            :style="item.parent_local_id ? 'padding-left: 2rem' : ''"
            @click="isMobileMenuOpen = false"
          >
            {{ item.label }}
          </UButton>
        </template>
      </nav>
    </div>
  </header>
</template>

<script setup lang="ts">
import type { NavMenuItem } from '~/types/api/navMenu'

const route = useRoute()
const { getOption } = useOption()
const authStore = useAuthStore()
const messageApi = useMessageApi()
const messageUnread = ref(0)

const siteName = computed(() => getOption('site_name'))
const headerActions = computed(() => getOption('header_actions'))

onMounted(async () => {
  if (!authStore.isLoggedIn || !authStore.user?.id) return
  try {
    const res = await messageApi.unreadCount()
    messageUnread.value = res.count
  } catch {}
})

// primary_menu is autoloaded — type inferred from OPTIONS_SCHEMA, no casting needed
const allItems = computed(() => getOption('primary_menu'))
const rootItems = computed(() => allItems.value.filter(i => !i.parent_local_id))

function getChildren(parentId: string): NavMenuItem[] {
  return allItems.value.filter(i => i.parent_local_id === parentId)
}

function isActive(url: string): boolean {
  if (url === '/') return route.path === '/'
  return route.path.startsWith(url)
}

const isMobileMenuOpen = ref(false)
</script>
