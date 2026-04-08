<template>
  <div class="space-y-3">
    <!-- Loading skeletons -->
    <template v-if="loading">
      <UCard v-for="i in skeletonCount" :key="i">
        <slot name="skeleton">
          <div class="flex items-center gap-3">
            <USkeleton class="size-10 rounded-full shrink-0" />
            <div class="flex-1 space-y-2">
              <USkeleton class="h-3.5 w-32" />
              <USkeleton class="h-3 w-full" />
              <USkeleton class="h-3 w-2/3" />
            </div>
          </div>
        </slot>
      </UCard>
    </template>

    <!-- Empty state -->
    <UCard v-else-if="empty">
      <slot name="empty">
        <div class="text-center py-12 text-muted">
          <UIcon name="i-tabler-inbox" class="size-12 mx-auto mb-3" />
          <p>{{ $t('common.no_data') }}</p>
        </div>
      </slot>
    </UCard>

    <!-- Items -->
    <template v-else>
      <slot />
    </template>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  loading: boolean
  empty: boolean
  skeletonCount?: number
}>(), {
  skeletonCount: 4,
})
</script>
