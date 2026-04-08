<script setup lang="ts">
import type { TimelineItem } from "@nuxt/ui";
import type { PostCard } from "~/types/models/post";

const props = withDefaults(
  defineProps<{
    posts: PostCard[];
    maxLength?: number;
  }>(),
  { maxLength: 10 },
);

const { defaultCover } = usePostCover();

const items = computed<(TimelineItem & PostCard)[]>(
  () => props.posts.slice(0, props.maxLength) as (TimelineItem & PostCard)[],
);

const getMonth = (dateStr: string) => {
  const d = new Date(dateStr);
  return `${d.getFullYear()}.${String(d.getMonth() + 1).padStart(2, "0")}`;
};
const getDay = (dateStr: string) => new Date(dateStr).getDate();
</script>

<template>
  <UTimeline :items="items">
    <template #indicator="{ item }">
      <div class="flex flex-col items-center text-center w-12 shrink-0 -mt-1">
        <div class="text-xs text-muted leading-none mb-0.5">
          {{ getMonth(item.date) }}
        </div>
        <div class="text-xl font-bold text-primary leading-none">
          {{ getDay(item.date) }}
        </div>
      </div>
    </template>

    <template #title="{ item }">
      <NuxtLink
        :to="`/posts/${item.slug}`"
        class="flex gap-3 overflow-hidden rounded-sm ring ring-default bg-default hover:shadow-md hover:ring-primary/30 transition-all group mb-2">
        <figure class="w-20 md:w-24 shrink-0 overflow-hidden rounded-l-sm">
          <BaseImg
            :src="item.cover?.trim() || defaultCover"
            :alt="item.title"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
        </figure>
        <div class="flex flex-col justify-center py-3 pr-3 flex-1 min-w-0">
          <UBadge
            color="primary"
            variant="soft"
            size="xs"
            class="mb-1.5 self-start">
            {{ item.category }}
          </UBadge>
          <h3
            class="text-sm font-semibold text-highlighted group-hover:text-primary transition-colors line-clamp-2">
            {{ item.title }}
          </h3>
          <p class="text-xs text-muted line-clamp-1 mt-1 hidden md:block">
            {{ item.excerpt }}
          </p>
          <div class="flex items-center gap-3 mt-2 text-xs text-muted">
            <span class="flex items-center gap-1">
              <UIcon name="i-tabler-eye" class="size-3" />{{ item.views }}
            </span>
            <span class="flex items-center gap-1">
              <UIcon name="i-tabler-message" class="size-3" />{{
                item.comments
              }}
            </span>
          </div>
        </div>
      </NuxtLink>
    </template>
  </UTimeline>
</template>
