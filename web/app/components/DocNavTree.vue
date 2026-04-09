<script setup lang="ts">
import type { DocItem } from '~/types/api/doc'

const props = defineProps<{
  docs: DocItem[]
  childDocsByParent: Record<number, DocItem[]>
  collectionSlug: string
  currentSlug: string
  depth?: number
}>()

const depth = computed(() => props.depth ?? 0)
</script>

<template>
  <ul :class="depth === 0 ? 'space-y-0.5' : 'ml-4 mt-0.5 space-y-0.5 border-l border-default pl-2'">
    <li v-for="item in docs" :key="item.id">
      <NuxtLink
        :to="`/docs/${collectionSlug}/${item.slug}`"
        class="flex items-center gap-1.5 px-2 py-1.5 rounded-md text-sm transition-colors hover:bg-elevated"
        :class="item.slug === currentSlug
          ? (depth === 0 ? 'bg-primary/10 text-primary font-medium' : 'text-primary font-medium')
          : 'text-muted'">
        <UIcon v-if="depth === 0" name="i-tabler-file-text" class="size-3.5 shrink-0" />
        <span class="truncate">{{ item.title }}</span>
      </NuxtLink>
      <DocNavTree
        v-if="childDocsByParent[item.id]?.length"
        :docs="childDocsByParent[item.id]!"
        :child-docs-by-parent="childDocsByParent"
        :collection-slug="collectionSlug"
        :current-slug="currentSlug"
        :depth="depth + 1"
      />
    </li>
  </ul>
</template>
