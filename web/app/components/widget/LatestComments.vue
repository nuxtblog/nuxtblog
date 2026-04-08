<script setup lang="ts">
import type { WidgetConfig } from "~/composables/useWidgetConfig";
const props = defineProps<{ config: WidgetConfig }>();
const { t } = useI18n()
const maxCount = computed(() => props.config.maxCount ?? 5);
const title = computed(() => props.config.title || t('site.widget.latest_comments.default_title'));

const commentApi = useCommentApi();
const { data } = await useAsyncData("widget-comments", () =>
  commentApi
    .getComments({
      status: "approved",
      page: 1,
      page_size: maxCount.value,
      order: "desc",
    })
    .catch(() => ({ list: [], total: 0, page: 1, page_size: maxCount.value })),
);
const comments = computed(() => data.value?.list ?? []);

const formatTime = (dateStr: string) => {
  const d = new Date(dateStr);
  const diff = Date.now() - d.getTime();
  const m = Math.floor(diff / 60000);
  if (m < 1) return t('site.activity.time_just_now');
  if (m < 60) return t('site.activity.time_minutes_ago', { n: m });
  const h = Math.floor(m / 60);
  if (h < 24) return t('site.activity.time_hours_ago', { n: h });
  const day = Math.floor(h / 24);
  return t('site.activity.time_days_ago', { n: day });
};
</script>

<template>
  <UCard
    :ui="{ header: 'p-2.5 sm:px-4', footer: 'p-2.5 sm:px-4 flex justify-end' }">
    <template #header>
      <h3 class="font-semibold text-highlighted">{{ title }}</h3>
    </template>

    <div v-if="comments.length > 0" class="-mx-4 -my-4 divide-y divide-default">
      <a
        v-for="comment in comments"
        :key="comment.comment_id"
        :href="`/posts/${comment.post_id}`"
        class="block px-4 py-3 hover:bg-muted transition-colors">
        <div class="flex gap-3">
          <BaseAvatar
            :src="comment.author?.avatar || ''"
            :alt="comment.author?.name || $t('site.widget.latest_comments.anonymous')"
            size="sm"
            class="shrink-0" />
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span class="text-sm font-medium text-highlighted">{{
                comment.author?.name || $t('site.widget.latest_comments.anonymous')
              }}</span>
              <span class="text-xs text-muted">{{
                formatTime(comment.created_at)
              }}</span>
            </div>
            <p class="text-sm text-muted line-clamp-2">{{ comment.content }}</p>
          </div>
        </div>
      </a>
    </div>
    <div v-else class="py-6 text-center text-sm text-muted">{{ $t('site.widget.latest_comments.no_comments') }}</div>

    <template #footer>
      <UButton
        variant="link"
        color="primary"
        trailing-icon="i-tabler-arrow-right"
        size="sm">
        {{ $t('site.widget.latest_comments.view_all') }}
      </UButton>
    </template>
  </UCard>
</template>
