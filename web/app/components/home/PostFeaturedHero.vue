<script setup lang="ts">
import type { PostCard } from "~/types/models/post";

const props = withDefaults(
  defineProps<{
    posts: PostCard[];
    maxLength?: number;
  }>(),
  { maxLength: 5 },
);

const { defaultCover } = usePostCover();
const items = computed(() => props.posts.slice(0, props.maxLength));
const featured = computed(() => items.value[0]);
const rest = computed(() => items.value.slice(1));
</script>

<template>
  <div v-if="featured" class="grid grid-cols-1 md:grid-cols-5 gap-3 md:gap-4">
    <!-- Featured large post -->
    <NuxtLink
      :to="`/posts/${featured.slug}`"
      class="md:col-span-3 group relative overflow-hidden rounded-md ring ring-default bg-default hover:shadow-xl transition-all block">
      <figure class="relative aspect-[4/3] md:aspect-auto md:h-full min-h-64 overflow-hidden">
        <BaseImg
          :src="featured.cover?.trim() || defaultCover"
          :alt="featured.title"
          class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-500" />
        <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-transparent" />
        <div class="absolute bottom-0 left-0 right-0 p-4 md:p-6">
          <UBadge color="primary" class="mb-2">{{ featured.category }}</UBadge>
          <h2 class="text-white font-bold text-base md:text-xl line-clamp-2 mb-2 group-hover:text-primary/90 transition-colors">
            {{ featured.title }}
          </h2>
          <p class="text-white/70 text-sm line-clamp-2 mb-3 hidden md:block">{{ featured.excerpt }}</p>
          <div class="flex items-center gap-3 text-xs text-white/60">
            <BaseAvatar :src="featured.authorAvatar" :alt="featured.authorName" size="2xs" />
            <span>{{ featured.authorName }}</span>
            <span>{{ formatDate(featured.date) }}</span>
            <span class="flex items-center gap-1 ml-auto">
              <UIcon name="i-tabler-eye" class="size-3" />{{ featured.views }}
            </span>
          </div>
        </div>
      </figure>
    </NuxtLink>

    <!-- Side list -->
    <div class="md:col-span-2 flex flex-col gap-3">
      <NuxtLink
        v-for="post in rest"
        :key="post.id"
        :to="`/posts/${post.slug}`"
        class="group flex gap-3 overflow-hidden rounded-md ring ring-default bg-default hover:shadow-md transition-all flex-1">
        <figure class="w-24 shrink-0 overflow-hidden rounded-l-lg">
          <BaseImg
            :src="post.cover?.trim() || defaultCover"
            :alt="post.title"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
        </figure>
        <div class="flex flex-col justify-center py-3 pr-3 flex-1 min-w-0">
          <UBadge color="primary" variant="soft" size="xs" class="mb-1 self-start">{{ post.category }}</UBadge>
          <h3 class="text-sm font-semibold text-highlighted group-hover:text-primary transition-colors line-clamp-2">
            {{ post.title }}
          </h3>
          <div class="flex items-center gap-2 mt-1.5 text-xs text-muted">
            <span>{{ formatDate(post.date) }}</span>
            <span class="flex items-center gap-1">
              <UIcon name="i-tabler-eye" class="size-3" />{{ post.views }}
            </span>
          </div>
        </div>
      </NuxtLink>
    </div>
  </div>
</template>
