<template>
  <UCard class="hover:shadow-md transition-shadow">
    <!-- 图标 + 趋势 -->
    <div class="flex items-center justify-between mb-4">
      <div
        :class="iconBg"
        class="h-12 w-12 rounded-md flex items-center justify-center"
      >
        <UIcon :name="icon" :class="iconColor" class="size-6" />
      </div>
      <Trend :value="trend" />
    </div>

    <!-- 主数字 -->
    <div class="text-2xl font-bold text-highlighted mb-1">{{ total }}</div>

    <!-- 描述文字 -->
    <div class="text-sm text-muted">{{ label }}</div>

    <!-- 辅助信息 -->
    <div v-if="$slots.description || description" class="mt-3 text-xs text-muted">
      <slot name="description">{{ description }}</slot>
    </div>
  </UCard>
</template>

<script setup lang="ts">
import Trend from "./Trend.vue";

const props = defineProps({
  icon: { type: String, required: true },
  iconColor: { type: String, default: "text-primary" },
  iconBg: { type: String, default: "bg-primary/10" },
  trend: { type: Number, default: 0 },
  total: { type: [Number, String], required: true },
  label: { type: String, required: true },
  description: { type: String, default: "" },
});
</script>
