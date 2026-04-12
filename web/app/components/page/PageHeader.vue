<template>
  <div :class="[containerClass, 'mx-auto px-4 md:px-6']">
    <div class="py-6 md:py-8 border-b border-default">
        <!-- Optional breadcrumb -->
        <div v-if="$slots.breadcrumb" class="mb-3">
          <slot name="breadcrumb" />
        </div>

        <!-- Top row: icon + title + actions -->
        <div class="flex flex-wrap items-start justify-between gap-4">
          <div class="flex items-center gap-3 min-w-0">
            <slot name="icon">
              <div v-if="icon" class="size-10 rounded-md bg-primary/10 flex items-center justify-center shrink-0">
                <UIcon :name="icon" class="size-5 text-primary" />
              </div>
            </slot>
            <div class="min-w-0">
              <h1 class="text-xl md:text-2xl font-bold text-highlighted truncate">{{ title }}</h1>
              <p v-if="subtitle" class="text-sm text-muted mt-0.5">{{ subtitle }}</p>
            </div>
          </div>
          <div v-if="$slots.actions" class="flex items-center gap-2 shrink-0">
            <slot name="actions" />
          </div>
        </div>

        <!-- Optional stats row -->
        <div v-if="stats?.length" class="flex items-center gap-3 mt-4 flex-wrap">
          <div
            v-for="stat in stats"
            :key="stat.label"
            class="flex items-center gap-1.5 text-sm">
            <span class="font-semibold" :class="stat.colorClass || 'text-highlighted'">{{ stat.value }}</span>
            <span class="text-muted">{{ stat.label }}</span>
          </div>
        </div>

        <!-- Toolbar slot (search bar, filters, tabs, etc.) -->
        <div v-if="$slots.toolbar" class="mt-4">
          <slot name="toolbar" />
        </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface PageHeaderStat {
  label: string
  value: string | number
  /** Full Tailwind class string, e.g. 'text-primary'. Avoids dynamic class concatenation. */
  colorClass?: string
}

defineProps<{
  title: string
  subtitle?: string
  icon?: string
  stats?: PageHeaderStat[]
}>()

const { containerClass } = useContainerWidth()
</script>
