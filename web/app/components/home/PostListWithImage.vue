<template>
  <div class="space-y-6">
    <article v-for="post in filteredPosts" :key="post.id" class="">
      <NuxtLink
        :to="`posts/${post.slug}`"
        class="group flex overflow-hidden rounded-sm ring ring-default bg-default hover:shadow-md transition-all cursor-pointer">
        <!-- 封面图 -->
        <figure class="w-48 shrink-0 overflow-hidden">
          <BaseImg
            :src="post.cover || defaultCover"
            :alt="post.title"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300" />
        </figure>

        <!-- 内容 -->
        <div class="flex flex-col justify-between p-4 flex-1 min-w-0">
          <div>
            <div class="flex items-center gap-2 mb-2">
              <UBadge color="primary" variant="soft" size="sm">{{
                post.category
              }}</UBadge>
              <span class="text-xs text-muted">{{
                formatDate(post.date)
              }}</span>
            </div>
            <h3
              class="font-semibold text-default group-hover:text-primary transition-colors line-clamp-2">
              {{ post.title }}
            </h3>
            <p class="text-sm text-muted line-clamp-2 mt-1">
              {{ post.excerpt }}
            </p>
          </div>
          <div class="flex items-center gap-4 text-xs text-muted mt-3">
            <span class="flex items-center gap-1">
              <UIcon name="i-tabler-eye" class="size-3.5" />
              {{ post.views }} 次阅读
            </span>
            <span class="flex items-center gap-1">
              <UIcon name="i-tabler-message" class="size-3.5" />
              {{ post.comments }} 条评论
            </span>
          </div>
        </div></NuxtLink
      >
    </article>
  </div>
</template>

<script setup lang="ts">
interface Post {
  id: number;
  title: string;
  excerpt: string;
  slug: string;
  cover: string;
  category: string;
  date: string;
  views: number;
  comments: number;
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

const { defaultCover } = usePostCover();
</script>
