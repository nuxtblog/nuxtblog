<template>
  <div class="min-h-screen pb-16">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-8']">

      <!-- Title -->
      <div class="flex items-center gap-2 mb-6">
        <UIcon name="i-tabler-users" class="size-6 text-primary shrink-0" />
        <h1 class="text-xl font-bold text-highlighted">{{ $t('site.followers.title') }}</h1>
        <span class="text-sm text-muted ml-1">{{ $t('site.followers.subtitle', { n: total }) }}</span>
      </div>

      <!-- Search and filter -->
      <div class="flex flex-col sm:flex-row gap-3 mb-6">
        <UInput
          v-model="searchQuery"
          :placeholder="$t('site.followers.search_placeholder')"
          leading-icon="i-tabler-search"
          class="flex-1"
          @input="currentPage = 1" />
        <USelect v-model="filterBy" :items="filterOptions" value-key="value" label-key="label" />
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
          <p class="font-semibold text-highlighted">{{ $t('site.followers.login_required') }}</p>
          <NuxtLink to="/login"><UButton color="primary" size="sm">{{ $t('common.login') }}</UButton></NuxtLink>
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="filteredList.length === 0" class="rounded-md bg-default ring-1 ring-default shadow-sm overflow-hidden">
        <div class="py-20 flex flex-col items-center gap-3 text-center px-4">
          <div class="size-16 rounded-full bg-muted flex items-center justify-center">
            <UIcon name="i-tabler-user-x" class="size-8 text-muted" />
          </div>
          <p class="font-semibold text-highlighted">
            {{ searchQuery ? $t('site.followers.no_results') : $t('site.followers.no_followers') }}
          </p>
          <p class="text-sm text-muted">{{ searchQuery ? $t('site.followers.try_keywords') : $t('site.followers.attract_more') }}</p>
        </div>
      </div>

      <!-- List -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div
          v-for="user in pagedList"
          :key="user.id"
          class="rounded-md ring-1 ring-default bg-default shadow-sm hover:shadow-md transition-shadow p-5">
          <div class="flex items-start gap-4">
            <NuxtLink :to="`/user/${user.id}`" class="shrink-0 relative">
              <BaseAvatar
                :src="user.avatar"
                :alt="user.display_name || user.username"
                size="lg"
                class="ring-2 ring-default hover:ring-primary transition-all" />
              <div
                v-if="user.is_following_back"
                class="absolute -bottom-1 -right-1 w-5 h-5 bg-primary rounded-full flex items-center justify-center border-2 border-default"
                :title="$t('site.followers.mutual')">
                <UIcon name="i-tabler-heart" class="size-3 text-white" />
              </div>
            </NuxtLink>

            <div class="flex-1 min-w-0">
              <NuxtLink :to="`/user/${user.id}`">
                <h3 class="text-base font-semibold text-highlighted hover:text-primary transition-colors mb-0.5 truncate">
                  {{ user.display_name || user.username }}
                </h3>
              </NuxtLink>
              <p class="text-xs text-muted mb-2 line-clamp-2">{{ user.bio || $t('site.followers.no_bio') }}</p>

              <div class="flex items-center gap-3 text-xs text-muted mb-3">
                <span class="flex items-center gap-1">
                  <UIcon name="i-tabler-file-text" class="size-3" />
                  {{ user.article_count }} {{ $t('site.followers.articles') }}
                </span>
                <span class="flex items-center gap-1">
                  <UIcon name="i-tabler-users" class="size-3" />
                  {{ user.follower_count }} {{ $t('site.followers.followers') }}
                </span>
              </div>

              <div class="flex items-center gap-2">
                <UButton
                  v-if="!user.is_following_back"
                  color="primary" size="sm" class="flex-1"
                  :loading="actionLoading[user.id]"
                  @click="followBack(user)">{{ $t('site.followers.follow_back') }}</UButton>
                <UButton
                  v-else
                  variant="ghost" color="neutral" size="sm" class="flex-1"
                  :loading="actionLoading[user.id]"
                  @click="unfollowUser(user)">{{ $t('site.followers.following') }}</UButton>
                <NuxtLink :to="`/user/${user.id}`">
                  <UButton variant="outline" color="neutral" size="sm">{{ $t('site.followers.home_page') }}</UButton>
                </NuxtLink>
                <UButton
                  variant="ghost" color="error" size="sm" square
                  :loading="removeLoading[user.id]"
                  :title="$t('site.followers.remove')"
                  @click="removeFollower(user)">
                  <UIcon name="i-tabler-user-x" class="size-4" />
                </UButton>
              </div>

              <p class="text-xs text-muted mt-2">
                {{ $t('site.followers.followed_at', { date: formatDate(user.followed_at) }) }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div v-if="totalPages > 1" class="flex justify-center mt-6">
        <UPagination v-model:page="currentPage" :total="filteredList.length" :items-per-page="pageSize" />
      </div>

      <div v-if="notFollowingBackCount > 0 && !displayLoading" class="mt-4 flex justify-end">
        <UButton color="primary" size="sm" :loading="batchLoading" @click="followAllBack">
          {{ $t('site.followers.follow_all', { n: notFollowingBackCount }) }}
        </UButton>
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

useHead({ title: t('site.followers.title') })

const rawLoading = ref(true)
const displayLoading = useMinLoading(rawLoading)
const searchQuery = ref('')
const filterBy = ref<'all' | 'following' | 'notFollowing'>('all')
const currentPage = ref(1)
const pageSize = 20
const total = ref(0)
const allFollowers = ref<FollowUserItem[]>([])
const actionLoading = ref<Record<number, boolean>>({})
const removeLoading = ref<Record<number, boolean>>({})
const batchLoading = ref(false)

const filterOptions = computed(() => [
  { label: t('site.followers.filter_all'), value: 'all' },
  { label: t('site.followers.filter_following'), value: 'following' },
  { label: t('site.followers.filter_not_following'), value: 'notFollowing' },
])

const filteredList = computed(() => {
  let list = allFollowers.value
  if (filterBy.value === 'following') list = list.filter(u => u.is_following_back)
  else if (filterBy.value === 'notFollowing') list = list.filter(u => !u.is_following_back)
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(u =>
      (u.display_name || u.username).toLowerCase().includes(q) || u.bio.toLowerCase().includes(q)
    )
  }
  return list
})

const pagedList = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredList.value.slice(start, start + pageSize)
})

const totalPages = computed(() => Math.ceil(filteredList.value.length / pageSize))
const notFollowingBackCount = computed(() => allFollowers.value.filter(u => !u.is_following_back).length)

const loadAll = async () => {
  if (!authStore.user?.id) { rawLoading.value = false; return }
  rawLoading.value = true
  try {
    const PAGE = 50
    let page = 1
    let fetched: FollowUserItem[] = []
    let tot = 0
    do {
      const res = await followApi.getFollowers(authStore.user.id, page, PAGE)
      fetched = fetched.concat(res.list)
      tot = res.total
      page++
    } while (fetched.length < tot)
    allFollowers.value = fetched
    total.value = tot
  } catch {
    toast.add({ title: t('common.load_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    rawLoading.value = false
  }
}

onMounted(loadAll)
watch(searchQuery, () => { currentPage.value = 1 })
watch(filterBy, () => { currentPage.value = 1 })

const followBack = async (user: FollowUserItem) => {
  actionLoading.value[user.id] = true
  try {
    await followApi.follow(user.id)
    user.is_following_back = true
  } catch {
    toast.add({ title: t('common.operation_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    actionLoading.value[user.id] = false
  }
}

const unfollowUser = async (user: FollowUserItem) => {
  actionLoading.value[user.id] = true
  try {
    await followApi.unfollow(user.id)
    user.is_following_back = false
  } catch {
    toast.add({ title: t('common.operation_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    actionLoading.value[user.id] = false
  }
}

const removeFollower = async (user: FollowUserItem) => {
  removeLoading.value[user.id] = true
  try {
    await followApi.removeFollower(user.id)
    allFollowers.value = allFollowers.value.filter(u => u.id !== user.id)
    total.value = Math.max(0, total.value - 1)
    toast.add({ title: t('site.followers.removed'), color: 'neutral', icon: 'i-tabler-check' })
  } catch {
    toast.add({ title: t('common.operation_failed'), color: 'error', icon: 'i-tabler-alert-circle' })
  } finally {
    removeLoading.value[user.id] = false
  }
}

const followAllBack = async () => {
  batchLoading.value = true
  try {
    const targets = allFollowers.value.filter(u => !u.is_following_back)
    await Promise.all(targets.map(u =>
      followApi.follow(u.id).then(() => { u.is_following_back = true }).catch(() => {})
    ))
    toast.add({ title: t('site.followers.followed_all'), color: 'success', icon: 'i-tabler-check' })
  } finally {
    batchLoading.value = false
  }
}

const formatDate = (iso: string) => {
  if (!iso) return ''
  return new Date(iso).toLocaleDateString(undefined, { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>
