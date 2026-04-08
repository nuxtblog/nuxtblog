<script setup lang="ts">
defineProps<{ tags: { term_taxonomy_id: number; name: string; slug: string }[] }>();
</script>
<template>
  <div v-if="tags.length" class="flex flex-wrap gap-2 pt-8 border-t border-default">
    <span class="text-sm text-muted mr-1 flex items-center gap-1">
      <UIcon name="i-tabler-tags" class="size-4" />{{ $t('site.post.tags_label') }}
    </span>
    <NuxtLink
      v-for="tag in tags"
      :key="tag.term_taxonomy_id"
      :to="`/posts?tag=${tag.term_taxonomy_id}`"
      class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-sm bg-muted hover:bg-primary/10 hover:text-primary transition-colors">
      # {{ tag.name }}
    </NuxtLink>
  </div>
</template>
