<template>
  <div :class="[containerClass, 'mx-auto px-4 md:px-6 py-4']">
    <slot />
  </div>
</template>

<script setup lang="ts">
const { containerClass } = useContainerWidth()
</script>
