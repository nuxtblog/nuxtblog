<template>
  <div class="shrink-0 border-b border-default bg-default px-4 md:px-6 py-3 md:py-4">
    <div class="flex items-start justify-between gap-4">
      <div class="min-w-0">
        <h1 class="text-lg md:text-xl font-semibold text-highlighted leading-tight">
          {{ title }}
        </h1>
        <p v-if="subtitle" class="text-sm text-muted mt-0.5">
          {{ subtitle }}
        </p>
      </div>

      <div class="shrink-0 flex items-center gap-2">
        <slot name="actions" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title: string;
  subtitle?: string;
}>();
</script>
