<script setup lang="ts">
import type { WidgetConfig } from '~/composables/useWidgetConfig'

defineProps<{ config: WidgetConfig }>()

const { t } = useI18n()
const authStore = useAuthStore()
const optionsStore = useOptionsStore()
const userApi = useUserApi()
const postApi = usePostApi()


// Load current user's post count (only when logged in)
const postCount = ref(0)
if (authStore.isLoggedIn && authStore.user?.id) {
  useAsyncData(`widget-userbox-posts-${authStore.user.id}`, () =>
    postApi
      .getPosts({ author_id: authStore.user!.id, status: 'published', page: 1, page_size: 1 })
      .then(r => { postCount.value = r.total; return r.total })
      .catch(() => {}),
  )
}

const bio = computed(
  () =>
    authStore.user?.bio ||
    optionsStore.get('default_author_bio', '') ||
    t('site.widget.user_box.default_bio'),
)

const userCoverBg = computed(
  () => (authStore.user as any)?.metas?.cover || optionsStore.get('default_user_bg', ''),
)

const siteName = computed(() => optionsStore.get('site_name', '') || t('site.widget.user_box.default_site_name'))
const siteDesc = computed(
  () => optionsStore.get('site_description', '') || t('site.widget.user_box.default_site_desc'),
)
</script>

<template>
  <!-- ── Logged in ──────────────────────────────────────── -->
  <UCard v-if="authStore.isLoggedIn && authStore.user" :ui="{ body: 'p-0 sm:p-0' }">
    <!-- Banner -->
    <div class="h-16 relative overflow-hidden rounded-t-[calc(var(--ui-radius)*1.5)]">
      <template v-if="userCoverBg">
        <div class="absolute inset-0 bg-cover bg-center" :style="`background-image:url('${userCoverBg}')`" />
        <div class="absolute inset-0 bg-black/25" />
      </template>
      <template v-else>
        <div class="absolute inset-0 bg-gradient-to-br from-primary/30 via-primary/10 to-violet-500/20" />
        <div class="absolute -top-4 -right-4 w-20 h-20 rounded-full bg-primary/15 blur-2xl" />
        <div class="absolute bottom-0 left-1/3 w-14 h-14 rounded-full bg-violet-400/15 blur-xl" />
      </template>
    </div>

    <!-- Avatar + info -->
    <div class="px-4 pb-4 flex flex-col items-center text-center -mt-7">
      <NuxtLink :to="`/user/${authStore.user.id}`">
        <div class="size-14 rounded-full p-0.5 bg-linear-to-br from-primary via-violet-500 to-indigo-400 shadow-md ring-2 ring-default mb-2">
          <BaseAvatar
            :src="authStore.user.avatar || ''"
            :alt="authStore.user.display_name || authStore.user.username"
            class="size-full"
            size="md" />
        </div>
      </NuxtLink>

      <NuxtLink :to="`/user/${authStore.user.id}`" class="font-bold text-highlighted text-sm leading-none hover:text-primary transition-colors">
        {{ authStore.user.display_name || authStore.user.username }}
      </NuxtLink>
      <p class="text-xs text-muted mt-1.5 leading-relaxed line-clamp-2 max-w-[180px]">{{ bio }}</p>
    </div>

    <USeparator />

    <!-- Quick links -->
    <div class="px-4 py-3 space-y-1.5">
      <NuxtLink
        :to="`/user/${authStore.user.id}`"
        class="flex items-center gap-2 px-2 py-1.5 rounded-md text-sm text-default hover:bg-muted transition-colors">
        <UIcon name="i-tabler-user" class="size-4 text-muted shrink-0" />
        <span>{{ $t('site.widget.user_box.my_profile') }}</span>
      </NuxtLink>
      <NuxtLink
        to="/user/profile"
        class="flex items-center gap-2 px-2 py-1.5 rounded-md text-sm text-default hover:bg-muted transition-colors">
        <UIcon name="i-tabler-settings" class="size-4 text-muted shrink-0" />
        <span>{{ $t('site.widget.user_box.my_settings') }}</span>
      </NuxtLink>
      <NuxtLink
        to="/user/profile?tab=favorites"
        class="flex items-center gap-2 px-2 py-1.5 rounded-md text-sm text-default hover:bg-muted transition-colors">
        <UIcon name="i-tabler-bookmark" class="size-4 text-muted shrink-0" />
        <span>{{ $t('site.widget.user_box.my_favorites') }}</span>
      </NuxtLink>
    </div>
  </UCard>

  <!-- ── Not logged in ──────────────────────────────────── -->
  <UCard v-else :ui="{ body: 'p-0 sm:p-0' }">
    <!-- Banner -->
    <div class="h-16 relative overflow-hidden rounded-t-[calc(var(--ui-radius)*1.5)]">
      <template v-if="userCoverBg">
        <div class="absolute inset-0 bg-cover bg-center" :style="`background-image:url('${userCoverBg}')`" />
        <div class="absolute inset-0 bg-black/25" />
      </template>
      <template v-else>
        <div class="absolute inset-0 bg-gradient-to-br from-primary/30 via-primary/10 to-violet-500/20" />
        <div class="absolute -top-4 -right-4 w-20 h-20 rounded-full bg-primary/15 blur-2xl" />
        <div class="absolute bottom-0 left-1/3 w-14 h-14 rounded-full bg-violet-400/15 blur-xl" />
      </template>
    </div>

    <div class="px-4 pb-4 pt-3 flex flex-col items-center text-center">
      <div class="size-12 rounded-full bg-primary/10 flex items-center justify-center mb-2">
        <UIcon name="i-tabler-user" class="size-6 text-primary" />
      </div>
      <p class="text-sm font-semibold text-highlighted mb-1">{{ siteName }}</p>
      <p class="text-xs text-muted leading-relaxed line-clamp-3 max-w-[180px] mb-3">{{ siteDesc }}</p>

      <div class="w-full space-y-2">
        <UButton to="/auth/login" color="primary" block size="sm" icon="i-tabler-login">
          {{ $t('site.widget.user_box.login') }}
        </UButton>
        <UButton to="/auth/register" color="neutral" variant="outline" block size="sm" icon="i-tabler-user-plus">
          {{ $t('site.widget.user_box.register') }}
        </UButton>
      </div>
    </div>
  </UCard>
</template>
