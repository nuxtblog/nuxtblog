<script setup lang="ts">
defineProps<{
  localValue: string;
}>();

const emit = defineEmits(["input"]);
</script>

<template>
  <div class="relative">
    <UInput
      :model-value="localValue"
      @input="$emit('input', $event)"
      :placeholder="$t('admin.media.search_placeholder')"
      size="sm"
      class="w-48 sm:w-full"
      :ui="{ base: 'pl-8' }">
      <template #leading>
        <UIcon name="i-tabler-search" class="size-4 text-muted" />
      </template>
    </UInput>
  </div>
</template>
