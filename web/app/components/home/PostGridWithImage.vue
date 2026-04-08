<template>
  <div class="grid grid-cols-2 md:grid-cols-3 gap-2 md:gap-6">
    <article
      v-for="post in filteredPosts"
      :key="post.id"
      class="group flex flex-col overflow-hidden rounded-sm ring ring-default bg-default hover:shadow-lg transition-all cursor-pointer">
      <NuxtLink :to="`posts/${post.slug}`" class="cursor-pointer">
        <!-- 封面图 -->
        <figure class="relative overflow-hidden aspect-video">
          <BaseImg
            :src="post.cover?.trim() || defaultCover"
            :alt="post.title"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
          <div class="absolute inset-2 sm:inset-3">
            <UBadge color="primary" size="sm">{{ post.category }}</UBadge>
          </div>
        </figure>

        <!-- 内容 -->
        <div class="p-2 md:p-4 flex flex-col flex-1">
          <h3
            class="text-default text-sm md:font-bold mb-2 group-hover:text-primary transition-colors">
            {{ post.title }}
          </h3>
          <p
            v-if="showExcerpt"
            class="text-xs text-muted mb-4 line-clamp-2 hidden md:block">
            {{ post.excerpt }}
          </p>
          <div
            class="flex items-center justify-between text-xs text-muted mt-auto">
            <span>{{ formatDate(post.date) }}</span>
            <div class="flex items-center gap-3">
              <span class="flex items-center gap-1">
                <UIcon name="i-tabler-eye" class="size-3.5" />
                {{ abbreviateNumber(post.views) }}
              </span>
              <span class="flex items-center gap-1">
                <UIcon name="i-tabler-message" class="size-3.5" />
                {{ abbreviateNumber(post.comments) }}
              </span>
            </div>
          </div>
        </div></NuxtLink
      >
    </article>
  </div>
</template>

<script setup lang="ts">
import type { PostCard } from "~/types/models/post";

const props = withDefaults(
  defineProps<{
    posts: PostCard[];
    maxLength?: number;
    showExcerpt?: boolean;
  }>(),
  {
    maxLength: 12,
    showExcerpt: true,
  },
);

const filteredPosts = computed(() => {
  return props.posts.slice(0, props.maxLength ?? props.posts.length);
});

const { defaultCover } = usePostCover();
</script>
