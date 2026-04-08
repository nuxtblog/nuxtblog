<template>
  <div class="min-h-screen pb-16">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-8']">

      <!-- Title -->
      <div class="flex items-center gap-2 mb-6">
        <UIcon name="i-tabler-user-plus" class="size-6 text-primary shrink-0" />
        <h1 class="text-xl font-bold text-highlighted">{{ $t('site.following.title') }}</h1>
        <span class="text-sm text-muted ml-1">{{ $t('site.following.subtitle', { n: total }) }}</span>
      </div>

      <!-- Search -->
      <div class="flex flex-col sm:flex-row gap-3 mb-6">
        <UInput
          v-model="searchQuery"
          :placeholder="$t('site.following.search_placeholder')"
          leading-icon="i-tabler-search"
          class="flex-1"
          @input="currentPage = 1" />
      </div>

      <!-- Loading skeleton -->
      <div v-if="displayLoading" class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div v-for="i in 4" :key="i" class="rounded-md ring-1 ring-default bg-default shadow-sm p-5">
          <div class="flex items-start gap-4">
            <USkeleton class="size-12 rounded-full shrink-0" />
            <div class="flex-1 space-y-2">
              <USkeleton class="h-5 w-32" />
              <USkeleton class="h-4 w-full" />
              <div class="flex gap-2 pt-1">
                <USkeleton class="h-8 flex-1 rounded-md" />
                <USkeleton class="h-8 w-16 rounded-md" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Not logged in -->
      <div v-else-if="!authStore.isLoggedIn" class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">
        <div class="py-20 flex flex-col items-center gap-3 text-center px-4">
          <div class="size-16 rounded-full bg-muted flex items-center justify-center">
            <UIcon name="i-tabler-lock" class="size-8 text-muted" />
          </div>
          <p class="font-semibold text-highlighted">{{ $t('site.following.login_required') }}</p>
          <NuxtLink to="/login"><UButton color="primary" size="sm">{{ $t('common.login') }}</UButton></NuxtLink>
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="filteredList.length === 0" class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">
        <div class="py-20 flex flex-col items-center gap-3 text-center px-4">
          <div class="size-16 rounded-full bg-muted flex items-center justify-center">
            <UIcon name="i-tabler-user-search" class="size-8 text-muted" />
          </div>
          <p class="font-semibold text-highlighted">
            {{ searchQuery ? $t('site.following.no_results') : $t('site.following.no_following') }}
          </p>
          <p class="text-sm text-muted">{{ searchQuery ? $t('site.following.try_keywords') : $t('site.following.discover') }}</p>
          <NuxtLink v-if="!searchQuery" to="/authors" class="mt-1">
            <UButton color="primary" size="sm">
              <UIcon name="i-tabler-search" class="size-4 mr-1" />
              {{ $t('site.following.discover_btn') }}
            </UButton>
          </NuxtLink>
        </div>
      </div>

      <!-- List -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div
          v-for="user in pagedList"
          :key="user.id"
          class="rounded-md ring-1 ring-default bg-default shadow-sm hover:shadow-md transition-shadow p-5">
          <div class="flex items-start gap-4">
            <NuxtLink :to="`/user/${user.id}`" class="shrink-0">
              <BaseAvatar
                :src="user.avatar"
                :alt="user.display_name || user.username"
                size="lg"
                class="ring-2 ring-default hover:ring-primary transition-all" />
            </NuxtLink>

            <div class="flex-1 min-w-0">
              <NuxtLink :to="`/user/${user.id}`">
                <h3 class="text-base font-semibold text-highlighted hover:text-primary transition-colors mb-0.5 truncate">
                  {{ user.display_name || user.username }}
                </h3>
              </NuxtLink>
              <p class="text-xs text-muted mb-2 line-clamp-2">{{ user.bio || $t('site.following.no_bio') }}</p>

              <div class="flex items-center gap-3 text-xs text-muted mb-3">
                <span class="flex items-center gap-1">
                  <UIcon name="i-tabler-file-text" class="size-3" />
                  {{ user.article_count }} {{ $t('site.following.articles') }}
                </span>
                <span class="flex items-center gap-1">
                  <UIcon name="i-tabler-users" class="size-3" />
                  {{ user.follower_count }} {{ $t('site.following.followers') }}
                </span>
              </div>

              <div class="flex items-center gap-2">
                <UButton
                  variant="ghost" color="neutral" size="sm" class="flex-1"
                  :loading="actionLoading[user.id]"
                  @click="unfollowUser(user)">{{ $t('site.following.unfollow') }}</UButton>
                <NuxtLink :to="`/user/${user.id}`">
                  <UButton color="primary" size="sm">{{ $t('site.following.home_page') }}</UButton>
                </NuxtLink>
              </div>

              <p class="text-xs text-muted mt-2">
                {{ $t('site.following.followed_at', { date: formatDate(user.followed_at) }) }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div v-if="totalPages > 1" class="flex justify-center mt-6">
        <UPagination v-model:page="currentPage" :total="filteredList.length" :items-per-page="pageSize" />
      </div>

    </div>
  </div>
</template>

<script setup lang="ts">
import type { FollowUserItem } from '~/composables/useFollowApi'

const { containerClass } = useContainerWidth()
const { t } = useI18n()
const authStore = useAuthStore()
const followApi = useFollowApi()
const toast = useToast()

useHead({ title: t('site.following.title') })

const rawLoading = ref(true)
const displayLoading = useMinLoading(rawLoading)
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = 20
const total = ref(0)
const allFollowing = ref<FollowUserItem[]>([])
const actionLoading = ref<Record<number, boolean>>({})

const filteredList = computed(() => {
  if (!searchQuery.value) return allFollowing.value
  const q = searchQuery.value.toLowerCase()
  return allFollowing.value.filter(u =>
    (u.display_name || u.username).toLowerCase().includes(q) || u.bio.toLowerCase().includes(q)
  )
})

const pagedList = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredList.value.slice(start, start + pageSize)
})

const totalPages = computed(() => Math.ceil(filteredList.value.length / pageSize))

const loadAll = async () => {
  if (!authStore.user?.id) { rawLoading.value = false; return }
  rawLoading.value = true
  try {
    const PAGE = 50
    let page = 1
    let fetched: FollowUserItem[] = []
    let tot = 0
    do {
      const res = await followApi.getFollowing(authStore.user.id, page, PAGE)
      fetched = fetched.concat(res.list)
      tot = res.total
      page++
    } while (fetched.length < tot)
    allFollowing.value = fetched
    total.value = tot
  } catch {
    toast.add({ title: t('common.load_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    rawLoading.value = false
  }
}

onMounted(loadAll)
watch(searchQuery, () => { currentPage.value = 1 })

const unfollowUser = async (user: FollowUserItem) => {
  actionLoading.value[user.id] = true
  try {
    await followApi.unfollow(user.id)
    allFollowing.value = allFollowing.value.filter(u => u.id !== user.id)
    total.value = Math.max(0, total.value - 1)
    toast.add({ title: t('site.following.unfollowed'), color: 'neutral', icon: 'i-tabler-check' })
  } catch {
    toast.add({ title: t('common.operation_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    actionLoading.value[user.id] = false
  }
}

const formatDate = (iso: string) => {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>
