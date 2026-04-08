<script setup lang="ts">
import type { PostCard } from "~/types/models/post";

const props = withDefaults(
  defineProps<{
    posts: PostCard[];
    maxLength?: number;
  }>(),
  { maxLength: 12 },
);

const { defaultCover } = usePostCover();
const items = computed(() => props.posts.slice(0, props.maxLength));
</script>

<template>
  <div class="columns-2 md:columns-3 gap-2 md:gap-4 space-y-0">
    <article
      v-for="post in items"
      :key="post.id"
      class="break-inside-avoid mb-2 md:mb-4">
      <NuxtLink
        :to="`/posts/${post.slug}`"
        class="group block overflow-hidden rounded-sm ring ring-default bg-default hover:shadow-lg transition-all">
        <figure class="relative overflow-hidden">
          <BaseImg
            :src="post.cover?.trim() || defaultCover"
            :alt="post.title"
            class="w-full object-cover group-hover:scale-105 transition-transform duration-300" />
          <div
            class="absolute inset-0 bg-gradient-to-t from-black/50 to-transparent opacity-0 group-hover:opacity-100 transition-opacity" />
          <div class="absolute top-2 left-2">
            <UBadge color="primary" size="xs">{{ post.category }}</UBadge>
          </div>
        </figure>
        <div class="p-3">
          <h3
            class="text-sm font-semibold text-highlighted group-hover:text-primary transition-colors line-clamp-2 mb-1.5">
            {{ post.title }}
          </h3>
          <p v-if="post.excerpt" class="text-xs text-muted line-clamp-2 mb-2">
            {{ post.excerpt }}
          </p>
          <div class="flex items-center gap-3 text-xs text-muted">
            <span class="flex items-center gap-1">
              <UIcon name="i-tabler-eye" class="size-3" />{{ post.views }}
            </span>
            <span class="flex items-center gap-1">
              <UIcon name="i-tabler-message" class="size-3" />{{
                post.comments
              }}
            </span>
            <span class="ml-auto">{{ formatDate(post.date) }}</span>
          </div>
        </div>
      </NuxtLink>
    </article>
  </div>
</template>
