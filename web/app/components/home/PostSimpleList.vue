<template>
  <div class="rounded-sm ring ring-default bg-default divide-y divide-default">
    <article v-for="post in filteredPosts" :key="post.id" class="">
      <NuxtLink
        :to="`posts/${post.slug}`"
        class="group flex items-center justify-between p-4 hover:bg-muted transition-colors cursor-pointer">
        <div class="flex-1 min-w-0 mr-4">
          <h3
            class="text-sm font-medium text-default group-hover:text-primary transition-colors truncate">
            {{ post.title }}
          </h3>
          <div class="flex items-center gap-3 mt-1 text-xs text-muted">
            <span>{{ formatDate(post.date) }}</span>
            <span class="flex items-center gap-1">
              <UIcon name="i-tabler-eye" class="size-3" />
              {{ post.views }}
            </span>
          </div>
        </div>
        <UIcon
          name="i-tabler-chevron-right"
          class="size-4 text-muted group-hover:text-default shrink-0" />
      </NuxtLink>
    </article>
  </div>
</template>

<script setup lang="ts">
interface Post {
  id: number;
  title: string;
  slug: string;
  date: string;
  views: number;
}

const props = withDefaults(
  defineProps<{
    posts: Post[];
    maxLength?: number;
  }>(),
  {
    maxLength: 12,
  },
);

const filteredPosts = computed(() => {
  return props.posts.slice(0, props.maxLength ?? props.posts.length);
});
</script>
