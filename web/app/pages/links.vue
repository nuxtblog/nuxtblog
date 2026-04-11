<script setup lang="ts">
import type { FriendlinkItem } from "~/composables/useFriendlinkApi";

const { t } = useI18n();
const { containerClass } = useContainerWidth();
const friendlinkApi = useFriendlinkApi();

const rawLoading = ref(true);
const loading = useMinLoading(rawLoading);
const list = ref<FriendlinkItem[]>([]);

onMounted(async () => {
  rawLoading.value = true;
  try {
    const res = await friendlinkApi.getFriendlinks();
    list.value = res.list ?? [];
  } finally {
    rawLoading.value = false;
  }
});
</script>

<template>
  <div class="min-h-screen pb-16">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-8']">
      <!-- Title -->
      <div class="flex items-center gap-2 mb-6">
        <UIcon name="i-tabler-link" class="size-6 text-primary shrink-0" />
        <h1 class="text-xl font-bold text-highlighted">{{ $t("friendlinks.title") }}</h1>
      </div>
      <p class="text-sm text-muted mb-6">{{ $t("friendlinks.subtitle") }}</p>

      <!-- Loading -->
      <div v-if="loading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="i in 6" :key="i" class="rounded-md bg-default ring-1 ring-default shadow-sm p-5">
          <div class="flex items-center gap-3 mb-3">
            <USkeleton class="size-12 rounded-md shrink-0" />
            <div class="flex-1 space-y-2">
              <USkeleton class="h-4 w-24" />
              <USkeleton class="h-3 w-32" />
            </div>
          </div>
          <USkeleton class="h-3 w-full" />
        </div>
      </div>

      <!-- Empty -->
      <div v-else-if="list.length === 0" class="text-center py-20">
        <div class="size-20 rounded-full bg-primary/10 flex items-center justify-center mx-auto mb-4">
          <UIcon name="i-tabler-link" class="size-10 text-primary" />
        </div>
        <p class="font-semibold text-highlighted mb-1">{{ $t("friendlinks.empty") }}</p>
        <p class="text-sm text-muted mb-4">{{ $t("friendlinks.empty_desc") }}</p>
      </div>

      <!-- Card Grid -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <a
          v-for="item in list"
          :key="item.id"
          :href="item.url"
          target="_blank"
          rel="noopener noreferrer"
          class="group rounded-md bg-default ring-1 ring-default shadow-sm p-5 transition-all hover:shadow-md hover:ring-primary/30">
          <div class="flex items-center gap-3 mb-2">
            <img
              v-if="item.logo"
              :src="item.logo"
              :alt="item.name"
              class="size-12 rounded-md object-cover shrink-0 ring-1 ring-default" />
            <div v-else class="size-12 rounded-md bg-primary/10 flex items-center justify-center shrink-0">
              <UIcon name="i-tabler-link" class="size-6 text-primary" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-semibold text-highlighted truncate group-hover:text-primary transition-colors">
                {{ item.name }}
              </p>
              <p class="text-xs text-muted truncate">{{ item.url }}</p>
            </div>
            <UIcon name="i-tabler-external-link" class="size-4 text-muted opacity-0 group-hover:opacity-100 transition-opacity shrink-0" />
          </div>
          <p v-if="item.description" class="text-xs text-muted line-clamp-2 leading-relaxed">
            {{ item.description }}
          </p>
        </a>
      </div>
    </div>
  </div>
</template>
