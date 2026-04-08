<script setup lang="ts">
import type { WidgetConfig } from '~/composables/useWidgetConfig'
const props = defineProps<{ config: WidgetConfig }>()
const { t } = useI18n()
const maxCount = computed(() => props.config.maxCount ?? 5)
const title = computed(() => props.config.title || t('site.widget.random.default_title'))

const postApi = usePostApi();
const posts = ref<Awaited<ReturnType<typeof postApi.getPosts>>["data"]>([]);
const loading = ref(false);

const fetchRandom = async () => {
  loading.value = true;
  try {
    const res = await postApi.getPosts({
      page: 1,
      page_size: maxCount.value,
      status: "published",
      sort_by: "random",
    });
    posts.value = res.data ?? [];
  } catch {
    posts.value = [];
  } finally {
    loading.value = false;
  }
};

await fetchRandom();
</script>

<template>
  <UCard :ui="{ header: 'p-2.5 sm:px-4' }">
    <template #header>
      <div class="flex items-center justify-between">
        <h3 class="font-semibold text-highlighted">{{ title }}</h3>
        <UButton
          icon="i-tabler-refresh"
          variant="ghost"
          color="neutral"
          size="xs"
          square
          :loading="loading"
          @click="fetchRandom" />
      </div>
    </template>

    <div v-if="posts.length > 0" class="-mx-4 -my-4 divide-y divide-default">
      <a
        v-for="post in posts"
        :key="post.id"
        :href="`/posts/${post.slug}`"
        :title="post.title"
        class="block px-4 py-3 hover:bg-muted transition-colors">
        <div class="text-sm text-default mb-1.5 line-clamp-2 leading-relaxed">
          {{ post.title }}
        </div>
        <div class="flex items-center text-xs text-muted">
          <UIcon name="i-tabler-eye" class="size-3 mr-1" />
          <span>{{ post.view_count ?? 0 }}</span>
          <span class="mx-1">·</span>
          <UIcon name="i-tabler-message" class="size-3 mr-1" />
          <span>{{ post.comment_count ?? 0 }}</span>
        </div>
      </a>
    </div>
    <div v-else class="py-6 text-center text-sm text-muted">{{ $t('site.widget.random.no_posts') }}</div>
  </UCard>
</template>
