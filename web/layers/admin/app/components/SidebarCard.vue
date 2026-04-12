<template>
  <div class="group/card rounded-lg border border-default">
    <button
      type="button"
      class="relative w-full flex items-center gap-2 px-3 py-2.5 hover:bg-elevated/50 transition-colors"
      @click="$emit('toggle')">
      <UIcon
        name="i-tabler-grip-vertical"
        class="sidebar-drag-handle size-4 cursor-grab text-muted hover:text-default" />
      <span class="flex-1 text-left text-sm font-semibold text-highlighted">{{
        title
      }}</span>
      <slot name="header-actions" />
      <UIcon
        name="i-tabler-chevron-down"
        class="size-4 text-muted transition-transform duration-200 shrink-0"
        :class="{ '-rotate-90': collapsed }" />
    </button>
    <div v-show="!collapsed" class="px-3 pb-3 space-y-3">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string;
  collapsed: boolean;
}>();

defineEmits<{
  toggle: [];
}>();
</script>
