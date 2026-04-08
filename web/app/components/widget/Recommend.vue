<script setup lang="ts">
import type { WidgetConfig } from '~/composables/useWidgetConfig'
const props = defineProps<{ config: WidgetConfig }>()
const { t } = useI18n()
const { abbreviate } = useNumberFormat()
const maxCount = computed(() => props.config.maxCount ?? 5)
const title = computed(() => props.config.title || t('site.widget.recommend.default_title'))

const postApi = usePostApi();
const { data } = await useAsyncData("widget-recommend", () =>
  postApi
    .getPosts({
      page: 1,
      page_size: maxCount.value,
      status: "published",
      sort_by: "view_count",
    })
    .catch(() => null),
);
const posts = computed(() => data.value?.data ?? []);

</script>

<template>
  <UCard :ui="{ header: 'p-2.5 sm:px-4' }">
    <template #header>
      <h3 class="font-semibold text-highlighted">{{ title }}</h3>
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
          <span>{{ $t('site.widget.recommend.reads', { n: abbreviate(post.view_count ?? 0) }) }}</span>
          <span class="mx-1">·</span>
          <span>{{ $t('site.widget.recommend.comments', { n: post.comment_count ?? 0 }) }}</span>
        </div>
      </a>
    </div>
    <div v-else class="py-6 text-center text-sm text-muted">{{ $t('site.widget.recommend.no_content') }}</div>
  </UCard>
</template>
