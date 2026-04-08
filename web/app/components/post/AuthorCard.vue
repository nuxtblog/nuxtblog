<script setup lang="ts">
const { t } = useI18n()
defineProps<{
  author: { nickname: string; avatar?: string; bio?: string };
}>();
</script>
<template>
  <div class="flex gap-4 p-6 rounded-md bg-muted/40 ring ring-default mt-8">
    <BaseAvatar :src="author.avatar || ''" :alt="author.nickname" size="lg" class="shrink-0 ring-2 ring-primary/20" />
    <div class="flex-1 min-w-0">
      <p class="text-xs text-muted mb-0.5">{{ $t('site.post.author') }}</p>
      <h4 class="font-bold text-highlighted text-base">{{ author.nickname }}</h4>
      <p class="text-sm text-muted mt-1 line-clamp-2">{{ author.bio || $t('site.post.no_bio') }}</p>
      <NuxtLink to="/user/profile" class="inline-flex items-center gap-1 text-xs text-primary mt-2 hover:underline">
        {{ $t('site.post.view_profile') }} <UIcon name="i-tabler-arrow-right" class="size-3" />
      </NuxtLink>
    </div>
  </div>
</template>
