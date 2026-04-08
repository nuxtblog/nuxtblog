<template>
  <UButton
    color="neutral"
    variant="outline"
    class="flex flex-col items-center gap-3 h-auto py-4 px-3 group hover:border-primary transition-all"
    @click="$emit('action')"
  >
    <div
      class="size-10 rounded-md bg-primary/10 flex items-center justify-center group-hover:bg-primary group-hover:text-white transition-colors text-primary"
    >
      <slot name="icon" />
    </div>
    <span class="text-xs font-medium">{{ label }}</span>
  </UButton>
</template>

<script setup lang="ts">
defineProps({
  label: { type: String, required: true },
});

defineEmits(["action"]);
</script>
