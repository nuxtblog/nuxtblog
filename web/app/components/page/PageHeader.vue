<template>
  <section class="relative left-1/2 right-1/2 -ml-[50vw] -mr-[50vw] w-screen bg-gradient-to-r from-primary/10 via-primary/5 to-transparent shadow-sm">
    <div :class="[containerClass, 'mx-auto px-4 md:px-6']">
        <!-- Optional breadcrumb -->
        <div v-if="$slots.breadcrumb" class="pt-4">
          <slot name="breadcrumb" />
        </div>

        <!-- Top row: icon + title + actions -->
        <div class="flex flex-wrap items-center justify-between gap-4 py-6">
          <div class="flex items-center gap-3 min-w-0">
            <slot name="icon">
              <div v-if="icon" class="size-12 rounded-xl bg-gradient-to-br from-primary to-primary/70 flex items-center justify-center shrink-0 shadow-md">
                <UIcon :name="icon" class="size-6 text-white" />
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
        <div v-if="stats?.length" class="flex items-center gap-3 pb-4 flex-wrap">
          <div
            v-for="stat in stats"
            :key="stat.label"
            class="flex items-center gap-1.5 text-sm">
            <span class="font-semibold" :class="stat.colorClass || 'text-highlighted'">{{ stat.value }}</span>
            <span class="text-muted">{{ stat.label }}</span>
          </div>
        </div>

        <!-- Toolbar slot (search bar, filters, tabs, etc.) -->
        <div v-if="$slots.toolbar" class="pb-4">
          <slot name="toolbar" />
        </div>
        <!-- Bottom decorative border -->
        <div class="h-1 bg-gradient-to-r from-primary via-primary/50 to-transparent rounded-full" />
    </div>
  </section>
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
