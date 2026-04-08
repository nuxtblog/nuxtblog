<template>
  <section>
    <UCarousel
      :items="posts"
      loop
      dots
      :autoplay="autoplayConfig"
      class="w-full rounded-md"
      :ui="{ item: 'w-full ' }">
      <template #default="{ item }">
        <div class="relative h-64 md:h-96 w-full">
          <NuxtLink :to="`posts/${item.slug}`" class="cursor-pointer">
            <BaseImg
              :src="item.cover || defaultCover"
              :alt="item.title"
              class="w-full h-full object-cover" />
            <div
              class="absolute inset-0 bg-linear-to-t from-black/60 to-transparent flex items-end p-6">
              <h2 class="text-white text-base md:text-2xl font-bold">
                {{ item.title }}
              </h2>
            </div></NuxtLink
          >
        </div>
      </template>
    </UCarousel>
  </section>
</template>

<script setup lang="ts">
const props = defineProps<{
  posts: Array<{
    id: number;
    title: string;
    cover: string;
    slug: string;
    [key: string]: unknown;
  }>;
  autoplay?: boolean;
  interval?: number;
}>();

const autoplayConfig = computed(() => {
  if (!props.autoplay) return false;

  return {
    delay: props.interval ?? 5000,
  };
});

const { defaultCover } = usePostCover();
</script>
