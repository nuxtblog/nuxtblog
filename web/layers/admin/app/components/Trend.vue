<template>
  <div v-if="value !== 0" class="flex items-center gap-1 text-xs">
    <UIcon
      :name="value > 0 ? 'i-tabler-arrow-up' : 'i-tabler-arrow-down'"
      :class="value > 0 ? 'text-success' : 'text-error'"
      class="size-3.5"
    />
    <span :class="value > 0 ? 'text-success' : 'text-error'" class="font-medium">
      {{ value > 0 ? `+${value}%` : `${value}%` }}
    </span>
  </div>
</template>

<script setup lang="ts">
defineProps<{ value: number }>();
</script>
