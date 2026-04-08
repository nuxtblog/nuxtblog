<template>
  <div class="space-y-2">
    <div
      v-for="(widget, index) in modelValue"
      :key="widget.id"
      class="rounded-md border border-default overflow-hidden"
    >
      <!-- 行 -->
      <div
        class="flex items-center gap-3 px-3 py-2.5"
        :class="expanded[widget.id] ? 'bg-elevated' : 'bg-default'"
      >
        <UCheckbox v-model="widget.enabled" />
        <span
          class="flex-1 text-sm select-none"
          :class="widget.enabled ? 'text-highlighted' : 'text-muted'"
        >{{ $t(widget.label) }}</span>
        <div class="flex items-center gap-0.5">
          <UButton
            :icon="expanded[widget.id] ? 'i-tabler-chevron-up' : 'i-tabler-settings'"
            variant="ghost"
            :color="expanded[widget.id] ? 'primary' : 'neutral'"
            size="xs"
            square
            @click="emit('toggle', widget.id)"
          />
          <UButton
            icon="i-tabler-arrow-up"
            variant="ghost"
            color="neutral"
            size="xs"
            square
            :disabled="index === 0"
            @click="moveUp(index)"
          />
          <UButton
            icon="i-tabler-arrow-down"
            variant="ghost"
            color="neutral"
            size="xs"
            square
            :disabled="index === modelValue.length - 1"
            @click="moveDown(index)"
          />
        </div>
      </div>

      <!-- 展开的配置面板 -->
      <div
        v-if="expanded[widget.id]"
        class="px-4 py-3 bg-muted/40 border-t border-default space-y-3"
      >
        <div class="flex items-center gap-3">
          <span class="text-sm text-default flex-1">{{ $t('admin.widgets.widget_title') }}</span>
          <UInput v-model="widget.title" :placeholder="$t(widget.label)" size="sm" class="w-40" />
        </div>

        <template v-if="getWidgetDef(widget.id)?.configFields.includes('showRecent')">
          <div class="flex items-center justify-between">
            <span class="text-sm text-default">{{ $t('admin.widgets.show_recent') }}</span>
            <UCheckbox v-model="widget.showRecent" />
          </div>
        </template>
        <template v-if="getWidgetDef(widget.id)?.configFields.includes('showHot')">
          <div class="flex items-center justify-between">
            <span class="text-sm text-default">{{ $t('admin.widgets.show_hot') }}</span>
            <UCheckbox v-model="widget.showHot" />
          </div>
        </template>
        <template v-if="getWidgetDef(widget.id)?.configFields.includes('maxCount')">
          <div class="flex items-center gap-3">
            <span class="text-sm text-default flex-1">{{ $t('admin.widgets.max_count') }}</span>
            <UInput v-model.number="widget.maxCount" type="number" :min="1" :max="20" size="sm" class="w-20" />
            <span class="text-xs text-muted">{{ $t('admin.widgets.items') }}</span>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { WidgetConfig } from '~/composables/useWidgetConfig'
import { getWidgetDef } from '~/config/widgets'

const props = defineProps<{
  modelValue: WidgetConfig[]
  expanded: Record<string, boolean>
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', v: WidgetConfig[]): void
  (e: 'toggle', id: string): void
}>()

const moveUp = (index: number) => {
  if (index === 0) return
  const arr = [...props.modelValue]
  ;[arr[index - 1], arr[index]] = [arr[index], arr[index - 1]]
  emit('update:modelValue', arr)
}

const moveDown = (index: number) => {
  if (index === props.modelValue.length - 1) return
  const arr = [...props.modelValue]
  ;[arr[index], arr[index + 1]] = [arr[index + 1], arr[index]]
  emit('update:modelValue', arr)
}
</script>
