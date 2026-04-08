<script setup lang="ts">
import type { PostCard } from "~/types/models/post";

const props = withDefaults(
  defineProps<{
    posts: PostCard[];
    maxLength?: number;
  }>(),
  { maxLength: 10 },
);

const { defaultCover } = usePostCover();
const items = computed(() => props.posts.slice(0, props.maxLength));

const medalColor = (index: number) => {
  if (index === 0) return "text-yellow-500";
  if (index === 1) return "text-slate-400";
  if (index === 2) return "text-amber-600";
  return "text-muted";
};
const medalBg = (index: number) => {
  if (index === 0)
    return "bg-yellow-50 dark:bg-yellow-500/10 ring-yellow-200 dark:ring-yellow-500/30";
  if (index === 1)
    return "bg-slate-50 dark:bg-slate-500/10 ring-slate-200 dark:ring-slate-500/30";
  if (index === 2)
    return "bg-amber-50 dark:bg-amber-500/10 ring-amber-200 dark:ring-amber-500/30";
  return "bg-muted/50 ring-default";
};
</script>

<template>
  <div class="space-y-2">
    <NuxtLink
      v-for="(post, index) in items"
      :key="post.id"
      :to="`/posts/${post.slug}`"
      class="group flex items-center gap-3 p-3 rounded-md ring ring-default bg-default hover:shadow-md transition-all">
      <!-- Rank badge -->
      <div
        class="shrink-0 size-8 rounded-md ring flex items-center justify-center font-bold text-sm"
        :class="[medalColor(index), medalBg(index)]">
        <UIcon v-if="index < 3" name="i-tabler-crown" class="size-4" />
        <span v-else>{{ index + 1 }}</span>
      </div>

      <!-- Cover -->
      <figure
        class="w-12 h-12 shrink-0 rounded-sm overflow-hidden hidden sm:block">
        <BaseImg
          :src="post.cover?.trim() || defaultCover"
          :alt="post.title"
          class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300" />
      </figure>

      <!-- Content -->
      <div class="flex-1 min-w-0">
        <h3
          class="text-sm font-semibold text-highlighted group-hover:text-primary transition-colors truncate">
          {{ post.title }}
        </h3>
        <div class="flex items-center gap-3 mt-0.5 text-xs text-muted">
          <UBadge color="primary" variant="soft" size="xs">{{
            post.category
          }}</UBadge>
          <span class="flex items-center gap-1">
            <UIcon name="i-tabler-eye" class="size-3" />{{ post.views }}
          </span>
          <span class="flex items-center gap-1">
            <UIcon name="i-tabler-message" class="size-3" />{{ post.comments }}
          </span>
        </div>
      </div>

      <UIcon
        name="i-tabler-chevron-right"
        class="size-4 text-muted group-hover:text-primary shrink-0 group-hover:translate-x-0.5 transition-all" />
    </NuxtLink>
  </div>
</template>
