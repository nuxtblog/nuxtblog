<script setup lang="ts">
import type { FollowUserItem } from '~/composables/useFollowApi'

const props = defineProps<{
  userId: number
  type: 'following' | 'followers'
}>()

const { t } = useI18n()
const followApi = useFollowApi()

const page = ref(1)
const size = 20
const rawLoading = ref(true)
const loading = useMinLoading(rawLoading)
const list = ref<FollowUserItem[]>([])
const total = ref(0)

async function load() {
  rawLoading.value = true
  try {
    const res = props.type === 'following'
      ? await followApi.getFollowing(props.userId, page.value, size)
      : await followApi.getFollowers(props.userId, page.value, size)
    list.value = res.list
    total.value = res.total
  } finally {
    rawLoading.value = false
  }
}

watch(page, load)
onMounted(load)

const emptyIcon = computed(() => props.type === 'following' ? 'i-tabler-user-off' : 'i-tabler-users')
const emptyText = computed(() =>
  props.type === 'following' ? t('site.user.no_following_yet') : t('site.user.no_followers_yet')
)
</script>

<template>
  <BaseCardList :loading="loading" :empty="list.length === 0">

    <template #skeleton>
      <div class="flex items-center gap-3">
        <USkeleton class="size-12 rounded-full shrink-0" />
        <div class="flex-1 space-y-2">
          <USkeleton class="h-4 w-32" />
          <USkeleton class="h-3 w-24" />
          <USkeleton class="h-3 w-full" />
        </div>
      </div>
    </template>

    <template #empty>
      <div class="text-center py-12 text-muted">
        <UIcon :name="emptyIcon" class="size-12 mx-auto mb-3" />
        <p>{{ emptyText }}</p>
      </div>
    </template>

    <NuxtLink
      v-for="user in list"
      :key="user.id"
      :to="`/user/${user.id}`"
      class="block group">
      <UCard class="hover:shadow-md transition-shadow">
        <div class="flex items-center gap-3">
          <BaseAvatar
            :src="user.avatar"
            :alt="user.display_name || user.username"
            size="lg"
            class="shrink-0" />
          <div class="flex-1 min-w-0">
            <p class="font-semibold text-highlighted group-hover:text-primary transition-colors">
              {{ user.display_name || user.username }}
            </p>
            <p class="text-xs text-muted">@{{ user.username }}</p>
            <p v-if="user.bio" class="text-sm text-muted line-clamp-1 mt-0.5">{{ user.bio }}</p>
          </div>
          <UIcon name="i-tabler-chevron-right" class="size-4 text-muted shrink-0" />
        </div>
      </UCard>
    </NuxtLink>

    <div v-if="total > size" class="flex justify-center pt-2">
      <UPagination v-model:page="page" :total="total" :items-per-page="size" />
    </div>

  </BaseCardList>
</template>
